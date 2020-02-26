package main

import (
	"encoding/binary"
	"errors"
)

type DataReader struct {
	data  	[]byte
	pos 	int32
	endian  binary.ByteOrder
}

var IndexError = errors.New("index out of range")

func NewDataReader(data []byte, endian binary.ByteOrder) *DataReader  {
	r := DataReader{ data: data, pos: 0, endian: endian }
	return &r
}

func (reader DataReader) checkRange(length int32)  (bool, error) {
	if len(reader.data) < int(reader.pos + length) {
		return false, IndexError
	}
	return true, nil
}


func (reader *DataReader) GetByte() (byte, error) {
	ok, err := reader.checkRange(1)
	if ok {
		r := reader.data[reader.pos]
		reader.pos++
		return r, nil
	} else {
		return 0, err
	}
}

func (reader *DataReader) GetUInt16() (uint16, error) {
	ok, err := reader.checkRange(2)
	if ok {
		var ret uint16 = 0
		switch reader.endian {
		case binary.BigEndian:
			ret = binary.BigEndian.Uint16(reader.data[reader.pos:])
		case binary.LittleEndian:
			ret = binary.LittleEndian.Uint16(reader.data[reader.pos:])
		}
		reader.pos += 2
		return ret, nil
	} else {
		return 0, err
	}
}

func (reader *DataReader) GetUInt32() (uint32, error) {
	ok, err := reader.checkRange(4)
	if ok {
		var ret uint32 = 0
		switch reader.endian {
		case binary.BigEndian:
			ret = binary.BigEndian.Uint32(reader.data[reader.pos:])
		case binary.LittleEndian:
			ret = binary.LittleEndian.Uint32(reader.data[reader.pos:])
		}
		reader.pos += 4
		return ret, nil
	} else {
		return 0, err
	}
}

func (reader *DataReader) GetUInt64() (uint64, error) {
	ok, err := reader.checkRange(8)
	if ok {
		var ret uint64 = 0
		switch reader.endian {
		case binary.BigEndian:
			ret = binary.BigEndian.Uint64(reader.data[reader.pos:])
		case binary.LittleEndian:
			ret = binary.LittleEndian.Uint64(reader.data[reader.pos:])
		}
		reader.pos += 8
		return ret, nil
	} else {
		return 0, err
	}
}

func (reader *DataReader) GetBytes(length int32) ([]byte, error) {
	ok, err := reader.checkRange(length)
	if ok {
		r := make([]byte, length)
		copy(r, reader.data[reader.pos:])
		reader.pos += length
		return r, nil
	} else {
		return []byte{}, err
	}
}

