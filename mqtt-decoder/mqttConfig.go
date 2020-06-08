package main

import (
	"encoding/json"
	"fmt"
)

type MqttConf struct {
	Server 		string		`json:"server"`
	Port 		int			`json:"port"`
	Topic 		string		`json:"topic"`
	ClientId 	string		`json:"client_id"`
	KeepAlive   int			`json:"keep_alive"`
	Qos			byte		`json:"qos"`
}

func (m MqttConf) connection() string {
	return fmt.Sprint("tcp://", m.Server, ":", m.Port)
}

func (m MqttConf) String() string  {
	if cont, err := json.Marshal(m); err == nil {
		return string(cont)
	}
	return ""
}