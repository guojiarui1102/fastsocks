package local

import (
	"fmt"
	"log"
	"net"

	"github.com/guojiarui1102/fastsocks"
)

type FsLocal struct {
	Cipher     *fastsocks.Cipher
	ListenAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

func NewFsLocal(password string, listenAddr, remoteAddr string) (*FsLocal, error) {
	bsPassword, err := fastsocks.ParsePassword(password)
	if err != nil {
		fmt.Printf("parse password failed, %s", err)
		return nil, err
	}
	structListenAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		fmt.Printf("resolve tcp listen addr failed, %s", err)
		return nil, err
	}
	structRemoteAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		fmt.Printf("resolve tcp remote addr failed, %s", err)
		return nil, err
	}
	return &FsLocal{
		Cipher:     fastsocks.NewCipher(bsPassword),
		ListenAddr: structListenAddr,
		RemoteAddr: structRemoteAddr,
	}, nil
}

func (local *FsLocal) Listen(didListen func(listenAddr *net.TCPAddr)) error {
	return fastsocks.ListenEncryptedTCP(local.ListenAddr, local.Cipher, local.handleConn, didListen)
}

func (local *FsLocal) handleConn(userConn *fastsocks.SecureTCPConn) {
	defer userConn.Close()

	proxyServer, err := fastsocks.DialEncryptedTCP(local.RemoteAddr, local.Cipher)
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyServer.Close()

	go func() {
		err := proxyServer.DecodeCopy(userConn)
		if err != nil {
			userConn.Close()
			proxyServer.Close()
		}
	}()
	userConn.EncodeCopy(proxyServer)
}
