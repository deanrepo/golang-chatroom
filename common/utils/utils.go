// Package utils provide utils for the program
package utils

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"net"
	"projects/chatroom1/common/message"
)

// WritePkg writes message package to the remote peer of the connection.
func WritePkg(conn net.Conn, data []byte) (err error) {
	// Send the length of the data
	var msgLen uint32
	msgLen = uint32(len(data))
	bufLen := msgLen + 4
	buf := make([]byte, bufLen)
	binary.BigEndian.PutUint32(buf[0:4], msgLen)

	n := copy(buf[4:], data)
	if n != int(msgLen) {
		err = errors.New("do not copy all data to buf")
		log.Printf("WritePkg -> copy err: %v\n", err)
		return
	}

	// Send the pkg.
	n, err = conn.Write(buf)
	if n != int(bufLen) || err != nil {
		log.Printf("WritePkg -> Write(buf) fail, err: %v\n", err)
		return
	}

	return
}

// ReadPkg reads message package from remote peer of the connection.
func ReadPkg(conn net.Conn) (msg *message.Message, err error) {

	buf := make([]byte, 8096)

	_, err = conn.Read(buf[:4])
	if err != nil {
		log.Printf("ReadPkg -> Read message len err: %v\n", err)
		return
	}

	// Get the length of the message
	msgLen := binary.BigEndian.Uint32(buf[:4])

	// Get the message
	n, err := conn.Read(buf[:msgLen])
	if err != nil || n != int(msgLen) {
		log.Printf("ReadPkg -> Read message err: %v\n", err)
		return
	}

	// initialize msg
	msg = &message.Message{}
	// De-serialize the message
	err = json.Unmarshal(buf[:msgLen], msg)
	if err != nil {
		log.Printf("ReadPkg -> Unmarshal Message err: %v\n", err)
		return
	}

	return
}
