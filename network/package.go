package network

import (
	"encoding/binary"
	"bufio"
	"bytes"
)

var (
	length     uint16 = 0
	headerSize int    = 2
)

type Package struct {
	BufferSize int
}

func (P *Package) Package(message string) ([]byte, error) {
	var length uint16 = uint16(len(message)) + uint16(headerSize)
	var pkg *bytes.Buffer = new(bytes.Buffer)
	binary.Write(pkg, binary.BigEndian, length)
	err := binary.Write(pkg, binary.BigEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

func (P *Package) UnPackage(reader *bufio.Reader) ([]byte, uint16, error) {

	// TODO :: if make []byte every time right
	var buffer = make([]byte, P.BufferSize)

	if length == 0 {
		header, err := reader.Peek(headerSize)
		if err != nil {
			return nil, 0, err
		}

		err = binary.Read(bytes.NewBuffer(header[0:headerSize]), binary.BigEndian, &length)
		if err != nil {
			return nil, 0, err
		}

		reader.Discard(headerSize)
		length -= uint16(headerSize)
	}

	// read
	len, err := reader.Read(buffer)
	if err != nil {
		return nil, 0, err
	}
	length -= uint16(len)
	if length < 0 {
		return nil, 0, err
	}

	return buffer[:len], length, nil
}
