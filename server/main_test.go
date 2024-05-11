package main

import (
	"net"
	"testing"
	"time"
)

type MockPacketConn struct {
    readFromUDP func(b []byte) (int, *net.UDPAddr, error)
}

func (m *MockPacketConn) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
    // このメソッドは、mockConn := の後ろで定義された関数を呼び出す
    return m.readFromUDP(b)
}

func TestMain(t *testing.T) {
    // 長さ、名前、メッセージを定義し、それらの文字列をバイトスライスにする
    username := "Alice"
    usernameLenByte := byte(len(username))
    namePart := append([]byte{usernameLenByte}, []byte(username)...)
    message := "test message"

    mockConn := &MockPacketConn{
        readFromUDP: func(b []byte) (int, *net.UDPAddr, error) {
            data := append(namePart, []byte(message)...)
            copy(b, data)
            return len(data), &net.UDPAddr{}, nil
        },
    }

    buf := make([]byte, 4096)

    // UDPソケットから読み取ったデータをバイトスライス（buf）に格納
    n, _, err := mockConn.ReadFromUDP(buf)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    // TEST: usernameLenByte が正しい値を取得しているか
    if buf[0] != usernameLenByte {
        t.Errorf("Expected %d, got %d", len("test message"), n)
    }
    // TEST: username が正しくスライスから取得できているか
    if string(buf[1:buf[0]+1]) != username {
        t.Errorf("Expected %s, got %s", username, string(buf[1:buf[0]+1]))
    }
    // TEST: message が正しくスライスから取得できているか
    if string(buf[len(username)+1:n]) != "test message" {
        t.Errorf("Expected %s, got %s", "test message", string(buf[len(username)+1:n]))
    }
}

// TEST: uniqueAddrList関数 が重複のない正しいアドレスリストを返すか
func TestUniqueAddrList(t *testing.T) {
    dummyClient1 := &ClientInfo{
        Address: &net.UDPAddr{
            IP:   net.ParseIP("127.0.0.1"),
            Port: 8080,
        },
        Username: []byte("dummyUser1"),
        Message:  []byte("dummyMessage1"),
        ReceivedTime: time.Now(),
    }
    dummyClient2 := &ClientInfo{
        Address: &net.UDPAddr{
            IP:   net.ParseIP("127.0.0.2"),
            Port: 8081,
        },
        Username: []byte("dummyUser2"),
        Message:  []byte("dummyMessage2"),
        ReceivedTime: time.Now(),
    }
    dummyClient3 := &ClientInfo{
        Address: &net.UDPAddr{
            IP:   net.ParseIP("127.0.0.1"),
            Port: 8080,
        },
        Username: []byte("dummyUser3"),
        Message:  []byte("dummyMessage3"),
        ReceivedTime: time.Now(),
    }

    clientInfos := make(map[string]*ClientInfo)
    clientInfos[dummyClient1.Address.String()] = dummyClient1
    clientInfos[dummyClient2.Address.String()] = dummyClient2
    clientInfos[dummyClient3.Address.String()] = dummyClient3

    // Call the uniqueAddrList function
    addrList := uniqueAddrList(clientInfos)

    if len(addrList) != 2 {
        t.Errorf("Expected address list length of 2, got %d", len(addrList))
    }

    seen := make(map[string]bool)
    for _, addr := range addrList {
        addrStr := addr.String()
        seen[addrStr] = true
    }
    if len(seen) != 2 {
        t.Errorf("Expected 3 unique addresses, got %d", len(seen))
    }
}