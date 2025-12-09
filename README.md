# PaintICMP

# English version

> *[Ler versão em Português](#versão-em-português)*

## About the Project

**PaintICMP** is a creative proof-of-concept application that turns network traffic into art. It allows users to draw on a real-time web canvas by sending **ICMP (Ping)** packets to specific IP addresses.

The system captures network traffic in real-time using `gopacket`, filters for specific ICMP requests, and decodes the destination IP address to determine the coordinates and color of the pixel to be drawn. The web interface is updated instantly via **WebSockets**.

This project demonstrates the power of low-level packet capture combined with modern web technologies to create unconventional interactive experiences.

## Key Features

### 🎨 Packet-Driven Drawing
- **Protocol-Based Art:** Use the standard `ping` command to draw pixels.
- **IP Encoding:** The destination IP address encodes the drawing logic: `10.<X>.<Y>.<Color>`.
- **Packet Sniffing:** Utilizes `libpcap` to inspect network traffic on the wire.

### ⚡ Real-Time Visualization
- **WebSocket Integration:** Instant propagation of drawn pixels to all connected web clients.
- **Live Canvas:** A 255x255 grid that updates dynamically as packets are received.

## Tech Stack

* **Language:** [Go (Golang)](https://go.dev/)
* **Web Framework:** [Gin Gonic](https://gin-gonic.com/)
* **Packet Capture:** [Gopacket](https://github.com/google/gopacket) & `libpcap`
* **Real-time:** [Gorilla WebSocket](https://github.com/gorilla/websocket)
* **Frontend:** HTML5, CSS Grid, Vanilla JavaScript.

## How it Works

The drawing logic is based on the destination IP address of the ICMP packet. The application listens for packets addressed to the `10.0.0.0/8` range.

**Format:** `10.<X>.<Y>.<Color>`

- **X:** Horizontal Coordinate (0-255)
- **Y:** Vertical Coordinate (0-255)
- **Color:** Color ID (0-8)

### Color Map
| ID | Color | ID | Color |
|----|-------|----|-------|
| 0  | Black | 5  | Orange|
| 1  | Blue  | 6  | Red   |
| 2  | Green | 7  | White |
| 3  | Yellow| 8  | Gray  |
| 4  | Purple|       |       |

**Example:** To draw a **Red** pixel at coordinates **50, 50**:
```bash
ping 10.50.50.6
```
## Configuration and Installation
### Prerequisites

- Go installed (v1.25+).

- libpcap installed on your system (Required for packet capture).

  - Ubuntu/Debian: ```sudo apt-get install libpcap-dev```

### 1. Clone the Repository
```Bash
git clone [https://github.com/yourusername/PaintICMP.git](https://github.com/yourusername/PaintICMP.git)
cd PaintICMP
```
### 2. Network Interface Configuration

Important: The code is currently hardcoded to listen on the ```wlo1``` interface. You may need to change this in ```internal/icmp/icmp.go``` to match your network interface (e.g., ```eth0```, ```en0```).
```go
// internal/icmp/icmp.go
handler, err := pcap.OpenLive("your_interface_name", 1600, true, pcap.BlockForever)
```
### 3. Running the Application

Since the application requires access to network devices for packet capturing, you must run it with root privileges (sudo).
```Bash
go mod tidy
sudo go run main.go
```
Access the canvas in your browser: ```http://localhost:8080/ws``` (Note: The frontend connects to the WebSocket endpoint).

Ensure you open the HTML file or serve the static content correctly if not embedded.
## Project Structure

- ```internal/icmp/```: Logic for packet capturing and filtering using pcap.

- ```internal/websocket/```: Manages WebSocket connections and broadcasts pixel data.

- ```web/```: Contains the frontend (index.html) for visualizing the grid.

- ```main.go```: Entry point, initializes the sniffer and the web server.

## License

This project is open-source. Please check the repository for specific license details.

------

# Versão em Português

> *[Read english version](#english-version)*

## Sobre o Projeto

**PaintICMP** é uma prova de conceito criativa que transforma tráfego de rede em arte. O projeto permite que os utilizadores desenhem numa tela web em tempo real enviando pacotes **ICMP (Ping)** para endereços IP específicos.

O sistema captura o tráfego de rede em tempo real usando `gopacket`, filtra pedidos ICMP específicos e descodifica o endereço IP de destino para determinar as coordenadas e a cor do píxel a ser desenhado. A interface web é atualizada instantaneamente via **WebSockets**.

Este projeto demonstra o poder da captura de pacotes de baixo nível combinada com tecnologias web modernas para criar experiências interativas não convencionais.

## Funcionalidades Principais

### 🎨 Desenho via Pacotes
- **Arte Baseada em Protocolo:** Utilize o comando padrão `ping` para desenhar píxeis.
- **Codificação via IP:** O endereço IP de destino define a lógica do desenho: `10.<X>.<Y>.<Cor>`.
- **Packet Sniffing:** Utiliza `libpcap` para inspecionar o tráfego de rede diretamente na interface.

### ⚡ Visualização em Tempo Real
- **Integração WebSocket:** Propagação instantânea dos píxeis desenhados para todos os clientes conectados.
- **Tela Viva:** Uma grelha de 255x255 que se atualiza dinamicamente conforme os pacotes são recebidos.

## Stack Tecnológica

* **Linguagem:** [Go (Golang)](https://go.dev/)
* **Web Framework:** [Gin Gonic](https://gin-gonic.com/)
* **Captura de Pacotes:** [Gopacket](https://github.com/google/gopacket) & `libpcap`
* **Real-time:** [Gorilla WebSocket](https://github.com/gorilla/websocket)
* **Frontend:** HTML5, CSS Grid, Vanilla JavaScript.

## Como Funciona

A lógica de desenho baseia-se no endereço IP de destino do pacote ICMP. A aplicação escuta pacotes endereçados à faixa `10.0.0.0/8`.

**Formato:** `10.<X>.<Y>.<Cor>`

- **X:** Coordenada Horizontal (0-255)
- **Y:** Coordenada Vertical (0-255)
- **Cor:** ID da Cor (0-8)

### Mapa de Cores
| ID | Cor      | ID | Cor       |
|----|----------|----|-----------|
| 0  | Preto    | 5  | Laranja   |
| 1  | Azul     | 6  | Vermelho  |
| 2  | Verde    | 7  | Branco    |
| 3  | Amarelo  | 8  | Cinzento  |
| 4  | Roxo     |    |           |

**Exemplo:** Para desenhar um píxel **Vermelho** nas coordenadas **50, 50**:
```bash
ping 10.50.50.6
```

## Configuração e Instalação
### Pré-requisitos

- Go instalado (v1.25+).

- libpcap instalado no sistema (Necessário para captura de pacotes).
    - Ubuntu/Debian: ```sudo apt-get install libpcap-dev```

### 1. Clonar o Repositório
```Bash
git clone [https://github.com/seuutilizador/PaintICMP.git](https://github.com/seuutilizador/PaintICMP.git)
cd PaintICMP
```

### 2. Configuração da Interface de Rede

Importante: O código atualmente define a interface wlo1 de forma fixa (hardcoded). Poderá ser necessário alterar isso no ficheiro internal/icmp/icmp.go para corresponder à sua interface de rede (ex: eth0, en0).
```Go
// internal/icmp/icmp.go
handler, err := pcap.OpenLive("nome_da_sua_interface", 1600, true, pcap.BlockForever)
```
### 3. Executar a Aplicação

Como a aplicação requer acesso aos dispositivos de rede para captura de pacotes, deve executá-la com privilégios de root (sudo).
```Bash
go mod tidy
sudo go run main.go
```
Aceda à tela no seu navegador: ```http://localhost:8080/ws``` (Nota: O frontend conecta-se a este endpoint WebSocket, certifique-se de abrir o ficheiro HTML ou servir o conteúdo estático).
## Estrutura do Projeto

- ```internal/icmp/```: Lógica para captura e filtragem de pacotes usando pcap.

- ```internal/websocket/```: Gere conexões WebSocket e transmite dados dos píxeis.

- ```web/```: Contém o frontend (index.html) para visualização da grelha.

- ```main.go```: Ponto de entrada, inicializa o sniffer e o servidor web.

## Licença

Este projeto é open-source. Consulte o repositório para detalhes específicos sobre a licença.
