package main

import (
	"fmt"
	"net"
	"time"
)

type ClientInfo struct {
    Address *net.UDPAddr
    Username []byte
    Message []byte
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

    // UDPのソケットを開く
	// ListenUDP acts like ListenPacket for UDP networks.
    conn, err := net.ListenUDP("udp", udpAddress)
    if err != nil {
        fmt.Println("Error opening UDP connection:", err)
        return
    }
    defer conn.Close()
    fmt.Println("Starting UDP server on", udpAddress)

    // バッファを作成
    buf := make([]byte, 4096)
    // クライアントの情報を保存するマップ
    clientInfos := make(map[string]*ClientInfo)

    for {
        // UDPソケットから読み取ったデータをバイトスライス（buf）に格納
        n, addr, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Error reading from UDP:", err)
            return
        }
        usernameLen := buf[0]
        username := buf[1:usernameLen+1]
        message := buf[usernameLen+1:n]

        // 受信したデータを出力
        fmt.Printf("Received from %s\n", addr)
        fmt.Printf("first byte: %d\n", usernameLen)
        fmt.Printf("username: %s\n", string(username))
        fmt.Printf("massage: %s\n", string(message))

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
            Username: username,
            Message: message,
            ReceivedTime: time.Now(),
        }
        for key, value := range clientInfos {
            // 一定時間経過したクライアント情報を削除
            if time.Since(value.ReceivedTime) > 100 * time.Second {
                delete(clientInfos, key)
            }
            fmt.Println("Client addr:", value.Address, "User:", string(value.Username))
            fmt.Println("Message:", string(value.Message))
        }
        fmt.Printf("Number of clients: %d\n", len(clientInfos))
    }
}