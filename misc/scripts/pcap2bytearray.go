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

func run(src gopacket.PacketDataSource) {
	dec := gopacket.DecodersByLayerName["Ethernet"]
	source := gopacket.NewPacketSource(src, dec)

	fmt.Println("var sample = [...][]byte{")
	for packet := range source.Packets() {
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}
		tcp := tcpLayer.(*layers.TCP)
		payload := tcpLayer.LayerPayload()
		if len(payload) == 0 {
			continue
		}
		fmt.Println(" {")
		fmt.Print("    ")
		if _, err := strconv.Atoi(tcp.SrcPort.String()); err != nil {
			//no brackets, count it as server
			fmt.Print(fromServer)
		} else {
			fmt.Print(fromClient)
		}
		fmt.Println(", //Direction")
		s := "    "
		for i, b := range payload {
			s += strconv.Itoa(int(b)) + ", "
			if i%12 == 0 && i != 0 {
				s += "\n    "
			}
		}
		fmt.Println(s)
		fmt.Println(" },")
	}
	fmt.Println("}")
}
