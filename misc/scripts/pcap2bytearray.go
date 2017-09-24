// +build !windows

// Converts pcap file to go byte array format to simpilify making test cases.

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const fromClient = "0"
const fromServer = "1"

func main() {
	if len(os.Args) == 2 {
		var handle *pcap.Handle
		var err error
		if handle, err = pcap.OpenOffline(os.Args[1]); err != nil {
			log.Fatal("PCAP OpenOffline error:", err)
		}
		run(handle)
	} else {
		println("go run pcam2bytearray.go file.pcap")
	}
}

func hex(num int) string {
	str := "0x"
	if num < 16 {
		str += "0"
	}
	return fmt.Sprintf("%s%X", str, int(num))
}

func run(src gopacket.PacketDataSource) {
	dec := gopacket.DecodersByLayerName["Ethernet"]
	source := gopacket.NewPacketSource(src, dec)
	var c, s string
	for packet := range source.Packets() {
		var str *string
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}
		tcp := tcpLayer.(*layers.TCP)
		payload := tcpLayer.LayerPayload()
		if len(payload) == 0 {
			continue
		}

		if _, err := strconv.Atoi(tcp.SrcPort.String()); err != nil {
			str = &s
		} else {
			str = &c
		}

		*str += " {\n    "
		for i, b := range payload {
			*str += hex(int(b)) + ", "
			if (i+1)%12 == 0 {
				*str += "\n    "
			}
		}
		*str += "\n },\n"
	}

	fmt.Println("c.Buffer = [][]byte{")
	fmt.Print(c)
	fmt.Println("}")
	fmt.Println("s.Buffer = [][]byte{")
	fmt.Print(s)
	fmt.Println("}")
}
