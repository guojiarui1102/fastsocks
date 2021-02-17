package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/guojiarui1102/fastsocks"
	"github.com/guojiarui1102/fastsocks/server"
	"github.com/guojiarui1102/fastsocks/src"
	"github.com/phayes/freeport"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)

	port, err := strconv.Atoi(os.Getenv("FASTSOCKS_SERVER_PORT"))
	if err != nil {
		port, err = freeport.GetFreePort()
	}
	if err != nil {
		port = 7448
	}
	config := &src.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		Password:   fastsocks.RandPassword(),
	}
	config.ReadConfig()
	config.SaveConfig()

	fsServer, err := server.NewFsServer(config.Password, config.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(fsServer.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(fmt.Sprintf(
			`fastsocks-server: %s 启动成功，配置如下：
			服务器监听地址：
			%s
			密码：
			%s`, version, listenAddr, config.Password))
	}))

}
