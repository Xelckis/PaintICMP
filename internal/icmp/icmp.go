package icmp

import (
	"log"
	"strings"

	"paint/internal/websocket"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func FilterICMP() {
	handler, err := pcap.OpenLive("wlo1", 1600, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	err = handler.SetBPFFilter("icmp")
	if err != nil {
		panic(err)
	}

	log.Println("Reading packets...")

	packetSource := gopacket.NewPacketSource(handler, handler.LinkType())
	for packet := range packetSource.Packets() {
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			ipAdress := strings.Split(ip.DstIP.String(), ".")
			if ipAdress[0] == "10" {
				log.Printf("Capturei esse pixel aqui: %s\n", ipAdress)
				pixel := websocket.Pixel{
					X:     ipAdress[1],
					Y:     ipAdress[2],
					Color: ipAdress[3],
				}
				websocket.PixelChan <- pixel
			}
		}
	}
	defer handler.Close()

}
