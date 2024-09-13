package main

import (
	"bufio"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"io"
	"os"
)

var show = true

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("CPU usage: 0%")
	systray.SetTooltip("PFinal南丞")
	Panel := systray.AddMenuItem("控制面板", "Panel")
	mQuit := systray.AddMenuItem("退出", "Quit the whole app")
	go func() {
		for {
			select {
			case <-Panel.ClickedCh:
				toggerPanel()
			case <-mQuit.ClickedCh:
				sendQuitMessage()
			}
		}
	}()
	mQuit.SetIcon(icon.Data)
	Panel.SetIcon(icon.Data)
	// 启动 TCP 服务器
	systray.Run(onReady, nil)
}

func ListenToMain() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("管道关闭")
				break // 管道关闭，结束读取
			}
			fmt.Println("读取标准输入失败:", err)
			continue
		}
		fmt.Println("从标准输入读取到:" + line)
		if line == "quit\n" {
			fmt.Println("收到退出请求")
			// 可以添加其他处理逻辑
			systray.Quit() // 根据需要执行退出操作
		}
	}
}

func sendQuitMessage() {
	// 发送退出消息给主程序
	_, _ = fmt.Fprintln(os.Stdout, "systray_run: quit") // 发送特定的退出消息
	systray.Quit()
}

var toggerPanel = func() {
	if show {
		show = false
		// 发送退出消息给主程序
		_, _ = fmt.Fprintln(os.Stdout, "systray_run: panel hide") // 发送特定的退出消息
	} else {
		show = true
		// 发送退出消息给主程序
		_, _ = fmt.Fprintln(os.Stdout, "systray_run: panel show") // 发送特定的退出消息
	}
}

func main() {
	onReady()
}
