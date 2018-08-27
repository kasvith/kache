/*
 * MIT License
 *
 * Copyright (c)  2018 Kasun Vithanage
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/kasvith/kache/internal/protcl"
)

// RunCli start kache-cli command
func RunCli(host string, port int) {
	if err := Dial(fmt.Sprintf("%s:%d", host, port)); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := prompt.New(
		Executor,
		Completer,
		prompt.OptionPrefix(fmt.Sprintf("%s:%d> ", host, port)),
		prompt.OptionTitle("kache-cli"),
	)
	p.Run()
}

// Executor used in CLI
func Executor(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	} else if s == "quit" || s == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}

	if err := c.Write(protcl.NewSliceResp3(strings.Split(s, " "))); err != nil {
		fmt.Println(err)
		return
	}

	resp, err := c.resp3Parser.Parse()
	if err != nil {
		fmt.Println(err)
		return
	} else if resp != nil {
		fmt.Println(resp.RenderString())
		return
	}
	fmt.Println("(empty)")
}

// Completer used in CLI
func Completer(document prompt.Document) []prompt.Suggest {
	return nil
}
