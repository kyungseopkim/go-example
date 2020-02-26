package main

import (
	"encoding/binary"
	"encoding/json"
	"log"
)

type MessageBody struct {
	dataReader *DataReader     `json:"-"`
	Vin        string          `json:"vin"`
	Packet     []MessagePacket `json:"packet"`
}

type MessagePacket struct {
	Timestamp int64            `json:"timestamp"`
	Usec      int32            `json:"usec"`
	Count     int32            `json:"count"`
	Payload   []MessagePayload `json:"payload"`
}

type MessagePayload struct {
	Msgid   int32  				`json:"msgid"`
	Payload []byte 				`json:"payload"`
}

func NewMessageBody(content []byte) (*MessageBody, error) {
	r := MessageBody{ dataReader: NewDataReader(content, binary.BigEndian)}
	if length, err := r.dataReader.GetByte(); err == nil {
		if chunk, err := r.dataReader.GetBytes(int32(length)); err== nil {
			r.Vin = string(chunk)
			for {
				if packet, e1 := NewMessagePacket(r.dataReader); e1 == nil {
					r.Packet = append(r.Packet, *packet)
				} else {
					break
				}
			}
			return &r, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}


func NewMessagePacket(reader *DataReader) (*MessagePacket, error) {
	if ts, e1 := reader.GetUInt64(); e1 == nil {
		if usec, e2 := reader.GetUInt32(); e2 == nil {
			if count, e3 := reader.GetUInt32(); e3 == nil {
				payloads := make([]MessagePayload, 0)
				for i:= uint32(0); i<count; i++ {
					if payload, e4 := NewMessagePayload(reader); e4 == nil {
						payloads = append(payloads, *payload)
					} else {
						break
					}
				}
				return &MessagePacket{Timestamp: int64(ts), Usec:int32(usec), Count:int32(count), Payload:payloads}, nil
			} else {
				return nil, e3
			}
		} else {
			return nil, e2
		}
	} else {
		return nil, e1
	}
}

func NewMessagePayload(reader *DataReader) (*MessagePayload, error)  {
	if id, err := reader.GetUInt32(); err == nil {
		if length, e1 := reader.GetUInt32(); e1 == nil {
			if payload, e2 := reader.GetBytes(int32(length)); e2 == nil {
				return &MessagePayload{Msgid: int32(id), Payload:payload}, nil
			} else {
				return nil, e2
			}
		} else {
			return nil, e1
		}
	} else {
		return nil, err
	}
}

func (payload MessagePayload) String() string  {
	if p, err := json.Marshal(payload); err == nil {
		return string(p)
	} else {
		log.Panic(err)
	}
	return ""
}

func (packet MessagePacket) String() string  {
	if p, err := json.Marshal(packet); err == nil {
		return string(p)
	} else {
		log.Fatal(err)
	}
	return ""
}


func (msg MessageBody) String() string {
	if m, err := json.Marshal(msg); err == nil {
		return string(m)
	} else {
		log.Fatal(err)
	}
	return ""
}
//func getVin(content []byte) :
