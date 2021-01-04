package core

import (
	"QuickScan/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"time"
)

// get the local ip and port based on our destination ip
func localIPPort(dstip net.IP) (net.IP, int, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", dstip.String()+":54321")
	if err != nil {
		return nil, 0, err
	}
	// We don't actually connect to anything, but we can determine
	// based on our destination ip what source ip we should use.
	if con, err := net.DialUDP("udp", nil, serverAddr); err == nil {
		if udpaddr, ok := con.LocalAddr().(*net.UDPAddr); ok {
			return udpaddr.IP, udpaddr.Port, nil
		}
	}
	return nil, -1, err
}

func SynScan(dstIp string, dstPort int) (string, int, error) {
	srcIp, srcPort, err := localIPPort(net.ParseIP(dstIp))
	dstAddrs, err := net.LookupIP(dstIp)
	if err != nil {
		return dstIp, 0, err
	}

	dstip := dstAddrs[0].To4()
	var dstport layers.TCPPort
	dstport = layers.TCPPort(dstPort)
	srcport := layers.TCPPort(srcPort)

	// Our IP header... not used, but necessary for TCP checksumming.
	ip := &layers.IPv4{
		SrcIP:    srcIp,
		DstIP:    dstip,
		Protocol: layers.IPProtocolTCP,
	}
	// Our TCP header
	tcp := &layers.TCP{
		SrcPort: srcport,
		DstPort: dstport,
		SYN:     true,
	}
	err = tcp.SetNetworkLayerForChecksum(ip)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		return dstIp, 0, err
	}

	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		return dstIp, 0, err
	}
	defer conn.Close()

	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstip}); err != nil {
		return dstIp, 0, err
	}

	// Set deadline so we don't wait forever.
	if err := conn.SetDeadline(time.Now().Add(time.Duration(utils.Timeout) * time.Second)); err != nil {
		return dstIp, 0, err
	}

	for {
		b := make([]byte, 4096)
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			return dstIp, 0, err
		} else if addr.String() == dstip.String() {
			// Decode a packet
			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
			// Get the TCP layer from this packet
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)

				if tcp.DstPort == srcport {
					if tcp.SYN && tcp.ACK {
						// log.Printf("%v:%d is OPEN\n", dstIp, dstport)
						return dstIp, dstPort, err
					} else {
						return dstIp, 0, err
					}
				}
			}
		}
	}
}

//
//func SynScan(ip string,port int)(string,int,error)  {
//	localIp,LocalPort,err := GetLocalIpPort(net.ParseIP(ip))
//	//判断ip是否解析
//	dstAddrs, err := net.LookupIP(ip)
//	if err != nil {
//		return ip, 0, err
//	}
//
//	dstip := dstAddrs[0].To4()
//	var dstport layers.TCPPort
//	dstport = layers.TCPPort(port)
//	srcport := layers.TCPPort(LocalPort)
//
//	// Our IP header... not used, but necessary for TCP checksumming.
//	ip2 := &layers.IPv4{
//		SrcIP:    localIp,
//		DstIP:    dstip,
//		Protocol: layers.IPProtocolTCP,
//	}
//	// Our TCP header
//	tcp := &layers.TCP{
//		SrcPort: srcport,
//		DstPort: dstport,
//		SYN:     true,
//	}
//	err = tcp.SetNetworkLayerForChecksum(ip2)
//
//	buf := gopacket.NewSerializeBuffer()
//	opts := gopacket.SerializeOptions{
//		ComputeChecksums: true,
//		FixLengths:       true,
//	}
//
//	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
//		return ip, 0, err
//	}
//
//	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
//	if err != nil {
//		return ip, 0, err
//	}
//	defer conn.Close()
//
//	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstip}); err != nil {
//		return ip, 0, err
//	}
//
//	// 设置超时时间
//	if err := conn.SetDeadline(time.Now().Add(time.Duration(utils.Timeout) * time.Second)); err != nil {
//		return ip, 0, err
//	}
//
//	for {
//		b := make([]byte, 4096)
//		n, addr, err := conn.ReadFrom(b)
//		if err != nil {
//			return ip, 0, err
//		} else if addr.String() == dstip.String() {
//			// Decode a packet
//			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
//			// Get the TCP layer from this packet
//			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
//				tcp, _ := tcpLayer.(*layers.TCP)
//
//				if tcp.DstPort == srcport {
//					if tcp.SYN && tcp.ACK {
//						log.Printf("%v:%d is OPEN\n", ip, dstport)
//						return ip, port, err
//					} else {
//						return ip, 0, err
//					}
//				}
//			}
//		}
//	}
//}
//
////根据目标ip获取本地ip和端口
//func GetLocalIpPort(dstip net.IP)(net.IP,int,error)  {
//	serverAddr, err := net.ResolveUDPAddr("udp", dstip.String()+":54321")
//	if err != nil {
//		return nil, 0, err
//	}
//	////实际上没建立连接，但可以基于我们的目标ip，我们应该使用什么源ip。
//	if con, err := net.DialUDP("udp", nil, serverAddr); err == nil {
//		if udpaddr, ok := con.LocalAddr().(*net.UDPAddr); ok {
//			return udpaddr.IP, udpaddr.Port, nil
//		}
//	}
//	return nil, -1, err
//}