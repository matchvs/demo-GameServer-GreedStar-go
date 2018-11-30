/*
 * @Company: Matchvs
 * @Author: Ville
 * @Date: 2018-11-28 14:30:33
 * @LastEditors: Ville
 * @LastEditTime: 2018-11-30 14:58:09
 * @Description: matchvs game server example
 */

package main

import (
	"demo-GameServer-GreedStar-go/app"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/matchvs/gameServer-go"
)

//程序函数入口
func main() {
	// 定义业务处理对象这个业务类需要 继承接口
	handler := &app.App{}
	// 创建 gameServer
	gsserver := matchvs.NewGameServer(handler, "")
	handler.SetPushHandler(gsserver.GetPushHandler())
	// 启动 gameSever 服务
	go gsserver.Start()
	//检测系统信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	processStr := <-sigCh
	gsserver.Stop()
	fmt.Printf("service close  %v", processStr)
}
