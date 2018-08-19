package freesocks

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
)

func ListenTcpServer(host string, port int) error {
	server, err := net.Listen("tcp4", host+":"+strconv.Itoa(port))

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Listen to tcp server on", server.Addr())
	defer server.Close()
	for {
		// listen for an incoming connection
		c, err := server.Accept()
		log.Println("Tcp Connection", c.LocalAddr(), c.RemoteAddr())
		if err != nil {
			return err
		}
		// handle connection
		go func(c net.Conn) {
			// step1. read username and password to verificate
			// step2. get request target information
			// step3. build remote connection
			// step4. perform bridge between client and remote server
			defer c.Close()
		}(c)
	}
	return nil
}

func readHeader(c net.Conn) (FreeSocketHeader, error) {
	var header FreeSocketHeader
	// read version and ip type
	verBuf := make([]byte, 1)
	if _, err := io.ReadFull(c, verBuf); err != nil {
		return header, err
	}
	header.version = readFst4Bits(verBuf[0])
	header.reserv = readLast4Bits(verBuf[0])

	// read user
	if userBytes, err := readBufferWithLen(c); err != nil {
		return header, err
	} else {
		header.user = userBytes
	}

	// read pass
	if passBytes, err := readBufferWithLen(c); err != nil {
		return header, err
	} else {
		header.pass = passBytes
	}

	// read host
	if hostBytes, err := readBufferWithLen(c); err != nil {
		return header, err
	} else {
		header.host = hostBytes
	}

	// read port
	portBytes := make([]byte, 2)
	if _, err := io.ReadFull(c, portBytes); err != nil {
		return header, err
	} else {
		port := binary.BigEndian.Uint16(portBytes)
		header.port = port
	}

	return header, nil
}

func readBufferWithLen(c net.Conn) ([]byte, error) {
	lenBytes := make([]byte, 1)
	if _, err := io.ReadFull(c, lenBytes); err != nil {
		return nil, err
	}

	len := int(lenBytes[0])
	if len == 0 {
		return nil, errors.New("missing content.")
	}

	cntBytes := make([]byte, len)
	if _, err := io.ReadFull(c, cntBytes); err != nil {
		return nil, err
	}

	return cntBytes, nil
}

func readFst4Bits(b byte) uint8 {
	return uint8((b & byte(240)) >> 4)
}

func readLast4Bits(b byte) uint8 {
	return uint8(b & byte(15))
}

type FreeSocketHeader struct {
	version uint8 // 0000
	reserv  uint8 // 0000

	user []byte // userLen 1 byte
	pass []byte // passLen 1 byte
	host []byte // hostLen 1 byte

	port uint16 // 2 bytes
}

/*
func validate(c net.Conn) error {
	bb := make([]byte, 2)
}
*/
