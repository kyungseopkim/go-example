package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/gobwas/glob"
	"github.com/pierrec/lz4"
)

func readFileAll(filename string) []byte {
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	if data, err := ioutil.ReadAll(reader); err != nil {
		log.Fatal(err)
	} else {
		return data
	}
	return nil
}

func decompress(payload []byte) []byte {
	dst := make([]byte, 10*1024*1024)
	size, err := lz4.UncompressBlock(payload, dst)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return dst[:size]
}

func decodePayload(dbc *Messages, payload []byte) []MessageSignal {
	content := decompress(payload)
	if content != nil {
		if message, err := NewMessageBody(content); err == nil {
			decoder := NewDecoder(dbc, message)
			return decoder.decode()
		}
	}
	return make([]MessageSignal,0)
}

func main() {
	msg := ArxmlReader("parsed_arxml_GPS_31.json")
	SortByBitStartOffset(msg)

	if len(os.Args) < 2 {
		fmt.Printf("%s dir", os.Args[0])
		os.Exit(0)
	}

	baseDir := os.Args[1]
	files, error := ioutil.ReadDir(baseDir)
	if error != nil {
		log.Fatal(error)
	}

	pattern := glob.MustCompile("*.baby")
	for _, file := range files {
		if pattern.Match(file.Name()) {
			fmt.Println(file.Name())
			content := readFileAll(path.Join(baseDir, file.Name()))
			for _, signal := range decodePayload(msg, content) {
				fmt.Println(signal)
			}
		}
	}

}
