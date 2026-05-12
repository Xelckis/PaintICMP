package icmp

import (
	"log"
	"strconv"
	"strings"

	"paint/internal/websocket"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var vram [255][255]string

/*
ipAdress[1] = X-cord
ipAdress[2] = Y-cord
ipAdress[3] = color
*/
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

			X, err := strconv.Atoi(ipAdress[1])
			if err != nil {
				log.Printf("Error converting X to int: %v", err)
				return
			}

			Y, err := strconv.Atoi(ipAdress[2])
			if err != nil {
				log.Printf("Error converting Y to int: %v", err)
				return
			}

			if ipAdress[0] == "10" {
				if vram[X][Y] == ipAdress[3] {
					log.Printf("Pixel color is equal, not sending request...")
				} else {
					log.Printf("Capturei esse pixel aqui: %s\n", ipAdress)
					vram[X][Y] = ipAdress[3]
					pixel := websocket.Pixel{
						X:     ipAdress[1],
						Y:     ipAdress[2],
						Color: ipAdress[3],
					}
					websocket.GlobalHub.Broadcast <- pixel
				}
			}
		}
	}
	defer handler.Close()

}
