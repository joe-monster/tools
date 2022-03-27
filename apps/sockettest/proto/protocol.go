package proto

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
)

//goim的tcp协议包struct
type protocol struct {
	Length       uint32
	HeaderLength uint16
	Version      uint16
	Act          uint32
	Seq          uint32

	Body []byte
}

func Encode(body []byte) ([]byte, error) {

	p := protocol{
		Length:       16 + uint32(len(body)),
		HeaderLength: 16,
		Version:      1,
		Act:          2,
		Seq:          3,
		Body:         body,
	}

	data := new(bytes.Buffer)

	var err error
	err = binary.Write(data, binary.BigEndian, &p.Length)
	err = binary.Write(data, binary.BigEndian, &p.HeaderLength)
	err = binary.Write(data, binary.BigEndian, &p.Version)
	err = binary.Write(data, binary.BigEndian, &p.Act)
	err = binary.Write(data, binary.BigEndian, &p.Seq)
	err = binary.Write(data, binary.BigEndian, &p.Body)

	return data.Bytes(), err
}

func Decode(data []byte) (*protocol, error) {
	p := protocol{}

	lengthBuffer := bytes.NewBuffer(data[:4])
	if err := binary.Read(lengthBuffer, binary.BigEndian, &p.Length); err != nil {
		return nil, errors.WithStack(err)
	}

	headerLengthBuffer := bytes.NewBuffer(data[4:6])
	if err := binary.Read(headerLengthBuffer, binary.BigEndian, &p.HeaderLength); err != nil {
		return nil, errors.WithStack(err)
	}

	versionBuffer := bytes.NewBuffer(data[6:8])
	if err := binary.Read(versionBuffer, binary.BigEndian, &p.Version); err != nil {
		return nil, errors.WithStack(err)
	}

	actBuffer := bytes.NewBuffer(data[8:12])
	if err := binary.Read(actBuffer, binary.BigEndian, &p.Act); err != nil {
		return nil, errors.WithStack(err)
	}

	seqBuffer := bytes.NewBuffer(data[12:16])
	if err := binary.Read(seqBuffer, binary.BigEndian, &p.Seq); err != nil {
		return nil, errors.WithStack(err)
	}

	p.Body = data[16:]

	return &p, nil
}
