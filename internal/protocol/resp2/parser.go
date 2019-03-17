package resp2

import (
	"bufio"
	"github.com/kasvith/kache/internal/protocol"
	"strconv"
)

type Parser struct {
	reader bufio.Reader
}

func NewParser(r bufio.Reader) *Parser {
	return &Parser{r}
}

// Parse reads commands as bulk strings
func (p Parser) Parse() (*Command, error) {
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
		}

	default:
	}

	return nil, nil
}

func (p Parser)readIntBeforeCRLF() (int, error) {
	buf, err := p.reader.ReadBytes('\n')
	if err != nil {
		return 0, err
	}

	buf, err = trimCRLF(buf)
	if err != nil {
		return 0, err
	}

	val, err := strconv.Atoi(string(buf))
	if err != nil {
		return 0, err
	}

	return val, nil
}

func trimCRLF(buf []byte) ([]byte,error) {
	bufLen := len(buf)

	if bufLen == 0 || bufLen <= 2 || buf[bufLen - 1] != '\n' || buf[bufLen - 2] != '\r' {
		return nil, &protocol.ErrUnexpectedLineEnd{}
	}

	return buf[:bufLen-2], nil
}
