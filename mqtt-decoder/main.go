package main

import (
	"encoding/json"
	"flag"
	"github.com/BurntSushi/toml"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/kyungseopkim/decode"
	"log"
	"time"

	//"time"
)

type Config struct {
	Mqtt 	MqttConf
	Influx 	InfluxConfig
	Arxml   ArxmlConfig
}

type Message struct {
	Topic 	string 		`json:"topic"`
	Payload []byte 		`json:"payload"`
}

func (m Message) String() string  {
	content, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(content)
}


func main()  {

	configFile := flag.String("config", "config.toml", "configuration toml file")
	flag.Parse()

	var conf Config
	if _, err := toml.DecodeFile(*configFile, &conf); err != nil {
		log.Fatal(err)
	}

	mqttConf := conf.Mqtt
	influxConf := conf.Influx

	log.Printf("%#v\n", mqttConf)
	log.Printf("%#v\n", influxConf)
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)

	dbc := decode.ArxmlReader(conf.Arxml.Filename)

	opts := MQTT.NewClientOptions().AddBroker(mqttConf.connection()).SetClientID(mqttConf.ClientId)
	opts.SetKeepAlive(time.Duration(mqttConf.KeepAlive) * time.Second)
	//opts.SetDefaultPublishHandler(f)
	//opts.SetPingTimeout(1 * time.Second)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		//choke <- [2]string{msg.Topic(), string(msg.Payload())}
		log.Println(msg.Topic())
	})

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}

	influxClient, err := influx.NewHTTPClient(influx.HTTPConfig{Addr:"http://"})
	if err != nil {
		log.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer influxClient.Close()

	choke := make(chan Message)
	callback := func(client MQTT.Client, msg MQTT.Message) {
		log.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", msg.Topic(), msg.Payload())
		choke <- Message{ Topic:msg.Topic(), Payload:msg.Payload()}
	}

	if token := client.Subscribe(mqttConf.Topic, byte(1), callback); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	buffer := make([]*influx.Point, 0, 1000)
	for {
		msg := <- choke

		messages:=decode.Decompress(msg.Payload)
		if msgBody, err := decode.NewMessageBody(messages); err != nil {
			log.Println(err)
			continue
		} else {
			decoder := decode.NewDecoder(dbc,msgBody)
			for _, signal := range decoder.Decode() {
				tags := make(map[string]string)
				fields := make(map[string]interface{})
				tags["signal_name"] = signal.SignalName
				tags["vin"] = signal.Vin
				tags["msg_name"] = signal.MsgName
				tags["msg_id"] = string(signal.MsgId)
				nano := (signal.Timestamp % 1000) * 1000000 // milliseconds
				point, err := influx.NewPoint(influxConf.Measurement, tags, fields, time.Unix(int64(signal.Epoch), nano))
				if err != nil {
					continue
				}
				buffer = append(buffer, point)
				if len(buffer) == cap(buffer) {
					if bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
						Database:        influxConf.DB,
						Precision:       influxConf.Precision,
						RetentionPolicy: influxConf.Retention,
					}); err == nil {
						bp.AddPoints(buffer)
						influxClient.Write(bp)
						buffer = buffer[:0]
					}
				}
			}
		}
	}

	client.Disconnect(250)
	log.Println("Subscriber Disconnected")

}
