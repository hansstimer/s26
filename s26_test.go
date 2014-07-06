package s26

import (
	"net"
	"testing"
	"time"
)

func TestS26(t *testing.T) {

}

func TestWriteToUDP(t *testing.T) {
	listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	l, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer l.Close()

	message := "UDP test message"

	done := make(chan bool, 1)
	go func() {
		var b [576]byte
		i, err := l.Read(b[:])
		if err != nil {
			t.Fatalf("Read err: %v", err)
		}
		if string(b[:i]) != message {
			t.Fatalf("Expected: %v, got: %v\n", message, string(b[:i]))
		}
		done <- true
	}()

	go func(raddr string) {
		udpAddr, err := net.ResolveUDPAddr("udp", raddr)
		if err != nil {
			t.Fatalf("Resolve failed: %v", err)
		}

		c, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			t.Fatalf("Dial failed: %v", err)
		}
		defer c.Close()

		_, err = c.Write([]byte(message))
		if err != nil {
			t.Fatalf("Write failed: %v", err)
		}
	}(l.LocalAddr().String())

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for UCP read")
	}

}
