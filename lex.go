package main

import (
	"bufio"
	"io"
	"log"
)

type lexState struct {
	r   *bufio.Reader
	ast Block
}

func newLexState(r io.Reader) *lexState {
	return &lexState{
		r: bufio.NewReader(r),
	}
}

func (ls *lexState) nextChar() int {
	for {
		ch, err := ls.r.ReadByte()
		if err == io.EOF {
			return 0
		}
		switch ch {
		case '+', '-', '<', '>', '[', ']', ',', '.':
			return int(ch)
		}
	}
}

func (ls *lexState) Lex(lval *yySymType) int {
	switch ch := ls.nextChar(); ch {
	case '+', '-':
		lval.intVal = ls.readPairedInst(ch, '+', '-')
		return '+'
	case '<', '>':
		lval.intVal = ls.readPairedInst(ch, '>', '<')
		return '>'
	case '[', ']', ',', '.', 0:
		return ch
	}
	panic("unreached")
}

func (ls *lexState) Error(s string) {
	log.Fatal(s)
}

func (ls *lexState) readPairedInst(ch, inc, dec int) (delta int) {
	for {
		switch ch {
		case inc:
			delta++
		case dec:
			delta--
		default:
			ls.r.UnreadByte()
			return
		}
		ch = ls.nextChar()
	}
}
