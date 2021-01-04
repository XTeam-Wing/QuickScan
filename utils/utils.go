package utils

import (
	"fmt"
	"github.com/malfunkt/iprange"
	"net"
	"os"
	"strconv"
	"strings"
)

func IsRoot() bool {
	return os.Geteuid() == 0
}

func CheckRoot() {
	if !IsRoot() {
		fmt.Println("Please Run QuickScan With Root!")
		os.Exit(0)
	}
}


//解析ip
func ParseIpList(ips string)([]net.IP,error) {
	addrList, err := iprange.ParseList(ips)
	if err != nil{
		return nil, err
	}
	list := addrList.Expand()
	return list,err
}

// 解析Port
func ParsePort(port string) ([]int,error) {
	//Py的数组
	ports := []int{}
	if port =="" {
		//不等于nil一般是错误
		return ports,nil
	}
	ranges := strings.Split(port,",")
	//类似python的range
	for _,r :=range ranges{
		r = strings.TrimSpace(r)
		if strings.Contains(r,"-"){
			SplitPort := strings.Split(r,"-")
			if len(SplitPort)!= 2 {
				return nil, fmt.Errorf("port error: %s",r)
			}
			//string转int
			pStart,err := strconv.Atoi(SplitPort[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'",SplitPort[0])
			}
			pEnd,err := strconv.Atoi(SplitPort[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'",SplitPort[1])
			}
			if pStart > pEnd {
				return nil, fmt.Errorf("Invalid port range: %d-%d",pStart , pEnd)
			}
			//类似C++
			for i:=pStart;i<pEnd;i++{
				ports = append(ports, i)
			}
		}else {
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", r)
			} else {
				ports = append(ports, port)
			}
		}
	}
	return ports,nil
}

