package protocol

import (
	"fmt"
	"log"
	"net"
)

const ProtocolVersion3 uint32 = 196608

func HandleStartup(conn net.Conn) error {

	length, err := ReadInt32(conn)
	if err != nil {
		return fmt.Errorf("read startup length: %w", err)
	}
	log.Printf("startup message length: %d", length)

	version, err := ReadInt32(conn)
	if err != nil {
		return fmt.Errorf("read protocol version: %w", err)
	}
	log.Printf("protocol version: %d", version)

	if version != ProtocolVersion3 {
		return fmt.Errorf("unsupported protocol version: %d (want %d)", version, ProtocolVersion3)
	}

	params := make(map[string]string)
	for {
		key, err := ReadString(conn)
		if err != nil {
			return fmt.Errorf("read param key: %w", err)
		}
		if key == "" {
			break
		}
		value, err := ReadString(conn)
		if err != nil {
			return fmt.Errorf("read param value: %w", err)
		}
		params[key] = value
		log.Printf("startup param: %s = %s", key, value)
	}

	// ----- AuthenticationOk -----
	// type 'R' + length 8 + auth type 0
	conn.Write([]byte{MsgAuthOk})
	WriteInt32(conn, 8)
	WriteInt32(conn, 0)

	// ----- ParameterStatus -----
	sendParamStatus(conn, "server_version", "14.0")
	sendParamStatus(conn, "client_encoding", "UTF8")
	sendParamStatus(conn, "server_encoding", "UTF8")

	// ----- BackendKeyData -----
	// type 'K' + length 12 + pid 4 bytes + secret 4 bytes
	conn.Write([]byte{MsgBackendKey})
	WriteInt32(conn, 12)
	WriteInt32(conn, 8080)
	WriteInt32(conn, 12345)

	// ----- ReadyForQuery -----
	// type 'Z' + length 5 + idle byte 'I'
	conn.Write([]byte{MsgReadyForQuery})
	WriteInt32(conn, 5)
	conn.Write([]byte{'I'})


	return nil
}

func sendParamStatus(conn net.Conn, key, value string) {
	// length = 4 (itself) + len(key) + 1 (\0) + len(value) + 1 (\0)
	totalLen := uint32(4 + len(key) + 1 + len(value) + 1)

	conn.Write([]byte{MsgParamStatus})
	WriteInt32(conn, totalLen)
	WriteString(conn, key)
	WriteString(conn, value)
}