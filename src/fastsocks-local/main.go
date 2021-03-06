package main

import (
	"fmt"
	"log"
	"net"

	"github.com/guojiarui1102/fastsocks/local"
	"github.com/guojiarui1102/fastsocks/src"
)

// DefaultListenAddr 默认端口7448
const (
	DefaultListenAddr = ":7448"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)

	config := &src.Config{
		ListenAddr: DefaultListenAddr,
	}
	config.ReadConfig()
	config.SaveConfig()

	fsLocal, err := local.NewFsLocal(config.Password, config.ListenAddr, config.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(fsLocal.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`fastsocks-local: %s 启动成功，配置:
		本地监听地址:
		%s
		远程服务地址:
		%s
		密码:
		%s`, version, listenAddr, config.RemoteAddr, config.Password))
	}))
}
