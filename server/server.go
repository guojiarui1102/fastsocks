package server

import (
	"encoding/binary"
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
func (fsServer *FsServer) Listen(didListen func(listenAddr *net.TCPAddr)) error {
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
	localConn.EncodeWrite([]byte{0x05, 0x00})

	n, err := localConn.DecodeRead(buf)
	if err != nil || n < 7 {
		return
	}

	if buf[1] != 0x01 {
		return
	}

	var dIP []byte
	switch buf[3] {
	case 0x01:
		dIP = buf[4 : 4+net.IPv4len]
	case 0x03:
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))
		if err != nil {
			return
		}
		dIP = ipAddr.IP
	case 0x04:
		dIP = buf[4 : 4+net.IPv6len]
	default:
		return
	}
	dPort := buf[n-2:]
	dstAddr := &net.TCPAddr{
		IP:   dIP,
		Port: int(binary.BigEndian.Uint16(dPort)),
	}

	dstServer, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		return
	} else {
		defer dstServer.Close()
		dstServer.SetLinger(0)

		localConn.EncodeWrite([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}

	go func() {
		err := localConn.DecodeCopy(dstServer)
		if err != nil {
			localConn.Close()
			dstServer.Close()
		}
	}()
	(&fastsocks.SecureTCPConn{
		Cipher:          localConn.Cipher,
		ReadWriteCloser: dstServer,
	}).EncodeCopy(localConn)
}
