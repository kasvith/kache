package resp2

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/kasvith/kache/internal/protocol"
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(r *bufio.Reader) *Parser {
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
		arrLen, err := p.readArrayLength()
		if err != nil {
			return nil, err
		}

		if arrLen == 0 {
			return &protocol.Command{}, nil
		}

		if arrLen == -1 {
			return nil, nil
		}

		args := make([]string, arrLen)
		for i := 0; i < arrLen; i++ {
			str, err := p.readBulkString()
			if err != nil {
				return nil, err
			}
			args[i] = str
		}
		return &protocol.Command{Name: strings.ToLower(args[0]), Args: args[1:]}, nil

	default:
		return nil, &protocol.ErrUnknownProtocol{}
	}

	return nil, nil
}

func (p Parser) readArrayLength() (int, error) {
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
		return 0, &protocol.ErrCastFailedToInt{Val: string(bs)}
	}

	return val, nil
}

func (p Parser) readBulkString() (string, error) {
	b, err := p.reader.ReadByte()

	if err != nil {
		return "", err
	}

	if b != TypeBulkString {
		return "", &protocol.ErrWrongType{}
	}

	lenBuf, err := p.reader.ReadBytes(LF)
	if err != nil {
		return "", err
	}

	bs, err := trimCRLF(lenBuf)
	if err != nil {
		return "", err
	}

	strLen, err := strconv.Atoi(string(bs))
	if err != nil {
		return "", &protocol.ErrCastFailedToInt{Val: string(bs)}
	}

	buf := make([]byte, strLen)
	n, err := io.ReadFull(p.reader, buf)
	if err != nil || n < strLen {
		return "", &protocol.ErrUnexpectedLineEnd{}
	}

	// eat CR
	b, err = p.reader.ReadByte()
	if b != CR {
		return "", &protocol.ErrUnexpectedLineEnd{}
	}

	// eat LF
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
