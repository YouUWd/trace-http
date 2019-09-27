package main

import (
	"./utils"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

var (
	device      string = "lo0"
	snapshotLen int32  = 1024
	promiscuous bool   = false
	err         error
	timeout     time.Duration = 30 * time.Second
	handle      *pcap.Handle
)

func main() {
	device = *flag.String("i", "lo0", "default device to listen on")
	filter := flag.String("f", "host 127.0.0.1 and port 8080", "default filter (eg: host 127.0.0.1 and port 8080)")
	flag.Parse()

	fmt.Println("You are listening on", device, "with filter string \""+*filter+"\"")
	// Open device
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	// Set filter
	err = handle.SetBPFFilter(*filter)

	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		printPacketInfo(packet)
	}
}

func printPacketInfo(packet gopacket.Packet) {

	packet.Layer(layers.LayerTypeTCP)

	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		netFlow := packet.NetworkLayer().NetworkFlow()
		transFlow := packet.TransportLayer().TransportFlow()
		fmt.Println(packet.Metadata().Timestamp,
			netFlow.Src().String()+":"+transFlow.Src().String(), "-->", netFlow.Dst().String()+":"+transFlow.Dst().String())
		utils.ParsePayload(applicationLayer.Payload())
	}

	// Check for errors
	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}
}
