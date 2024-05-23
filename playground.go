package main

import (
	"encoding/binary"
	"log"
	"os"
)

func main() {
	f, err := os.Create("file.bin")
	if err != nil {
		log.Fatal("Couldn't open file")
	}
	defer f.Close()

	var data = struct {
		n1 uint16
		n2 uint8
		n3 uint8
	}{1200, 2, 4}
	err = binary.Write(f, binary.LittleEndian, data)
	if err != nil {
		log.Fatal("Write failed")
	}
}
