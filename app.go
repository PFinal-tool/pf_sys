package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os"
	"os/exec"
)

var (
	closeType = 0
)

// App struct
type App struct {
	ctx          context.Context
	sendPipeW    *os.File
	sendPipeR    *os.File
	receivePipeR *os.File
	receivePipeW *os.File
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// 创建两个管道：一个用于发送数据，一个用于接收数据
	sendPipeR, sendPipeW, err := os.Pipe() // 用于主程序向子进程发送数据
	if err != nil {
		fmt.Println("创建发送管道失败:", err)
		return
	}
	receivePipeR, receivePipeW, err := os.Pipe() // 用于子进程向主程序发送数据
	if err != nil {
		fmt.Println("创建接收管道失败:", err)
		return
	}
	a.sendPipeW = sendPipeW       // 主程序向子进程写
	a.sendPipeR = sendPipeR       // 主程序向子进程写
	a.receivePipeR = receivePipeR // 主程序从子进程读
	a.receivePipeW = receivePipeW // 主程序从子进程读

	go func() {
		// 启动 systray_run.go 程序
		cmd := exec.Command("go", "run", "./pak/sys_run/systray_run.go")
		cmd.Stdin = a.sendPipeR
		cmd.Stdout = a.receivePipeW
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Println("启动 systray_run.go 失败:", err)
			return
		}
		if err := cmd.Wait(); err != nil {
			fmt.Println("systray_run.go 执行失败:", err)
		}
	}()

	a.monitorPipe()
}

func (a *App) monitorPipe() {
	reader := bufio.NewReader(a.receivePipeR)
	for {
		line, err := reader.ReadString('\n') // 读取一行直到换行符
		if err != nil {
			if err == io.EOF {
				fmt.Println("管道关闭")
				break // 管道关闭，结束读取
			}
			fmt.Println("读取管道数据失败:", err)
			continue
		}
		// 处理从 systray_run.go 中接收到的输出
		fmt.Printf("从 systray_run.go 接收到: %s", line)
		switch line {
		case "systray_run: quit\n":
			fmt.Println("收到退出请求")
			closeType = 1
			runtime.Quit(a.ctx)
			break
		case "systray_run: panel show\n":
			fmt.Println("收到控制面板请求")
			runtime.WindowShow(a.ctx)
			break
		case "systray_run: panel hide\n":
			fmt.Println("收到控制面板请求")
			runtime.WindowHide(a.ctx)
			break
		}
	}
}

func (a *App) closeup(ctx context.Context) {
	fmt.Println(closeType)
	if closeType != 1 {
		// 发送退出消息给主程序
		_, _ = fmt.Fprintln(a.sendPipeW, "quit") // 发送特定的退出消息
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	runtime.WindowSetSize(a.ctx, 1000, 500)
	runtime.WindowReload(a.ctx)
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
