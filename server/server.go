package server

import (
	"net"

	"github.com/guojiarui1102/fastsocks"
)

// FsServer struct
type FsServer struct {
	Cipher     *fastsocks.Cipher
	ListenAddr *net.TCPAddr
}

// NewFsServer 新建一个服务端
func NewFsServer(pw string, listenAddr string) (*FsServer, error) {
	bsPw, err := fastsocks.ParsePassword(pw)
	if err != nil {
		return nil, err
	}
	structListenAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &FsServer{
		Cipher:     fastsocks.NewCipher(bsPw),
		ListenAddr: structListenAddr,
	}, nil
}

// Listen 监听客户端请求
func (fsServer *FsServer) listen(didListen func(listenAddr *net.TCPAddr)) error {
	return fastsocks.ListenEncryptedTCP(fsServer.ListenAddr, fsServer.Cipher, fsServer.handleConn, didListen)
}

// handleConn 解析SOCKS5协议
func (fsServer *FsServer) handleConn(localConn *fastsocks.SecureTCPConn) {
	defer localConn.Close()
	buf := make([]byte, 256)

	_, err := localConn.DecodeRead(buf)
	if err != nil || buf[0] != 0x05 {
		return
	}
	//...
}
