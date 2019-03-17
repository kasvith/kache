package resp2

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/kasvith/kache/internal/protocol"
)

type Parser struct {
	reader bufio.Reader
}

func NewParser(r bufio.Reader) *Parser {
	return &Parser{r}
}

// Parse reads commands as bulk strings
func (p Parser) Parse() (*protocol.Command, error) {
	r := p.reader

	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch b {
	case TypeArray:
		arrLen, err := p.readIntBeforeCRLF()
		if err != nil {
			return nil, err
		}

		args := make([]string, arrLen)
		for i := 0; i < arrLen; i++ {
			b, err = r.ReadByte()
			if err != nil {
				return nil, err
			}

			if b != TypeBulkString {
				return nil, &protocol.ErrWrongType{}
			}

			strLen, err := p.readIntBeforeCRLF()
			if err != nil {
				return nil, err
			}

			str, err := p.readBulkString(strLen)
			if err != nil {
				return nil, err
			}

			args = append(args, str)
		}

		return &protocol.Command{Name: strings.ToLower(args[0]), Args: args[1:]}, nil

	default:
		return nil, &protocol.ErrUnknownProtocol{}
	}

	return nil, nil
}

func (p Parser) readIntBeforeCRLF() (int, error) {
	buf, err := p.reader.ReadBytes(LF)
	if err != nil {
		return 0, err
	}

	bs, err := trimCRLF(buf)
	if err != nil {
		return 0, err
	}

	val, err := strconv.Atoi(string(bs))
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (p Parser) readBulkString(strLen int) (string, error) {
	buf := make([]byte, strLen)
	n, err := io.ReadFull(&p.reader, buf)
	if err != nil || n < strLen {
		return "", &protocol.ErrUnexpectedLineEnd{}
	}

	// eat CRLF
	b, err := p.reader.ReadByte()
	if b != CR {
		return "", &protocol.ErrUnexpectedLineEnd{}
	}

	b, err = p.reader.ReadByte()
	if b != LF {
		return "", &protocol.ErrUnexpectedLineEnd{}
	}

	return string(buf), nil
}

func trimCRLF(buf []byte) ([]byte, error) {
	bufLen := len(buf)

	if bufLen == 0 || bufLen <= 2 || buf[bufLen-1] != '\n' || buf[bufLen-2] != '\r' {
		return nil, &protocol.ErrUnexpectedLineEnd{}
	}

	return buf[:bufLen-2], nil
}
