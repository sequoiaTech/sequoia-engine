package network

import (
	"net"
	"bufio"
	"encoding/binary"
	"bytes"
)

var headerSize = 2

type Message struct {
	BufferSize int
	Payload    []byte
}

func (M *Message) Read(conn *net.Conn) ([]byte, error) {

	var messageSize uint16 = 0

	defer func() {
		M.Payload = []byte{}
	}()

	reader := bufio.NewReaderSize(*conn, M.BufferSize)

	for {
		var buffer = make([]byte, M.BufferSize)

		if messageSize == 0 {
			header, err := reader.Peek(headerSize)
			if err != nil {
				return nil, err
			}

			err = binary.Read(bytes.NewBuffer(header[0:2]), binary.BigEndian, &messageSize)
			if err != nil {
				return nil, err
			}

			reader.Discard(headerSize)
			messageSize -= uint16(headerSize)
		}

		len, err := reader.Read(buffer)
		if err != nil {
			return nil, err
		}

		messageSize -= uint16(len)
		if messageSize < 0 {
			return nil, err
		}

		M.Payload = append(M.Payload, buffer[:len] ...)

		if messageSize == 0 {
			break
		}
	}

	return M.Payload, nil
}

func (M *Message) Write(conn *net.Conn, msg []byte) error {
	var pkg *bytes.Buffer = new(bytes.Buffer)
	var messageSize uint16 = uint16(len(msg)) + uint16(headerSize)

	binary.Write(pkg, binary.BigEndian, messageSize)
	err := binary.Write(pkg, binary.BigEndian, []byte(msg))
	if err != nil {
		return err
	}

	_, err = (*conn).Write(pkg.Bytes())
	if err != nil {
		return err
	}

	return nil
}
