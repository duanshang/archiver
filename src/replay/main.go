package main

import (
	"github.com/yuin/gopher-lua"
	"gopkg.in/readline.v1"
	"strings"
)

type REPL struct {
	L       *lua.LState // the lua virtual machine
	toolbox *ToolBox
	rl      *readline.Instance
}

func (repl *REPL) init() {
	repl.L = lua.NewState()
	repl.toolbox = &ToolBox{}
	repl.toolbox.init("/data")
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	repl.rl = rl
}

func (repl *REPL) doREPL() {
	for {
		str, ok := repl.loadline()
		repl.rl.SetPrompt("> ")
		if !ok {
			break
		}
		repl.toolbox.exec(str)
	}
}

func (repl *REPL) dtor() {
	repl.rl.Close()
}

func incomplete(err error) bool {
	if strings.Index(err.Error(), "EOF") != -1 {
		return true
	}
	return false
}

func (repl *REPL) loadline() (string, bool) {
	line, err := repl.rl.Readline()
	if err != nil {
		return "", false
	}
	// try add return
	_, err = repl.L.LoadString("return " + line)
	if err == nil { // syntax ok
		return line, true
	} else { // syntax error
		return repl.multiline(line)
	}
}

func (repl *REPL) multiline(ml string) (string, bool) {
	repl.rl.SetPrompt(">> ")
	for {
		line, err := repl.rl.Readline()
		if err != nil {
			return "", false
		}
		ml = ml + "\n" + line

		_, err = repl.L.LoadString(ml)
		if err == nil { // syntax ok
			return ml, true
		} else if !incomplete(err) { // syntax error
			return ml, true
		}
	}
}

func main() {
	repl := &REPL{}
	repl.init()
	repl.doREPL()
	repl.dtor()
}
