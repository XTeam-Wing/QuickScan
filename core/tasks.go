package core

import (
	"QuickScan/utils"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)
//[]map[string]int 相当于一个列表字典
func GenerateTask(ipList []net.IP,ports []int)([]map[string]int,int)  {
	//任务格式是ip:port
	tasks := make([]map[string]int,0)
	for _,ip := range ipList{
		for _,port := range ports{
			//字典赋值
			ipPort := map[string]int{ip.String():port}
			tasks = append(tasks,ipPort)
		}
	}
	//返回任务列表和任务数量
	return tasks,len(tasks)
}
func RunTask(tasks []map[string]int)  {
	//Create file
	var err1 error
	_, err1 = os.Create(utils.OutPutFile) //创建文件
	check(err1)
	//运用channal
	//WaitGroup 对象内部有一个计数器，最初从0开始，它有三个方法：Add(), Done(), Wait() 用来控制计数器的数量。Add(n) 把计数器设置为n ，Done() 每次把计数器-1 ，wait() 会阻塞代码的运行，直到计数器地值减为0。
	//wg的目的应该是为了保证程序能够有序的进行,结果正确
	wg := &sync.WaitGroup{}

	// 创建一个buffer为threadNum 的channel,woker数量
	taskChan := make(chan map[string]int, utils.ThreadNum)
	//每个worker起一个协程
	for i:=0;i<utils.ThreadNum;i++{
		go RunScan(taskChan,wg)
	}
	// 生产者，不断地往taskChan channel发送数据，直接channel阻塞
	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()
}

//WaitGroup对象不是一个引用类型，在通过函数传值的时候需要使用地址：
// 一定要通过指针传值，不然进程会进入死锁状态
func RunScan(taskchan chan map[string]int,wg *sync.WaitGroup)  {
	// 每个协程都从channel中读取数据后开始扫描并入库
	for task := range taskchan{
		for ip,port := range task{
			if strings.ToLower(utils.Mode) == "syn"{
				err := SaveResult(SynScan(ip,port))
				_ = err
			}else {
				err := SaveResult(TcpScan(ip,port))
				_ = err
			}
			wg.Done()
		}
	}
}

func SaveResult(ip string, port int, err error) error {
	if err != nil {
		return err
	}

	if port > 0 {
		//查询key,ip作为key
		v, ok := utils.Result.Load(ip)
		//这个ip已经有了的话,就存端口
		if ok {
			ports, ok1 := v.([]int)
			if ok1 {
				//存储进去
				ports = append(ports, port)
				utils.Result.Store(ip, ports)
			}
		} else {//没有这个ip需要新加进去
			ports := make([]int, 0)
			ports = append(ports, port)
			utils.Result.Store(ip, ports)
		}
	}

	return err
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func PrintResult() {
	//遍历打印出来
	utils.Result.Range(func(key, value interface{}) bool {
		fmt.Printf("ip:%v\n", key)
		fmt.Printf("open ports: %v\n", value)
		fmt.Println(strings.Repeat("-", 100))
		tmpString := fmt.Sprintf("%v", key)
		tmpPorts := fmt.Sprintf("%v",value)
		tmp1 := strings.ReplaceAll(tmpPorts,"[","")
		tmp2 := strings.ReplaceAll(tmp1,"]","")
		tmpPorts2 := strings.Split(tmp2," ")
		for _,v := range tmpPorts2{
			WriteString := tmpString+":"+v +"\n"
			SaveToFile(WriteString,utils.OutPutFile)
		}

		return true
	})

}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func SaveToFile(wireteString string,filename string)  {
	var f *os.File
	var err1 error
	f, err1 = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //打开文件
	check(err1)
	n, err1 := io.WriteString(f, wireteString) //写入文件(字符串)
	check(err1)
	_ = n
}