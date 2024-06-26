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

    fmt.Print("Enter username: ")
	// 標準入力からデータを取得する Scanner を作成
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	// ユーザー名を取得
	username := scanner.Text()
    fmt.Println("Starting chat with", username)

    // ユーザー名をバイト化し、長さも取得
    encodedUsername := []byte(username)
    usernameLenByte := byte(len(encodedUsername))
    namePart := append([]byte{usernameLenByte}, encodedUsername...)
    err_count := 0

    fmt.Println("username:", username)
    // fmt.Println("username length:", usernameLenByte)

    // サーバーからのメッセージを受信するためのゴルーチンを開始
    go func() {
        buf := make([]byte, 4096)
        for {
            n, _, err := conn.ReadFromUDP(buf)
            if err != nil {
                fmt.Println("\nError reading from UDP:", err)
                fmt.Print("Enter message: ")
                return
            }
            usernameLen := buf[0]
            username := buf[1:usernameLen+1]
            message := buf[usernameLen+1:n]
            fmt.Println("\nMessage from", string(username), ":", string(message))
            fmt.Print("Enter message: ")
        }
    }()

    // メッセージを送信するループ
    for {
        fmt.Print("Enter message: ")
        // コマンドラインからメッセージを読み取る
        scanner.Scan()
        if err := scanner.Err(); err != nil {
            fmt.Println("Error reading input:", err)
            return
        }
        text := scanner.Text()
        encodedText := []byte(text)

        // ユーザー名とメッセージを結合
        message := append(namePart, encodedText...)
        if len(message) > 4096 {
            fmt.Println("Message too long")
            continue
        }

        // メッセージをサーバーに送信
        _, err = conn.Write(message)
        if err != nil {
            fmt.Println("Error writing to UDP:", err)
            err_count++
            if err_count >= 3 {
                fmt.Println("Too many errors. Exiting.")
                os.Exit(1)
            }
            return
        }
        err_count = 0
        fmt.Print("Message sent!\n\n")
    }
}