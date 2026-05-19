# PaintICMP 🎨⚡

**PaintICMP** is a creative networking proof-of-concept that turns standard network traffic into real-time art and video. It allows users to draw on a live web canvas by sending **ICMP (Ping)** packets to specific IP addresses. 

In its latest version, the system has been supercharged to handle massive packet floods, allowing it to **stream actual video over ICMP pings**! 

> [Youtube Video](https://www.youtube.com/watch?v=-ZcB5TelY8k)

This project demonstrates the power of low-level packet capture (`gopacket` / `libpcap`) combined with modern, high-performance web technologies (Binary WebSockets, Canvas API) to create unconventional interactive experiences. Perfect for networking enthusiasts and developers!

---

## 🚀 Key Features

### 🎬 Video Over ICMP (New!)
Turn your network into a GPU! The project now includes a **Network-GPU Injector** that reads any video file, scales it down, maps its frames to a 256-color palette, and broadcasts it to the canvas in real-time purely via Ping requests.

### 🎨 256-Color Packet Drawing
- **Protocol-Based Art:** Use the standard `ping` command to draw pixels manually.
- **IP Encoding:** The destination IP address handles all the drawing logic: `10.<X>.<Y>.<Color_ID>`.
- **Extended Palette:** Upgraded from basic colors to a full Xterm-like **256-color palette**.

### ⚡ Extreme Real-Time Performance
- **Binary WebSockets:** Pixel data is sent to the browser via fast binary ArrayBuffers instead of heavy JSON text.
- **Memory Pooling:** Under the hood, Go leverages `sync.Pool` to reuse memory, preventing garbage collector lag even when thousands of ping packets arrive per second.
- **Live Canvas:** A 256x256 grid updated dynamically without any frame drops.

---

## 🛠️ Tech Stack

* **Language:** [Go (Golang)](https://go.dev/)
* **Web Framework:** [Gin Gonic](https://gin-gonic.com/)
* **Packet Capture:** [Gopacket](https://github.com/google/gopacket) & `libpcap`
* **Real-time:** [Gorilla WebSocket](https://github.com/gorilla/websocket) (Binary Mode)
* **Frontend:** HTML5 Canvas, Vanilla JavaScript
* **Video Processing:** `FFmpeg` (Used by the Injector)

---

## ⚙️ Configuration and Installation

### Prerequisites

1. **Go installed** (v1.25+).
2. **libpcap** installed on your system (Required for packet sniffing).
   - Ubuntu/Debian: `sudo apt-get install libpcap-dev`
3. **FFmpeg** installed (Required *only* if you want to use the Video Injector).
   - Ubuntu/Debian: `sudo apt-get install ffmpeg`

### 1. Clone the Repository
```bash
git clone [https://github.com/yourusername/PaintICMP.git](https://github.com/yourusername/PaintICMP.git)
cd PaintICMP
go mod tidy

```

### 2. Network Interface Configuration (Important ⚠️)

The code is currently set to listen on the `wlo1` interface. You **must** change this to match your active network interface (e.g., `eth0`, `wlan0`, `en0`).
Open `internal/icmp/icmp.go` and update the following line:

```go
// internal/icmp/icmp.go
handler, err := pcap.OpenLive("your_interface_name", 1600, true, pcap.BlockForever)

```

---

## ▶️ Running the Application

Because the application requires direct access to network devices to capture packets, **both scripts must be run with root privileges (`sudo`)**.

### Step 1: Start the PaintICMP Server

This runs the packet sniffer and hosts the Web Canvas.

```bash
sudo go run cmd/paintICMP/main.go

```

🌐 **Open your browser and navigate to:** `http://localhost:8080/`

### Step 2: Draw on the Canvas

**Option A: Manual Drawing (The Classic Way)**
Open a new terminal and send a ping to the `10.0.0.0/8` range.
*Format:* `10.<X_Coord>.<Y_Coord>.<Color_ID_0-255>`

Example: To draw a red pixel (Color ID 196) at coordinates X: 50, Y: 50:

```bash
ping -c 1 10.50.50.196

```

**Option B: Play a Video (The Awesome Way)**
Use the provided injector to stream a `.mp4` file directly to the canvas using ICMP packets!

```bash
sudo go run cmd/injector/main.go -video path/to/your_video.mp4 -width 128 -height 96

```

*(Note: Keep the resolution low (e.g., 64x48 or 128x96) depending on your machine's networking limits!)*

---

## 📂 Project Structure

* `cmd/paintICMP/`: Entry point for the main server and packet sniffer.
* `cmd/injector/`: Entry point for the FFmpeg video-to-ICMP injector.
* `internal/icmp/`: Logic for packet capturing and filtering using pcap.
* `internal/websocket/`: Manages WebSocket connections and broadcasts binary pixel data.
* `internal/common/`: Memory management (`sync.Pool`) for high-performance packet handling.
* `web/`: Contains the frontend (`index.html`) with the HTML5 Canvas logic and 256-color map.

---

## 📄 License

This project is open-source. Feel free to fork, modify, and feature it in your videos!
