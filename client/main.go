package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
    // サーバーのアドレスを設定
    serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8080")
    if err != nil {
        fmt.Println("Error resolving UDP address:", err)
        return
    }

    // UDP接続を作成
    conn, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        fmt.Println("Error dialing UDP:", err)
        return
    }
    defer conn.Close()

    // コマンドラインからユーザー名を読み取る
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter username: ")
    username, _ := reader.ReadString('\n')
    fmt.Println("Starting chat with", username)

    // サーバーからのメッセージを受信するためのゴルーチンを開始
    go func() {
        buf := make([]byte, 1024)
        for {
            n, _, err := conn.ReadFromUDP(buf)
            if err != nil {
                fmt.Println("Error reading from UDP:", err)
                return
            }
            fmt.Println("Received message:", string(buf[:n]))
        }
    }()

    // メッセージを送信するループ
    for {
        // コマンドラインからメッセージを読み取る
        fmt.Print("Enter message: ")
        text, _ := reader.ReadString('\n')

        // ユーザー名とメッセージを結合
        message := username + "::: " + text

        // メッセージが4096バイト以上の場合はエラーを表示
        if len([]byte(message)) > 4096 {
            fmt.Println("Error: Message is too large. Please enter a message of 4096 bytes or less.")
            continue
        }

        // メッセージをサーバーに送信
        _, err = conn.Write([]byte(message))
        if err != nil {
            fmt.Println("Error writing to UDP:", err)
            return
        }

        fmt.Println("Message sent: ", username)
    }
}