package protcl

import (
	"bufio"
	"math/big"
	"strconv"
	"strings"
)

// Resp3Parser is for parser resp3 protocol
type Resp3Parser struct {
	reader *bufio.Reader
}

// NewResp3Parser return a Resp3Parser
func NewResp3Parser(r *bufio.Reader) *Resp3Parser {
	return &Resp3Parser{reader: r}
}

// Commands parse resp3 message to kache command
func (r *Resp3Parser) Commands() (*RespCommand, error) {
	resp3, err := r.Parse()
	if err != nil {
		return nil, err
	}
	args, err := resp3.commands()
	if err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, &ErrInvalidCommand{}
	}

	return &RespCommand{Name: strings.ToLower(args[0]), Args: args[1:]}, nil
}

// Parse return Resp3
func (r *Resp3Parser) Parse() (*Resp3, error) {
	resp, err := r.parse()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *Resp3Parser) parse() (*Resp3, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch b {
	case Resp3SimpleString, Resp3SimpleError:
		str, err := r.stringBeforeLF()
		if err != nil {
			return nil, err
		}
		return &Resp3{Type: b, Str: str}, nil
	case Resp3BlobString, Resp3BolbError:
		length, err := r.intBeforeLF()
		if err != nil {
			return nil, err
		}

		bs, err := r.readLengthBytesWithLF(length)
		if err != nil {
			return nil, err
		}

		return &Resp3{Type: b, Str: string(bs)}, nil
	case Resp3Number:
		integer, err := r.intBeforeLF()
		if err != nil {
			return nil, err
		}
		return &Resp3{Type: b, Integer: integer}, nil
	case Resp3Double:
		str, err := r.stringBeforeLF()
		if err != nil {
			return nil, err
		}
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, &ErrConvertType{Type: "double", Value: str, Err: err}
		}
		return &Resp3{Type: b, Double: f}, nil
	case Resp3BigNumber:
		str, err := r.stringBeforeLF()
		if err != nil {
			return nil, err
		}
		bigInt, ok := big.NewInt(0).SetString(str, 10)
		if !ok {
			return nil, &ErrConvertType{Type: "Big Number", Value: str}
		}
		return &Resp3{Type: b, BigInt: bigInt}, nil
	case Resp3Null:
		if _, err := r.readLengthBytesWithLF(0); err != nil {
			return nil, err
		}
		return &Resp3{Type: b}, nil
	case Resp3Boolean:
		buf, err := r.readLengthBytesWithLF(1)
		if err != nil {
			return nil, err
		}

		switch buf[0] {
		case 't':
			return &Resp3{Type: b, Boolean: true}, nil
		case 'f':
			return &Resp3{Type: b, Boolean: false}, nil
		}
		return nil, &ErrUnexpectString{Str: "t/f"}
	case Resp3Array, Resp3Set:
		length, err := r.intBeforeLF()
		if err != nil {
			return nil, err
		}
		resp := &Resp3{Type: b}
		for i := 0; i < length; i++ {
			elem, err := r.Parse()
			if err != nil {
				return nil, err
			}
			resp.Elems = append(resp.Elems, elem)
		}
		return resp, nil
	}

	return nil, &ErrProtocolType{Type: b}
}

func (r *Resp3Parser) stringBeforeLF() (string, error) {
	buf, err := r.reader.ReadBytes(LF)
	if err != nil {
		return "", err
	}
	bs, err := trimLastLF(buf)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (r *Resp3Parser) intBeforeLF() (int, error) {
	buf, err := r.reader.ReadBytes(LF)
	if err != nil {
		return 0, err
	}
	bs, err := trimLastLF(buf)
	if err != nil {
		return 0, err
	}
	s := string(bs)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, &ErrCastFailedToInt{Val: s}
	}
	return i, nil
}

func (r *Resp3Parser) readLengthBytesWithLF(length int) ([]byte, error) {
	if length == 0 {
		if b, err := r.reader.ReadByte(); err != nil {
			return nil, err
		} else if b != LF {
			return nil, &ErrUnexpectString{Str: "<LF>"}
		}
		return nil, nil
	}

	buf := make([]byte, length+1)
	n, err := r.reader.Read(buf)
	if err != nil {
		return nil, err
	} else if n < length+1 {
		return nil, &ErrUnexpectedLineEnd{}
	}

	return trimLastLF(buf)
}

func trimLastLF(buf []byte) ([]byte, error) {
	bufLen := len(buf)
	if len(buf) == 0 || buf[bufLen-1] != LF {
		return nil, &ErrUnexpectedLineEnd{}
	}

	return buf[:bufLen-1], nil
}
