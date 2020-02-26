package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)
// Signal is the definition of Signal
type Signal struct {
	BitLength   int32   	`json:"bit_length"`
	Factor      float32 	`json:"factor"`
	IsBigEndian bool    	`json:"is_big_endian"`
	IsSigned    bool    	`json:"is_signed"`
	Name        string  	`json:"name"`
	Offset      float32 	`json:"offset"`
	StartBit    int32   	`json:"start_bit"`
	Minimum     float32 	`json:"minimum"`
	Maximum     float32 	`json:"maximum"`
	Unit        string  	`json:"unit"`
	RecvNodes   string  	`json:"recv_nodes"`
	ValDesc     string  	`json:"val_desc"`
}

// Message is the definition of Message
type Message struct {
	Id              int32    `json:"id"`
	Description     string   `json:"description"`
	IsExtendedFrame bool     `json:"is_extended_frame"`
	Name            string   `json:"name"`
	Length          int32    `json:"length"`
	Signals         []Signal `json:"signals"`
}

// Messages is a wrapper of Message Array
type Messages struct {
	Messages []Message
}

// ArxmlReader returns New Arxml Reader
func ArxmlReader(fileName string) *Messages {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(file)
	var m Messages
	err = dec.Decode(&m)
	if err != nil {
		log.Fatal(err)
	}
	return &m
}

// SortByBitStartOffset is doing sorting by startOffset
func SortByBitStartOffset(messages *Messages)  {
	for _, m := range messages.Messages {
		sort.SliceStable(m.Signals, func(i, j int) bool {
			return m.Signals[i].StartBit < m.Signals[j].StartBit
		})
	}
}

// GetLooup return a lookup table
func (msg Messages) GetLookup() map[int32]Message {
	ret := make(map[int32]Message)
	for _, m := range msg.Messages {
		ret[m.Id]=m
	}
	return ret
}
