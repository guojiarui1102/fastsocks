package fastsocks

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

const (
	bufSize = 1024
)

var bpool sync.Pool

func init() {
	bpool.New = func() interface{} {
		return make([]byte, bufSize)
	}
}

func bufferPoolGet() []byte {
	return bpool.Get().([]byte)
}

func bufferPoolPut(b []byte) {
	bpool.Put(b)
}

// SecureTCPConn 加密TCP
type SecureTCPConn struct {
	io.ReadWriteCloser
	Cipher *Cipher
}

// DecodeRead 读取加密数据，解密后放回bs
func (secureSocket *SecureTCPConn) DecodeRead(bs []byte) (n int, err error) {
	n, err = secureSocket.Read(bs)
	if err != nil {
		fmt.Printf("read failed, %s", err)
		return
	}
	secureSocket.Cipher.Decode(bs[:n])
	return
}

// EncodeWrite 把bs的数据加密后，写回
func (secureSocket *SecureTCPConn) EncodeWrite(bs []byte) (int, error) {
	secureSocket.Cipher.Encode(bs)
	return secureSocket.Write(bs)
}

// EncodeCopy 从对象池读取数据，加密后写入dst
func (secureSocket *SecureTCPConn) EncodeCopy(dst io.ReadWriteCloser) error {
	buf := bufferPoolGet()
	defer bufferPoolPut(buf)
	for {
		readCnt, errRead := secureSocket.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			}
			return nil
		}
		if readCnt > 0 {
			writeCnt, errWrite := (&SecureTCPConn{
				ReadWriteCloser: dst,
				Cipher:          secureSocket.Cipher,
			}).EncodeWrite(buf[0:readCnt])
			if errWrite != nil {
				return errWrite
			}
			if readCnt != writeCnt {
				return io.ErrShortWrite
			}
		}
	}
}

// DecodeCopy 从对象池读取加密后的数据，解密后写入dst
func (secureSocket *SecureTCPConn) DecodeCopy(dst io.Writer) error {
	buf := bufferPoolGet()
	defer bufferPoolPut(buf)
	for {
		readCnt, errRead := secureSocket.DecodeRead(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			}
			return nil
		}
		if readCnt > 0 {
			writeCnt, errWrite := dst.Write(buf[0:readCnt])
			if errWrite != nil {
				return errWrite
			}
			if readCnt != writeCnt {
				return io.ErrShortWrite
			}
		}
	}
}

// DialEncryptedTCP DialTCP
func DialEncryptedTCP(raddr *net.TCPAddr, cipher *Cipher) (*SecureTCPConn, error) {
	remoteConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		fmt.Printf("dial failed，%s", err)
		return nil, err
	}
	remoteConn.SetLinger(0)

	return &SecureTCPConn{
		ReadWriteCloser: remoteConn,
		Cipher:          cipher,
	}, nil
}

// ListenEncryptedTCP LisenTCP
func ListenEncryptedTCP(laddr *net.TCPAddr, cipher *Cipher, handleConn func(localConn *SecureTCPConn), didListen func(listenAddr *net.TCPAddr)) error {
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Printf("listen failed, %s", err)
		return err
	}
	defer listener.Close()

	if didListen != nil {
		go didListen(listener.Addr().(*net.TCPAddr))
	}

	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		localConn.SetLinger(0)
		go handleConn(&SecureTCPConn{
			ReadWriteCloser: localConn,
			Cipher:          cipher,
		})
	}
}
