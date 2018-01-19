package network

import (
	"encoding/binary"
	"bufio"
	"bytes"
)

var length int32 = 0

type Package struct {
	BufferSize int
	HeaderSize int
}

func (p *Package) Package(message string) ([]byte, error) {
	var length int32 = int32(len(message)) + 5
	var pkg *bytes.Buffer = new(bytes.Buffer)
	binary.Write(pkg, binary.BigEndian, []byte{0x0a})
	binary.Write(pkg, binary.BigEndian, length)
	err := binary.Write(pkg, binary.BigEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

func (p *Package) UnPackage(reader *bufio.Reader) ([]byte, int32, error) {

	// TODO :: if make []byte every time right
	var buffer = make([]byte, p.BufferSize)

	// header
	header, err := reader.Peek(p.HeaderSize)
	if err != nil {
		return nil, 0, err
	}

	// reader.Buffered()
	if length == 0 && header[0] == 0x0a {
		err := binary.Read(bytes.NewBuffer(header[1:5]), binary.BigEndian, &length)
		if err != nil {
			return nil, 0, err
		}
		reader.Discard(p.HeaderSize)
		length -= int32(p.HeaderSize)
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
