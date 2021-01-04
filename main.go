package main

import (
	"QuickScan/cmd"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"runtime"
)

func main()  {
	//初始化命令行
	app := cli.NewApp()
	app.Name = "QuickScan"
	app.Version = "1.0 beta"
	app.Usage = "PortScanner 1.0 Beta"
	app.Commands = []cli.Command{cmd.Cmd}
	app.Flags = append(app.Flags,cmd.Cmd.Flags...)
	//接收参数
	err := app.Run(os.Args)
	if err !=nil{
		fmt.Println(err)
		os.Exit(0)
	}
}

func init() {
	//Go 1.5 版本之前，默认使用的是单核心执行。从 Go 1.5 版本开始，默认执行上面语句以便让代码并发执行，最大效率地利用 CPU。
	runtime.GOMAXPROCS(runtime.NumCPU())
}