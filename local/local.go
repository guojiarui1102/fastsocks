package local

import (
	"log"
	"net"

	"github.com/guojiarui1102/fastsocks"
)

type LsLocal struct {
	Cipher     *fastsocks.Cipher
	ListenAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

func NewLsLocal(password string, listenAddr, remoteAddr string) (*LsLocal, error) {
	bsPassword, err := fastsocks.ParsePassword(password)
	if err != nil {
		return nil, err
	}
	structListenAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	structRemoteAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return nil, err
	}
	return &LsLocal{
		Cipher:     fastsocks.NewCipher(bsPassword),
		ListenAddr: structListenAddr,
		RemoteAddr: structRemoteAddr,
	}, nil
}

func (local *LsLocal) Listen(didListen func(listenAddr *net.TCPAddr)) error {
	return fastsocks.ListenEncryptedTCP(local.ListenAddr, local.Cipher, local.handleConn, didListen)
}

func (local *LsLocal) handleConn(userConn *fastsocks.SecureTCPConn) {
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
