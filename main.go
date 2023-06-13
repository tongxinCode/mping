package main

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"mping/multicast"

	"golang.org/x/net/ipv4"
	"golang.org/x/text/encoding/simplifiedchinese"

	lua "github.com/yuin/gopher-lua"
)

const (
	usage = `mping version: mping/1.8.2
Usage: ./mping [-h] [-s sendGroup] [-r receiveGroup] [-l localAddress] [-S sourceAddress] [-m message] [-i interval] [-log path]

Options:
`
)

const (
	MAX_DATA_SIZE = 65504
	FIT_DATA_SIZE = 1472
)

var (
	help           bool
	test           bool
	realtime       bool
	hexdata        bool
	count          bool
	logPath        string
	protoPath      string
	sendAddress    string
	receiveAddress string
	localAddress   string
	sourceAddress  string
	content        string
	encoding       string
	contentByte    []byte
	interval       int
	dataSize       int
	sendLimit      int

	clock_start time.Time
	clock_end   time.Time
	clock_mutex bool

	bytes_send_sum     float32
	bytes_rev_sum      float32
	packet_rev_sum     uint32
	packet_rev_theory  uint32
	packet_number_cur  uint32
	packet_number_last uint32
	packet_number_send uint32

	rawlog *log.Logger

	luaState *lua.LState

	ipReg   *regexp.Regexp
	addrReg *regexp.Regexp
)

