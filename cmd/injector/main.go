package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os/exec"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var (
	videoPath = flag.String("video", "video.mp4", "Path to the input video file")
	width     = flag.Int("width", 64, "Video width to scale to (max 255)")
	height    = flag.Int("height", 48, "Video height to scale to (max 255)")
	fps       = 30
)

type RGB struct {
	R, G, B uint8
}

var palette [256]RGB

func init() {
	basic := []RGB{
		{0, 0, 0}, {0, 0, 170}, {0, 170, 0}, {0, 170, 170},
		{170, 0, 0}, {170, 0, 170}, {170, 85, 0}, {170, 170, 170},
		{85, 85, 85}, {85, 85, 255}, {85, 255, 85}, {85, 255, 255},
		{255, 85, 85}, {255, 85, 255}, {255, 255, 85}, {255, 255, 255},
	}

	for i, c := range basic {
		palette[i] = c
	}

	idx := 16
	vals := []uint8{0, 95, 135, 175, 215, 255}
	for _, r := range vals {
		for _, g := range vals {
			for _, b := range vals {
				palette[idx] = RGB{r, g, b}
				idx++
			}
		}
	}

	for i := range 24 {
		v := uint8(8 + i*10)
		palette[idx] = RGB{v, v, v}
		idx++
	}
}

func findClosestColor(r, g, b uint8) int {
	minDist := math.MaxFloat64
	bestIdx := 0

	for i, c := range palette {
		dr := float64(r) - float64(c.R)
		dg := float64(g) - float64(c.G)
		db := float64(b) - float64(c.B)

		dist := dr*dr + dg*dg + db*db

		if dist < minDist {
			minDist = dist
			bestIdx = i
		}
	}
	return bestIdx
}

func SendPixel(conn *icmp.PacketConn, x, y, color int) {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   1234,
			Seq:  1,
			Data: []byte("GAMBIARRA"),
		},
	}

	bin, err := msg.Marshal(nil)
	if err != nil {
		return
	}

	targetIP := net.IPv4(10, byte(x), byte(y), byte(color))
	_, _ = conn.WriteTo(bin, &net.IPAddr{IP: targetIP})
}

func main() {
	flag.Parse()

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("Failed to listen on raw socket (Did you run with sudo?): %v", err)
	}
	defer conn.Close()

	log.Printf("🌈 Starting 256-COLOR Network-GPU Injector")
	log.Printf("▶️  Video: %s | Resolution: %dx%d | FPS: %d", *videoPath, *width, *height, fps)

	cmd := exec.Command("ffmpeg",
		"-i", *videoPath,
		"-vf", fmt.Sprintf("scale=%d:%d", *width, *height),
		"-f", "image2pipe",
		"-pix_fmt", "rgb24",
		"-vcodec", "rawvideo",
		"-loglevel", "quiet",
		"-")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create FFmpeg pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start FFmpeg: %v", err)
	}
	defer cmd.Process.Kill()

	framePixels := *width * *height
	frameBuffer := make([]byte, framePixels*3)

	prevFrame := make([]byte, framePixels)

	for i := range prevFrame {
		prevFrame[i] = 255
	}

	frameDuration := time.Second / time.Duration(fps)
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()

	frameIndex := 0

	log.Println("📡 Broadcasting 256-COLOR frames via ICMP...")

	for {
		_, err := io.ReadFull(stdout, frameBuffer)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Println("✅ End of video stream reached.")
				break
			}
			log.Fatalf("Error reading frame: %v", err)
		}

		<-ticker.C
		pixelsSent := 0

		for y := 0; y < *height; y++ {
			for x := 0; x < *width; x++ {
				pixelIdx := y*(*width) + x
				byteIdx := pixelIdx * 3

				r := frameBuffer[byteIdx]
				g := frameBuffer[byteIdx+1]
				b := frameBuffer[byteIdx+2]

				colorID := findClosestColor(r, g, b)

				if prevFrame[pixelIdx] != byte(colorID) {
					SendPixel(conn, x, y, colorID)
					prevFrame[pixelIdx] = byte(colorID)
					pixelsSent++

					if pixelsSent%800 == 0 {
						time.Sleep(1 * time.Microsecond)
					}
				}

			}
		}

		frameIndex++
		if frameIndex%(fps) == 0 {
			log.Printf("🎞️  Processed Frame %d | Packets injected: %d", frameIndex, pixelsSent)
		}
	}

	cmd.Wait()
	log.Println("🔌 Injection complete. Shutting down.")
}
