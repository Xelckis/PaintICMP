package icmp

import (
	"log"
	"strconv"
	"strings"
	"time"

	"paint/internal/common"
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
	defer handler.Close()

	err = handler.SetBPFFilter("icmp")
	if err != nil {
		panic(err)
	}

	log.Println("Reading packets...")

	ticker := time.NewTicker(33 * time.Millisecond)

	packetSource := gopacket.NewPacketSource(handler, handler.LinkType())
	for {

		select {

		case <-ticker.C:
			if common.PixelPool != nil && len(*common.PixelPool) > 0 {
				websocket.GlobalHub.Broadcast <- common.PixelPool
				common.NilPixelPool()
			}

		case packet := <-packetSource.Packets():
			if packet == nil {
				continue
			}

			if common.PixelPool == nil {
				common.SetPixelPool()
			}

			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer != nil {
				ip, _ := ipLayer.(*layers.IPv4)
				ipAdress := strings.Split(ip.DstIP.String(), ".")

				X, err := strconv.Atoi(ipAdress[1])
				if err != nil {
					continue
				}

				Y, err := strconv.Atoi(ipAdress[2])
				if err != nil {
					continue
				}

				Color, err := strconv.Atoi(ipAdress[3])
				if err != nil {
					continue
				}

				if ipAdress[0] == "10" {
					if vram[X][Y] == ipAdress[3] {
					} else {
						vram[X][Y] = ipAdress[3]

						pixel := common.Pixel{
							X:     byte(X),
							Y:     byte(Y),
							Color: byte(Color),
						}

						*common.PixelPool = append(*common.PixelPool, pixel)
					}
				}

			}
		}
	}

}
