package protcl

import (
	"bufio"
	"fmt"
	"strconv"
)

// resp3 protocol type
const (
	Resp3SimpleString   = '+' // +<string>\n
	Resp3BlobString     = '$' // $<length>\n<bytes>\n
	Resp3VerbatimString = '=' // =<length>\n<format(3 bytes)>\n<bytes>\n
	Resp3SimpleError    = '-' // -<string>\n
	Resp3BolbError      = '!' // !<length>\n<bytes>\n
	Resp3Number         = ':' // :<number>\n
	Resp3Double         = ',' // ,<floating-point-number>\n
	Resp3BigNumber      = '(' // (<big number>\n
	Resp3Null           = '_' // _\n
	Resp3Boolean        = '#' // #t\n or #f\n
	Resp3Array          = '*' // *<elements number>\n... numelements other types ...
	Resp3Map            = '%' // %<elements number>\n... numelements other types ...
	Resp3Set            = '~' // ~<elements number>\n... numelements other types ...
)

// LF is \n
const LF = '\n'

// Resp3 the response of resp3 protocol
type Resp3 struct {
	Type    byte
	Str     string
	Integer int
	Boolean bool
}

func (r *Resp3) String() string {
	switch r.Type {
	case RepSimpleString, Resp3BlobString:
		return fmt.Sprintf("%q", r.Str)
	case Resp3SimpleError, Resp3BolbError:
		return "(error) " + r.Str
	case Resp3Number:
		return "(integer) " + strconv.Itoa(r.Integer)
	case Resp3Null:
		return "(null)"
	case Resp3Boolean:
		if r.Boolean {
			return "(boolean) true"
		}
		return "(boolean) false"
	}

	return "(error) unknown protocol type: " + string(r.Type)
}

// Resp3Parser is for parser resp3 protocol
type Resp3Parser struct {
	reader *bufio.Reader
}

// NewResp3Parser return a Resp3Parser
func NewResp3Parser(r *bufio.Reader) *Resp3Parser {
	return &Resp3Parser{reader: r}
}

// Parse return Resp3
func (r *Resp3Parser) Parse() (*Resp3, error) {
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
