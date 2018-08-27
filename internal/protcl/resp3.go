package protcl

import (
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

// NewSliceResp3 convert string slice to resp3 protocol raw string
func NewSliceResp3(slices []string) string {
	buf := new(strings.Builder)
	buf.WriteByte('*')
	buf.WriteString(strconv.Itoa(len(slices)))
	buf.WriteByte('\n')
	for _, v := range slices {
		buf.WriteByte('$')
		buf.WriteString(strconv.Itoa(len(v)))
		buf.WriteByte('\n')
		buf.WriteString(v)
		buf.WriteByte('\n')
	}
	return buf.String()
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

func (r *Resp3) commands() ([]string, error) {
	switch r.Type {
	case RepSimpleString, Resp3BlobString, Resp3Number, Resp3Double, Resp3BigNumber, Resp3Boolean:
		c, err := r.command()
		if err != nil {
			return nil, err
		}
		return []string{c}, nil
	case Resp3Array:
		var slices = make([]string, len(r.Elems))
		i := 0
		for _, v := range r.Elems {
			if v.Type != RepSimpleString && v.Type != Resp3BlobString && v.Type != Resp3Number && v.Type != Resp3Double && v.Type != Resp3BigNumber && v.Type != Resp3Boolean {
				return nil, &ErrInvalidCommand{}
			}
			c, err := v.command()
			if err != nil {
				return nil, err
			}
			slices[i] = c
			i++
		}
		return slices, nil
	}
	return nil, &ErrInvalidCommand{}
}

func (r *Resp3) command() (string, error) {
	switch r.Type {
	case Resp3SimpleString, Resp3BlobString:
		return r.Str, nil
	case Resp3Number:
		return strconv.Itoa(r.Integer), nil
	case Resp3Double:
		return strconv.FormatFloat(r.Double, 'f', -1, 64), nil
	case Resp3BigNumber:
		return r.BigInt.String(), nil
	case Resp3Boolean:
		if r.Boolean {
			return "true", nil
		}
		return "false", nil
	}
	return "", &ErrInvalidCommand{}
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
