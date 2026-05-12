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
	fps       = flag.Int("fps", 30, "Target framerate")
)

type RGB struct {
	R, G, B uint8
}

var palette = []RGB{
	{0, 0, 0},       // 0: Black
	{0, 0, 170},     // 1: Dark Blue
	{0, 170, 0},     // 2: Dark Green
	{0, 170, 170},   // 3: Dark Cyan
	{170, 0, 0},     // 4: Dark Red
	{170, 0, 170},   // 5: Dark Magenta
	{170, 85, 0},    // 6: Brown
	{170, 170, 170}, // 7: Light Gray
	{85, 85, 85},    // 8: Dark Gray
	{85, 85, 255},   // 9: Bright Blue
	{85, 255, 85},   // 10: Bright Green
	{85, 255, 255},  // 11: Bright Cyan
	{255, 85, 85},   // 12: Bright Red
	{255, 85, 255},  // 13: Bright Magenta
	{255, 255, 85},  // 14: Yellow
	{255, 255, 255}, // 15: White
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

	log.Printf("🌈 Starting COLOR Network-GPU Injector")
	log.Printf("▶️  Video: %s | Resolution: %dx%d | FPS: %d", *videoPath, *width, *height, *fps)

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

	frameDuration := time.Second / time.Duration(*fps)
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()

	frameIndex := 0

	log.Println("📡 Broadcasting COLOR frames via ICMP...")

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
				}
			}
		}

		frameIndex++
		if frameIndex%(*fps) == 0 {
			log.Printf("🎞️  Processed Frame %d | Packets injected: %d", frameIndex, pixelsSent)
		}
	}

	cmd.Wait()
	log.Println("🔌 Injection complete. Shutting down.")
}
