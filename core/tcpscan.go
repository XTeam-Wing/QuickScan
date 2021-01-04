package core

import (
	"QuickScan/utils"
	"fmt"
	"net"
	"time"
)

func TcpScan(ip string,port int)(string,int,error)  {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), time.Duration(utils.Timeout)*time.Second)

	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	return ip, port, err
}