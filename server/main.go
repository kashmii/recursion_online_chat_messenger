package main

import (
	"fmt"
	"net"
	"time"
)

type ClientInfo struct {
    Address *net.UDPAddr
    ReceivedTime time.Time
}

func main() {
    // UDPのアドレスを設定
	// ResolveUDPAddr returns an address of UDP end point.
    udpAddress, err := net.ResolveUDPAddr("udp", ":8080")
    if err != nil {
        fmt.Println("Error resolving UDP address:", err)
        return
    }

	fmt.Println(udpAddress.Port)

    // UDPのソケットを開く
	// ListenUDP acts like ListenPacket for UDP networks.
    conn, err := net.ListenUDP("udp", udpAddress)
    if err != nil {
        fmt.Println("Error opening UDP connection:", err)
        return
    }
    defer conn.Close()

    // バッファを作成
    buf := make([]byte, 1024)
    // クライアントの情報を保存するマップ
    clientInfos := make(map[string]*ClientInfo)

    for {
        // データを読み取る
        n, addr, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Error reading from UDP:", err)
            return
        }

        // 受信したデータを出力
        fmt.Printf("Received from %s\n", addr)
        fmt.Printf("massage: %s\n", string(buf[:n]))

        // クライアントにメッセージを送信
        message := []byte("Hello from server")
        // メッセージが4096バイト以上の場合はエラーを表示
        if len([]byte(message)) > 4096 {
            fmt.Println("Error: Message is too large. Please enter a message of 4096 bytes or less.")
            continue
        }
        _, err = conn.WriteToUDP(message, addr)
        if err != nil {
            fmt.Println("Error sending message to client:", err)
            return
        }

        // クライアントの情報を保存
        clientInfos[addr.String()] = &ClientInfo{
            Address: addr,
            ReceivedTime: time.Now(),
        }
        for key, value := range clientInfos {
            if time.Since(value.ReceivedTime) > 100 * time.Second {
                delete(clientInfos, key)
            }
            fmt.Println("Client addr: ", value.Address, "Client time: ", value.ReceivedTime)
        }
        fmt.Printf("Number of clients: %d\n", len(clientInfos))
    }
}