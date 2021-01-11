package multicast

import (
	"errors"
	"net"
	"regexp"

	"golang.org/x/net/ipv4"
)

const (
	maxDatagramSize = 8 * 64 * 1024
)

// Receive is a function providing ASM and SSM receive function
func Receive(address string, sourceAddress string, ifi *net.Interface, handler func(*ipv4.ControlMessage, net.Addr, int, []byte)) error {
	group, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	if matched, _ := regexp.MatchString(`232(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`, group); matched {
		conn, err := JoinSSM(address, sourceAddress, ifi)
		if err != nil {
			return err
		}
		err = Listen(conn, handler)
		if err != nil {
			return err
		}
		err = LeaveSSM(address, sourceAddress, ifi, conn)
		if err != nil {
			return err
		}
	} else if matched, _ := regexp.MatchString(`2((2[4-9])|(3\d))(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`, group); matched {
		conn, err := JoinASM(address, ifi)
		if err != nil {
			return err
		}
		err = Listen(conn, handler)
		if err != nil {
			return err
		}
		err = LeaveASM(address, ifi, conn)
		if err != nil {
			return err
		}
	} else {
		err := errors.New("Check your multicast address.")
		return err
	}
	return nil
}

// Listen: loop and handle the log
func Listen(packetConn *ipv4.PacketConn, handler func(*ipv4.ControlMessage, net.Addr, int, []byte)) error {
	err := packetConn.SetMulticastLoopback(true)
	if err != nil {
		return nil
	}
	_ = packetConn.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true)
	buffer := make([]byte, maxDatagramSize)
	defer packetConn.Close()

	for {
		numBytes, controlMessage, src, err := packetConn.ReadFrom(buffer)
		if err != nil {
			return err
		}
		handler(controlMessage, src, numBytes, buffer)
	}
}

// JoinASM Join the ASM group
func JoinASM(address string, ifi *net.Interface) (*ipv4.PacketConn, error) {
	c, err := net.ListenPacket("udp", address)
	if err != nil {
		return nil, err
	}
	p := ipv4.NewPacketConn(c)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	err = p.JoinGroup(ifi, addr)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// LeaveASM Leave the ASM
func LeaveASM(address string, ifi *net.Interface, conn *ipv4.PacketConn) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	err = conn.LeaveGroup(ifi, addr)
	if err != nil {
		return err
	}
	return nil
}

// JoinSSM Join the SSM group
func JoinSSM(address string, sourceAddress string, ifi *net.Interface) (*ipv4.PacketConn, error) {
	c, err := net.ListenPacket("udp", address)
	if err != nil {
		return nil, err
	}
	p := ipv4.NewPacketConn(c)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	sourceAddr, err := net.ResolveUDPAddr("udp", sourceAddress)
	if err != nil {
		return nil, err
	}
	err = p.JoinSourceSpecificGroup(ifi, addr, sourceAddr)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// LeaveSSM: Leave the SSM group
func LeaveSSM(address string, sourceAddress string, ifi *net.Interface, conn *ipv4.PacketConn) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	sourceAddr, err := net.ResolveUDPAddr("udp", sourceAddress)
	if err != nil {
		return err
	}
	err = conn.LeaveSourceSpecificGroup(ifi, addr, sourceAddr)
	if err != nil {
		return err
	}
	return nil
}
