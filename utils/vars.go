package utils

import "sync"

//变量
var (
	ThreadNum = 500
	//自1.9版本以后提供了sync.Map，支持多线程并发读写，比之前的加锁map性能要好一点。
	Result    *sync.Map

	Host    string
	Port    = "22,23,53,80-139"
	Mode    = "tcp"
	Timeout = 2
	OutPutFile string
)