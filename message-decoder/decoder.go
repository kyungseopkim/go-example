package main

import (
	"log"
)

// Decoder
type Decoder struct {
	Dbc map[int32]Message
	Msg *MessageBody
}

// NewDecoder is
func NewDecoder(dbc *Messages, msg *MessageBody) *Decoder {
	return &Decoder{Dbc: dbc.GetLookup(), Msg: msg}
}

func (decoder Decoder) decode() []MessageSignal {
	ret := make([]MessageSignal, 0)
	vin := decoder.Msg.Vin
	for _, packet := range decoder.Msg.Packet {
		for _, payload := range packet.Payload {
			if msg, ok := decoder.Dbc[payload.Msgid]; ok {
				if len(payload.Payload) != int(msg.Length) {
					log.Printf("DEFINED MSG LENGTH %d != REAL %d\n", msg.Length, len(payload.Payload))
					continue
				}
				for _, s := range msg.Signals {
					milli := (packet.Timestamp * 1000000) + int64(packet.Usec)
					bitDecoder := NewBitDecoder(payload.Payload, s)
					log.Println(s)
					value := bitDecoder.GetValue()
					ret = append(ret, MessageSignal{MsgId: payload.Msgid, Timestamp: milli, Epoch: int32(packet.Timestamp),
						Vin: vin, MsgName: msg.Name, SignalName: s.Name, Value: value})
				}
			} else {
				log.Printf("MSG ID[%d] NOT DEFINED IN ARXML\n", payload.Msgid)
			}
		}
	}
	return ret
}

type BitDecoder struct {
	Data   []byte
	Signal Signal
}

func NewBitDecoder(payload []byte, signal Signal) *BitDecoder {
	return &BitDecoder{Data: payload, Signal: signal}
}

func reverseByteArray(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	return append(reverseByteArray(data[1:]), data[0])
}

func (decoder BitDecoder) GetValue() float32 {
	idx := decoder.Signal.StartBit / 64
	startByte := idx * 8
	sliceStart := startByte * 8
	var startBit  = decoder.Signal.StartBit - sliceStart
	last := startByte+8
	if int(last) > len(decoder.Data) {
		last = int32(len(decoder.Data))
	}
	var unit = decoder.Data[startByte : last]
	if decoder.Signal.IsBigEndian {
		unit = reverseByteArray(unit)
		startBit = 64 - startBit - decoder.Signal.BitLength
	}
	data := NewBitset()
	data.From(decoder.Data[startByte:last])
	value := data.GetRange(startBit, decoder.Signal.BitLength)

	return (float32(value.ToUint64()) * decoder.Signal.Factor) + decoder.Signal.Offset
}
