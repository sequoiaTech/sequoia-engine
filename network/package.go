package network

import (
	"encoding/binary"
	"bufio"
	"bytes"
)

var length int32 = 0
var headerSize int = 5

type Package struct {
	BufferSize int
}

func (P *Package) Package(message string) ([]byte, error) {
	var length int32 = int32(len(message)) + int32(headerSize)
	var pkg *bytes.Buffer = new(bytes.Buffer)
	binary.Write(pkg, binary.BigEndian, []byte{0x0a})
	binary.Write(pkg, binary.BigEndian, length)
	err := binary.Write(pkg, binary.BigEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

func (P *Package) UnPackage(reader *bufio.Reader) ([]byte, int32, error) {

	// TODO :: if make []byte every time right
	var buffer = make([]byte, P.BufferSize)

	// header
	header, err := reader.Peek(headerSize)
	if err != nil {
		return nil, 0, err
	}

	// reader.Buffered()
	if length == 0 && header[0] == 0x0a {
		err := binary.Read(bytes.NewBuffer(header[1:5]), binary.BigEndian, &length)
		if err != nil {
			return nil, 0, err
		}
		reader.Discard(headerSize)
		length -= int32(headerSize)
	}

	// read
	len, err := reader.Read(buffer)
	if err != nil {
		return nil, 0, err
	}
	length -= int32(len)
	if length < 0 {
		return nil, 0, err
	}

	return buffer[:len], length, nil
}