func init() {
	clock_mutex = false
	bytes_send_sum = 0
	bytes_rev_sum = 0
	packet_rev_sum = 0
	packet_rev_theory = 0
	ipReg, _ = regexp.Compile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
	addrReg, _ = regexp.Compile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}:(([2-9]\d{3})|([1-5]\d{4})|(6[0-4]\d{3})|(65[0-4]\d{2})|(655[0-2]\d)|(6553[0-5]))`)
	flagSettup()
}

func main() {
	flag.Parse()
	logSettup()
	luaSettup()
	processArgs()
	processCommands()
	if luaState != nil {
		defer luaState.Close()
	}
}

func msgReceiveHandler(cm *ipv4.ControlMessage, src net.Addr, n int, b []byte) {
	if cm != nil {
		log.Println(cm.String())
		packet_rev_sum++
		if packet_number_cur == 0 {
			packet_number_cur = binary.BigEndian.Uint32(b[0:4])
			packet_rev_theory = 1
		} else {
			packet_number_last = packet_number_cur
			packet_number_cur = binary.BigEndian.Uint32(b[0:4])
			packet_rev_theory = packet_rev_theory + (packet_number_cur - packet_number_last)
		}
		if count {
			log.Printf("Total packets received:%d/%d\n", packet_rev_sum, packet_rev_theory)
		} else {
			log.Printf("Total packets received:%d\n", packet_rev_sum)
		}
	}
	if !clock_mutex {
		clock_start = time.Now()
		clock_mutex = true
	} else {
		clock_end = time.Now()
		bytes_rev_sum = bytes_rev_sum + float32(n)
		rates_rev := bytes_rev_sum * 1000000000 / float32(clock_end.Sub(clock_start).Nanoseconds())
		if rates_rev < 1000 {
			log.Println(rates_rev, "Bps")
		} else if rates_rev < 1000000 {
			log.Println(rates_rev/1024, "KBps")
		} else if rates_rev < 1000000000 {
			log.Println(rates_rev/1024/1024, "MBps")
		}
	}
	log.Println(n, "bytes read from", src)
	if hexdata {
		rawlog.Println(hex.Dump(b[:n]))
	}
	if protoPath != "*.lua" {
		// 调用Lua函数
		err := luaState.CallByParam(lua.P{
			Fn:      luaState.GetGlobal("Decode"),
			NRet:    0,
			Protect: true,
		}, lua.LString(b[:n]))
		if err != nil {
			log.Panic(err)
		}
	}
}

func msgSendHandler(n int, b []byte) {
	if !clock_mutex {
		clock_start = time.Now()
		clock_mutex = true
	} else {
		clock_end = time.Now()
		bytes_send_sum = bytes_send_sum + float32(n)
		rates_send := bytes_send_sum * 1000000000 / float32(clock_end.Sub(clock_start).Nanoseconds())
		if rates_send < 1000 {
			log.Println(rates_send, "Bps")
		} else if rates_send < 1000000 {
			log.Println(rates_send/1024, "KBps")
		} else if rates_send < 1000000000 {
			log.Println(rates_send/1024/1024, "MBps")
		}
	}
	log.Println(n, "bytes has been sent")
	if hexdata {
		rawlog.Println(hex.Dump(b[:n]))
	}
}

func getifi(addr string) (*net.Interface, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if host == "127.0.0.1" {
		return nil, nil
	}
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				ipv4 := ipReg.FindString(address.String())
				if ipv4 == host {
					ifi := &netInterfaces[i]
					// index := netInterfaces[i].Index
					// ifi, err := net.InterfaceByIndex(index)
					// if err != nil {
					// 	return nil, err
					// }
					return ifi, nil
				}
			}
		}
	}
	return nil, nil
}

func logSettup() {
	// set the formatflag of log
	// log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetFlags(log.LstdFlags)
	// define the log writer
	if logPath != "/" {
		file := logPath + time.Now().Format("2006-01-02 15-04") + ".log"
		logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			log.Fatal(err)
		}
		writers := []io.Writer{
			logFile,
			os.Stdout,
		}
		fileAndStdoutWriter := io.MultiWriter(writers...)
		log.SetOutput(fileAndStdoutWriter)
		rawlog = log.New(fileAndStdoutWriter, "", 0)
	} else {
		rawlog = log.New(os.Stdout, "", 0)
	}
}

func luaSettup() {
	if protoPath != "*.lua" {
		// 创建一个新的Lua虚拟机
		luaState = lua.NewState()
		// 加载Lua文件
		err := luaState.DoFile(protoPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func flagSettup() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&test, "test", false, "send and receive locally to examinate a test(default false)")
	flag.BoolVar(&realtime, "time", false, "send real time as the content to examinate(default false)")
	flag.BoolVar(&hexdata, "x", false, "whether to show the hex data(default false)")
	flag.BoolVar(&count, "c", false, "whether to count Packet loss rate(default false)")
	flag.StringVar(&logPath, "log", "/", "[/tmp/] or [C:\\] determine whether to log, Path e.g ./, Forbidden /")
	flag.StringVar(&protoPath, "proto", "*.lua", "choose a lua script to decode/encode udp data, function Decode(dataBytes)/Encode() must be included")
	flag.StringVar(&sendAddress, "s", "239.255.255.255:9999", "[group:port] send packet to group")
	flag.StringVar(&receiveAddress, "r", "239.255.255.255:9999", "[group:port] receive packet from group")
	flag.StringVar(&localAddress, "l", "127.0.0.1:8888", "[ip[:port]] must choose your local using interface")
	flag.StringVar(&sourceAddress, "S", "127.0.0.1:8888", "[ip[:port]] must determine the peer source ip if using SSM")
	flag.StringVar(&content, "m", "Init Data", "[string] change the content of sending")
	flag.StringVar(&encoding, "e", "ascii", "['ascii','utf8','gbk'] change the content of sending")
	flag.IntVar(&interval, "i", 1000000000, "[number] change the interval between package sent (unit:Nanosecond)")
	flag.IntVar(&dataSize, "p", -1, "[number] the size of payload data(0 means use 1472 Bytes payloads)")
	flag.IntVar(&sendLimit, "C", -1, "[number] the limit number of sending packets(-1 means no limits)")
	flag.Usage = flagUsage
}

func flagUsage() {
	fmt.Fprintf(os.Stderr, usage)
	flag.PrintDefaults()
}

func processCommands() {
	var wg sync.WaitGroup
	if help {
		flag.Usage()
		return
	}
	// determine the selected interface
	ifi, err := getifi(localAddress)
	if ifi != nil {
		log.Println("The index of interface used is", ifi.Index+1)
		log.Println("The name of interface used is", ifi.Name)
	} else {
		log.Println("[Tips:determine your using interface IP]")
		log.Println("[Otherwise the result may be incorrect]")
	}
	if err != nil {
		log.Fatal(err)
	}
	if dataSize == -1 {
		if count {
			var data []byte = make([]byte, 4)
			if encoding == "ascii" {
				contentByte = strconv.AppendQuoteToASCII(data, content)
			} else if encoding == "utf8" {
				contentByte = append(data, []byte(content)...)
			} else if encoding == "gbk" {
				contentByteTmp, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(content))
				contentByte = append(data, contentByteTmp...)
			} else {
				log.Fatal("unsupported encoding format")
			}
		} else {
			if encoding == "ascii" {
				contentByte = []byte(strconv.QuoteToASCII(content))
			} else if encoding == "utf8" {
				contentByte = []byte(content)
			} else if encoding == "gbk" {
				contentByte, _ = simplifiedchinese.GBK.NewEncoder().Bytes([]byte(content))
			} else {
				log.Fatal("unsupported encoding format")
			}
		}
	} else if dataSize == 0 {
		dataSize = FIT_DATA_SIZE
		var data []byte = make([]byte, dataSize-len(content))
		if encoding == "ascii" {
			contentByte = strconv.AppendQuoteToASCII(data, content)
		} else if encoding == "utf8" {
			contentByte = append(data, []byte(content)...)
		} else if encoding == "gbk" {
			contentByteTmp, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(content))
			contentByte = append(data, contentByteTmp...)
		} else {
			log.Fatal("unsupported encoding format")
		}
	} else if dataSize > 0 && dataSize < 4 {
		if count {
			log.Fatal("small packet")
		} else {
			contentByte = make([]byte, dataSize)
		}
	} else if dataSize > len(content) && dataSize <= MAX_DATA_SIZE {
		var data []byte = make([]byte, dataSize-len(content))
		if encoding == "ascii" {
			contentByte = strconv.AppendQuoteToASCII(data, content)
		} else if encoding == "utf8" {
			contentByte = append(data, []byte(content)...)
		} else if encoding == "gbk" {
			contentByteTmp, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(content))
			contentByte = append(data, contentByteTmp...)
		} else {
			log.Fatal("unsupported encoding format")
		}
	} else if dataSize > MAX_DATA_SIZE {
		log.Fatal("big packet")
	}
	if (sendAddress != "239.255.255.255:9999") && (receiveAddress != "239.255.255.255:9999") {
		log.Println("Send to ", sendAddress)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			p, err := multicast.Broadcast(sendAddress, localAddress)
			if err != nil || p.UdpConn == nil || p.PacketConn == nil {
				log.Fatal(err)
			}
			for {
				packet_number_send++
				if sendLimit > 0 && sendLimit < math.MaxInt && packet_number_send > uint32(sendLimit) {
					return
				}
				if protoPath != "*.lua" {
					// 调用Lua函数
					err := luaState.CallByParam(lua.P{
						Fn:      luaState.GetGlobal("Encode"),
						NRet:    1,
						Protect: true,
					}, lua.LNil)
					if err != nil {
						log.Panic(err)
					}
					// 获取返回值
					returnValue := luaState.Get(-1)
					str := returnValue.String()
					// 发送Lua函数的返回值
					err = multicast.Send(p, []byte(str), interval, msgSendHandler)
					if err != nil {
						log.Fatal(err)
					}
					// 弹出返回值
					luaState.Pop(1)
					continue
				}
				if realtime {
					content = time.Now().Format("2006-01-02 15:04:05")
					if encoding == "ascii" {
						contentByte = strconv.AppendQuoteToASCII(contentByte[0:4], content)
					} else if encoding == "utf8" {
						contentByte = append(contentByte[0:4], []byte(content)...)
					} else if encoding == "gbk" {
						contentByteTmp, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(content))
						contentByte = append(contentByte[0:4], contentByteTmp...)
					} else {
						log.Fatal("unsupported encoding format")
					}
				}
				if count {
					binary.BigEndian.PutUint32(contentByte[0:4], packet_number_send)
				}
				err := multicast.Send(p, contentByte, interval, msgSendHandler)
				if err != nil {
					log.Fatal(err)
				}
			}
		}(&wg)
		log.Println("Receive from ", receiveAddress)
		err := multicast.Receive(receiveAddress, sourceAddress, ifi, msgReceiveHandler)
		if err != nil {
			log.Fatal(err)
		}
	} else if sendAddress != "239.255.255.255:9999" && (receiveAddress == "239.255.255.255:9999") {
		log.Println("Send to ", sendAddress)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			p, err := multicast.Broadcast(sendAddress, localAddress)
			if err != nil || p.UdpConn == nil || p.PacketConn == nil {
				log.Fatal(err)
			}
			for {
				packet_number_send++
				if sendLimit > 0 && sendLimit < math.MaxInt && packet_number_send > uint32(sendLimit) {
					return
				}
				if protoPath != "*.lua" {
					// 调用Lua函数
					err := luaState.CallByParam(lua.P{
						Fn:      luaState.GetGlobal("Encode"),
						NRet:    1,
						Protect: true,
					}, lua.LNil)
					if err != nil {
						log.Panic(err)
					}
					// 获取返回值
					returnValue := luaState.Get(-1)
					str := returnValue.String()
					// 发送Lua函数的返回值
					err = multicast.Send(p, []byte(str), interval, msgSendHandler)
					if err != nil {
						log.Fatal(err)
					}
					// 弹出返回值
					luaState.Pop(1)
					continue
				}
				if realtime {
					content = time.Now().Format("2006-01-02 15:04:05")
					if encoding == "ascii" {
						contentByte = strconv.AppendQuoteToASCII(contentByte[0:4], content)
					} else if encoding == "utf8" {
						contentByte = append(contentByte[0:4], []byte(content)...)
					} else if encoding == "gbk" {
						contentByteTmp, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(content))
						contentByte = append(contentByte[0:4], contentByteTmp...)
					} else {
						log.Fatal("unsupported encoding format")
					}
				}
				if count {
					binary.BigEndian.PutUint32(contentByte[0:4], packet_number_send)
				}
				err := multicast.Send(p, contentByte, interval, msgSendHandler)
				if err != nil {
					log.Fatal(err)
				}
			}
		}(&wg)
	} else if receiveAddress != "239.255.255.255:9999" && (sendAddress == "239.255.255.255:9999") {
		log.Println("Receive from ", receiveAddress)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err := multicast.Receive(receiveAddress, sourceAddress, ifi, msgReceiveHandler)
			if err != nil {
				log.Fatal(err)
			}
		}(&wg)
	}
	if test {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			p, err := multicast.Broadcast(sendAddress, localAddress)
			if err != nil || p.UdpConn == nil || p.PacketConn == nil {
				log.Fatal(err)
			}
			for {
				packet_number_send++
				if sendLimit > 0 && sendLimit < math.MaxInt && packet_number_send > uint32(sendLimit) {
					return
				}
				if realtime {
					content = time.Now().Format("2006-01-02 15:04:05")
					if encoding == "ascii" {
						contentByte = strconv.AppendQuoteToASCII(contentByte[0:4], content)
					} else if encoding == "utf8" {
						contentByte = append(contentByte[0:4], []byte(content)...)
					} else if encoding == "gbk" {
						contentByteTmp, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(content))
						contentByte = append(contentByte[0:4], contentByteTmp...)
					} else {
						log.Fatal("unsupported encoding format")
					}
				}
				if count {
					binary.BigEndian.PutUint32(contentByte[0:4], packet_number_send)
				}
				err := multicast.Send(p, contentByte, interval, msgSendHandler)
				if err != nil {
					log.Fatal(err)
				}
			}
		}(&wg)
		err = multicast.Receive(receiveAddress, sourceAddress, ifi, msgReceiveHandler)
		if err != nil {
			log.Fatal(err)
		}
	}
	wg.Wait()
}

func processArgs() {
	if !addrReg.MatchString(localAddress) {
		conn, err := net.ListenUDP("udp", nil)
		if err != nil {
			log.Fatal(err)
		}
		port := conn.LocalAddr().(*net.UDPAddr).Port
		localAddress = net.JoinHostPort(localAddress, strconv.Itoa(port))
		conn.Close()
	}
	if !addrReg.MatchString(sourceAddress) {
		sourceAddress = net.JoinHostPort(sourceAddress, "0")
	}
}
