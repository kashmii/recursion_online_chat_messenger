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

    // クライアントの情報を保存するマップ
    clientInfos := make(map[string]*ClientInfo)
    // バッファを作成
    buf := make([]byte, 4096)

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
        // クライアントにメッセージを送信
        sendMessageToOtherMembers(conn, buf, addr, clientInfos)
        buf = make([]byte, 4096)

        saveClientInfos(clientInfos, addr, username, message)
        RemoveInactiveClients(clientInfos)
        fmt.Printf("Number of clients: %d\n\n", len(clientInfos))
    }
}

func sendMessageToOtherMembers(conn *net.UDPConn, buf []byte, senderAddr *net.UDPAddr, clientInfos map[string]*ClientInfo) error {
    for key, v := range clientInfos {
        if key != senderAddr.String() {
            _, err := conn.WriteToUDP(buf, v.Address)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func saveClientInfos(clientInfos map[string]*ClientInfo, addr *net.UDPAddr, username []byte, message []byte) {
    // クライアントの情報を保存
    clientInfos[addr.String()] = &ClientInfo{
        Address: addr,
        Username: username,
        Message: message,
        ReceivedTime: time.Now(),
    }
}

func RemoveInactiveClients(clientInfos map[string]*ClientInfo) {
    deleteTime := 100 * time.Second
    for key, value := range clientInfos {
        // 一定時間経過したクライアント情報を削除
        if time.Since(value.ReceivedTime) > deleteTime {
            delete(clientInfos, key)
        }
    }
}