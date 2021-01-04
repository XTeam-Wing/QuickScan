package cmd

import (
	"QuickScan/core"
	"QuickScan/utils"
	"fmt"
	"github.com/urfave/cli"
	"strings"
	"sync"
)

//命令行解析参数
var Cmd = cli.Command{
	Name: "Scan",
	Usage: "portscan tool",
	Action: GetArgs,
	Flags: []cli.Flag{
		stringFlag("host", "", "ip list"),
		stringFlag("port, p", "", "port list"),
		stringFlag("mode, m", "", "scan mode"),
		stringFlag("output, o", "scanres.txt", "output file"),
		intFlag("timeout, t", 3, "timeout"),
		intFlag("rate, r", 100, "thread"),
	},
}




func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func GetArgs(ctx *cli.Context)error  {
	if ctx.IsSet("host"){
		utils.Host = ctx.String("host")
	}
	if ctx.IsSet("port"){
		utils.Port = ctx.String("port")
	}
	if ctx.IsSet("mode") {
		utils.Mode = ctx.String("mode")
	}
	if ctx.IsSet("timeout") {
		utils.Timeout = ctx.Int("timeout")
	}
	if ctx.IsSet("rate"){
		utils.ThreadNum = ctx.Int("rate")
	}
	if ctx.IsSet("output"){
		utils.OutPutFile = ctx.String("output")
	}
	if strings.ToLower(utils.Mode) == "syn"{
		utils.CheckRoot()
	}

	ips,err := utils.ParseIpList(utils.Host)
	ports,err := utils.ParsePort(utils.Port)
	fmt.Println(ips,ports)
	tasks,tasknum := core.GenerateTask(ips,ports)
	fmt.Println("任务数量为:",tasknum)
	core.RunTask(tasks)
	core.PrintResult()
	return err
}



func init() {
	utils.Result = &sync.Map{}
}
