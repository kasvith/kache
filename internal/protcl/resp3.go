package protcl

import (
	"bufio"
	"fmt"
	"math/big"
	"strconv"
	"strings"
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
	Double  float64
	BigInt  *big.Int
	Elems   []*Resp3
}

// RenderString convert resp3 to show-message on client
func (r *Resp3) RenderString() string {
	return r.renderString("")
}

// ProtocolString convert resp3 to protocol raw string
func (r *Resp3) ProtocolString() string {
	buf := new(strings.Builder)
	buf.WriteByte(r.Type)
	r.protocolString(buf)
	return buf.String()
}

func (r *Resp3) protocolString(buf *strings.Builder) {
	switch r.Type {
	case RepSimpleString, Resp3SimpleError:
		buf.WriteString(r.Str)
	case Resp3BlobString, Resp3BolbError:
		buf.WriteString(strconv.Itoa(len(r.Str)))
		buf.WriteByte('\n')
		buf.WriteString(r.Str)
	case Resp3Number:
		buf.WriteString(strconv.Itoa(r.Integer))
	case Resp3Double:
		buf.WriteString(strconv.FormatFloat(r.Double, 'f', -1, 64))
	case Resp3BigNumber:
		buf.WriteString(r.BigInt.String())
	case Resp3Null:
	case Resp3Boolean:
		if r.Boolean {
			buf.WriteByte('t')
		} else {
			buf.WriteByte('f')
		}
	case Resp3Array, Resp3Set:
		buf.WriteString(strconv.Itoa(len(r.Elems)))
		buf.WriteByte('\n')

		for _, v := range r.Elems {
			buf.WriteByte(v.Type)
			v.protocolString(buf)
		}
		return
	}

	buf.WriteByte('\n')
}

func (r *Resp3) renderString(pre string) string {
	switch r.Type {
	case RepSimpleString, Resp3BlobString:
		return fmt.Sprintf("%s%q", pre, r.Str)
	case Resp3SimpleError, Resp3BolbError:
		return pre + "(error) " + r.Str
	case Resp3Number:
		return pre + "(integer) " + strconv.Itoa(r.Integer)
	case Resp3Double:
		return pre + "(double) " + strconv.FormatFloat(r.Double, 'f', -1, 64)
	case Resp3BigNumber:
		return pre + "(big number) " + r.BigInt.String()
	case Resp3Null:
		return pre + "(null)"
	case Resp3Boolean:
		if r.Boolean {
			return pre + "(boolean) true"
		}
		return pre + "(boolean) false"
	case Resp3Array, Resp3Set:
		str := new(strings.Builder)
		str.WriteString(pre)
		if r.Type == Resp3Array {
			str.WriteString("(array)")
		} else {
			str.WriteString("(set)")
		}
		for _, elem := range r.Elems {
			str.WriteString("\n")
			str.WriteString(elem.renderString(pre + "\t"))
		}
		return str.String()
	}

	return pre + "(error) unknown protocol type: " + string(r.Type)
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
