package main

import (
	"encoding/binary"
	"encoding/hex"
)

type Bitset []byte

func NewBitset() Bitset {
	return make(Bitset, 8)
}

func (bit Bitset) Set(index int32)  {
	bit[index/8] |= 1 << byte(index % 8)
}

func (bit Bitset)  Unset(index int32)  {
	bit[index/8] &= ^(1 << byte(index % 8))
}


func (bit Bitset) Sets(indexs ...int32) {
	for _, i := range indexs {
		bit.Set(i)
	}
}

func (bit Bitset) From(data []byte)  {
	for i:=0; i < len(data); i++ {
		bit[i] = data[i]
	}
}

func (bit Bitset) String () string {
	return hex.EncodeToString(bit[:])
}

func (bit Bitset) ToUint64() uint64  {
	return binary.LittleEndian.Uint64(bit)
}

func (bit Bitset) Get(index int32) byte {
	return bit[index/8] & (1 << byte(index % 8))
}

func (bit Bitset) GetRange(start int32, len int32) Bitset  {
	last := start + len
	ret := NewBitset()
	var index = int32(0)
	for i:= start; i< last; i++  {
		if bit.Get(i) != 0 {
			ret.Set(index)
		}
		index++
	}
	return ret
}