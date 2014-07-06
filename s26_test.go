package s26

import (
	"net"
	"runtime"
	"testing"
)

func TestS26(t *testing.T) {

}

func TestWriteToUDP(t *testing.T) {
	switch runtime.GOOS {
	case "plan9":
		t.Skipf("skipping test on %q", runtime.GOOS)
	}

	l, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer l.Close()

	testWriteToConn(t, l.LocalAddr().String())
	testWriteToPacketConn(t, l.LocalAddr().String())
}

func testWriteToConn(t *testing.T, raddr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	c, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer c.Close()

	_, err = c.Write([]byte("Connection-oriented mode socket"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
}

func testWriteToPacketConn(t *testing.T, raddr string) {
	c, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenPacket failed: %v", err)
	}
	defer c.Close()

	ra, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		t.Fatalf("ResolveUDPAddr failed: %v", err)
	}

	_, err = c.(*net.UDPConn).WriteToUDP([]byte("Connection-less mode socket"), ra)
	if err != nil {
		t.Fatalf("WriteToUDP failed: %v", err)
	}

	_, err = c.WriteTo([]byte("Connection-less mode socket"), ra)
	if err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}

	_, err = c.(*net.UDPConn).Write([]byte("Connection-less mode socket"))
	if err == nil {
		t.Fatal("Write should fail")
	}
}
