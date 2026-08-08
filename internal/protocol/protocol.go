package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
)

func ReadInt32(conn net.Conn) (uint32, error){
	buf := make([]byte, 4)
	n, err := conn.Read(buf)
	if err != nil{
		return 0, err
	}

	if n != 4 {
		return 0, fmt.Errorf("ReadInt32: expected 4 bytes but got %d", n)
	}
	return binary.BigEndian.Uint32(buf), nil
}

func WriteInt32(conn net.Conn, n uint32) error{
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, n)
	_, err := conn.Write(buf)
	return err
}

func ReadString(conn net.Conn) (string, error){
	var bytes []byte
	b := make([]byte, 1)
	for {
		_, err := conn.Read(b)
		if err != nil {
			return "", err
		}
		if b[0] == 0x00 {
			return string(bytes), nil
		}
		bytes = append(bytes, b[0])
	}
}

func WriteString(conn net.Conn, s string) error {
	data := append([]byte(s), 0x00)
	_, err := conn.Write(data)
	return err
}

const (
	MsgAuthOk        byte = 'R'
	MsgBackendKey    byte = 'K'
	MsgParamStatus   byte = 'S'
	MsgReadyForQuery byte = 'Z'
	MsgQuery         byte = 'Q'
	MsgEmptyQuery    byte = 'I'
	MsgTerminate     byte = 'X'
	MsgRowDesc       byte = 'T'
	MsgDataRow       byte = 'D'
	MsgCmdComplete   byte = 'C'
	MsgError         byte = 'E'
)

