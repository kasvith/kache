package protcl

import (
	"bufio"
	"math/big"
	"strings"
	"testing"

	testifyAssert "github.com/stretchr/testify/assert"
)

func testResp3Parser(t *testing.T, input, expect string) {
	assert := testifyAssert.New(t)

	parser := NewResp3Parser(bufio.NewReader(strings.NewReader(input)))
	result, err := parser.Parse()
	assert.Nil(err)
	assert.Equal(expect, result.RenderString())
}

func testResp3Error(t *testing.T, input, e string) {
	assert := testifyAssert.New(t)

	parser := NewResp3Parser(bufio.NewReader(strings.NewReader(input)))
	result, err := parser.Parse()
	assert.NotNil(err)
	assert.Nil(result)
	assert.Equal(e, err.Error())
}

type TestResp3 struct {
	resp3    *Resp3
	protocol string
	render   string
	err      string
}

func TestResp3Parser(t *testing.T) {
	assert := testifyAssert.New(t)

	bigNumber, _ := big.NewInt(0).SetString("3492890328409238509324850943850943825024385", 10)
	testCases := []TestResp3{
		// simple renderString
		{protocol: "+renderString", err: "EOF"},
		{resp3: &Resp3{Type: Resp3SimpleString, Str: ""}, protocol: "+\n", render: `""`},
		{resp3: &Resp3{Type: Resp3SimpleString, Str: "hello"}, protocol: "+hello\n", render: `"hello"`},
		{resp3: &Resp3{Type: Resp3SimpleString, Str: "hello world"}, protocol: "+hello world\n", render: `"hello world"`},

		// blob renderString
		{protocol: "$1\n\n", err: "unexpected line end"},
		{protocol: "$1\naa\n", err: "unexpected line end"},
		{resp3: &Resp3{Type: Resp3BlobString, Str: ""}, protocol: "$0\n\n", render: `""`},
		{resp3: &Resp3{Type: Resp3BlobString, Str: "hello"}, protocol: "$5\nhello\n", render: `"hello"`},
		{resp3: &Resp3{Type: Resp3BlobString, Str: "hello\nworld"}, protocol: "$11\nhello\nworld\n", render: "\"hello\\nworld\""},

		// simple error
		{protocol: "-renderString", err: "EOF"},
		{resp3: &Resp3{Type: Resp3SimpleError, Str: ""}, protocol: "-\n", render: `(error) `},
		{resp3: &Resp3{Type: Resp3SimpleError, Str: "hello"}, protocol: "-hello\n", render: `(error) hello`},
		{resp3: &Resp3{Type: Resp3SimpleError, Str: "hello world"}, protocol: "-hello world\n", render: `(error) hello world`},

		// blob error
		{protocol: "!1\n\n", err: "unexpected line end"},
		{protocol: "!1\naa\n", err: "unexpected line end"},
		{resp3: &Resp3{Type: Resp3BolbError, Str: ""}, protocol: "!0\n\n", render: `(error) `},
		{resp3: &Resp3{Type: Resp3BolbError, Str: "hello"}, protocol: "!5\nhello\n", render: `(error) hello`},
		{resp3: &Resp3{Type: Resp3BolbError, Str: "hello\nworld"}, protocol: "!11\nhello\nworld\n", render: "(error) hello\nworld"},

		// number
		{protocol: ":invalid", err: "EOF"},
		{protocol: ":invalid\n", err: "ERR: error casting invalid to int"},
		{resp3: &Resp3{Type: Resp3Number, Integer: -1}, protocol: ":-1\n", render: `(integer) -1`},
		{resp3: &Resp3{Type: Resp3Number, Integer: 0}, protocol: ":0\n", render: `(integer) 0`},
		{resp3: &Resp3{Type: Resp3Number, Integer: 100}, protocol: ":100\n", render: `(integer) 100`},

		// double
		{protocol: ",invalid", err: "EOF"},
		{protocol: ",invalid\n", err: "convert invalid to double fail, because of strconv.ParseFloat: parsing \"invalid\": invalid syntax"},
		{resp3: &Resp3{Type: Resp3Double, Double: -1}, protocol: ",-1\n", render: "(double) -1"},
		{resp3: &Resp3{Type: Resp3Double, Double: 0}, protocol: ",0\n", render: "(double) 0"},
		{resp3: &Resp3{Type: Resp3Double, Double: 10}, protocol: ",10\n", render: "(double) 10"},
		{resp3: &Resp3{Type: Resp3Double, Double: 1.23}, protocol: ",1.23\n", render: "(double) 1.23"},
		{protocol: ",.1\n", render: "(double) 0.1"},
		{protocol: ",1.\n", render: "(double) 1"},

		// big number
		{protocol: "(invalid", err: "EOF"},
		{protocol: "(invalid\n", err: "convert invalid to Big Number fail"},
		{resp3: &Resp3{Type: Resp3BigNumber, BigInt: bigNumber}, protocol: "(3492890328409238509324850943850943825024385\n", render: "(big number) 3492890328409238509324850943850943825024385"},

		// null
		{protocol: "_invalid", err: "unexpect string: <LF>"},
		{protocol: "_invalid\n", err: "unexpect string: <LF>"},
		{resp3: &Resp3{Type: Resp3Null}, protocol: "_\n", render: "(null)"},

		// boolean
		{protocol: "#", err: "EOF"},
		{protocol: "#\n", err: "unexpected line end"},
		{protocol: "#x\n", err: "unexpect string: t/f"},
		{protocol: "#invalid", err: "unexpected line end"},
		{protocol: "#invalid\n", err: "unexpected line end"},
		{resp3: &Resp3{Type: Resp3Boolean, Boolean: true}, protocol: "#t\n", render: "(boolean) true"},
		{resp3: &Resp3{Type: Resp3Boolean, Boolean: false}, protocol: "#f\n", render: "(boolean) false"},

		// array
		{protocol: "*", err: "EOF"},
		{protocol: "*\n", err: "ERR: error casting  to int"},
		{protocol: "*invalid", err: "EOF"},
		{protocol: "*invalid\n", err: "ERR: error casting invalid to int"},
		{protocol: "*1\n\n", err: "unknown protocol type: \n"},
		{protocol: "*1\ninvalid\n", err: "unknown protocol type: i"},
		{protocol: "*3\n:1\n:2\n", err: "EOF"},
		{resp3: &Resp3{Type: Resp3Array, Elems: []*Resp3{
			{Type: Resp3Number, Integer: 1},
			{Type: Resp3Number, Integer: 2},
			{Type: Resp3Number, Integer: 3},
		}}, protocol: "*3\n:1\n:2\n:3\n", render: "(array)\n\t(integer) 1\n\t(integer) 2\n\t(integer) 3"},
		{resp3: &Resp3{Type: Resp3Array, Elems: []*Resp3{
			{Type: Resp3Array, Elems: []*Resp3{
				{Type: Resp3Number, Integer: 1},
				{Type: Resp3BlobString, Str: "hello"},
				{Type: Resp3Number, Integer: 2},
			}},
			{Type: Resp3Boolean, Boolean: false},
		}}, protocol: "*2\n*3\n:1\n$5\nhello\n:2\n#f\n", render: "(array)\n\t(array)\n\t\t(integer) 1\n\t\t\"hello\"\n\t\t(integer) 2\n\t(boolean) false"},

		// set
		{protocol: "~", err: "EOF"},
		{protocol: "~\n", err: "ERR: error casting  to int"},
		{protocol: "~invalid", err: "EOF"},
		{protocol: "~invalid\n", err: "ERR: error casting invalid to int"},
		{protocol: "~1\n\n", err: "unknown protocol type: \n"},
		{protocol: "~1\ninvalid\n", err: "unknown protocol type: i"},
		{protocol: "~3\n:1\n:2\n", err: "EOF"},
		{resp3: &Resp3{Type: Resp3Set, Elems: []*Resp3{
			{Type: Resp3Number, Integer: 1},
			{Type: Resp3Number, Integer: 2},
			{Type: Resp3Number, Integer: 3},
		}}, protocol: "~3\n:1\n:2\n:3\n", render: "(set)\n\t(integer) 1\n\t(integer) 2\n\t(integer) 3"},
		{resp3: &Resp3{Type: Resp3Set, Elems: []*Resp3{
			{Type: Resp3Array, Elems: []*Resp3{
				{Type: Resp3Number, Integer: 1},
				{Type: Resp3BlobString, Str: "hello"},
				{Type: Resp3Number, Integer: 2},
			}},
			{Type: Resp3Boolean, Boolean: false},
		}}, protocol: "~2\n*3\n:1\n$5\nhello\n:2\n#f\n", render: "(set)\n\t(array)\n\t\t(integer) 1\n\t\t\"hello\"\n\t\t(integer) 2\n\t(boolean) false"},
	}

	for _, testCase := range testCases {
		if testCase.resp3 == nil {
			if testCase.render != "" {
				// protocol renderString -> render renderString
				testResp3Parser(t, testCase.protocol, testCase.render)
			} else {
				// err test
				testResp3Error(t, testCase.protocol, testCase.err)
			}
		} else {
			// resp3 -> protocol renderString -> render renderString
			assert.Equal(testCase.protocol, testCase.resp3.ProtocolString())
			testResp3Parser(t, testCase.protocol, testCase.render)
		}
	}
}
