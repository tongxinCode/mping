package multicast

import (
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

type Packet struct {
	TTL          int
	Port         int
	Address      net.IP
	RouterAlert  bool
	Raw          bool
	IGMPVersion  int // 1, 2, or 3
	Interface    *net.Interface
	Message      []byte
	Protocol     string // 'udp' or 'ip:2'/'ip4:2'
	TargetAddr   *net.UDPAddr
	LocalAddress *net.UDPAddr
	UdpConn      *net.UDPConn
	PacketConn   *ipv4.PacketConn
	IpConn       net.PacketConn
	RawConn      *ipv4.RawConn
	Padding      []byte
	TOS          int
}

// Create a Packet Struct init-instance function
func newPacket() *Packet {
	return &Packet{
		TTL:         50,
		RouterAlert: false,
		IGMPVersion: 3,
		Protocol:    "udp",
		TOS:         20,
	}
}

// Use struct to define a multicast packet
func Broadcast(address string, localaddress string) (*Packet, error) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	laddr, err := net.ResolveUDPAddr("udp4", localaddress)
	if err != nil {
		return nil, err
	}
	// make a new packet and conn
	p := newPacket()
	p.TargetAddr = addr
	p.LocalAddress = laddr
	p.UdpConn, err = net.DialUDP("udp", p.LocalAddress, addr)
	if err != nil {
		return nil, err
	}
	p.PacketConn = ipv4.NewPacketConn(p.UdpConn)
	p.PacketConn.SetMulticastTTL(p.TTL)
	p.PacketConn.SetTOS(p.TOS)
	err = p.PacketConn.SetMulticastLoopback(true)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// write buffer([]byte) into the connection in a loop
func SendLoop(address string, localaddress string, content_byte []byte, intervalns int, handler func(int, []byte)) error {
	p, err := Broadcast(address, localaddress)
	if err != nil || p.UdpConn == nil || p.PacketConn == nil {
		return err
	}

	for {
		p.UdpConn.Write(content_byte)
		time.Sleep(time.Duration(intervalns) * time.Nanosecond)
		handler(len(content_byte), content_byte)
	}
}

// write buffer([]byte) into the connection
func Send(p *Packet, content_byte []byte, intervalns int, handler func(int, []byte)) error {
	p.UdpConn.Write(content_byte)
	time.Sleep(time.Duration(intervalns) * time.Nanosecond)
	handler(len(content_byte), content_byte)
	return nil
}
