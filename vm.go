package main

import (
	"log"
	"os"
)

const (
	OINC = iota
	OSHIFT
	OINPUT
	OOUTPUT
	OJMP
	OCLEAR
	OMULADD
)

type Instruction struct {
	op, arg1 uint8
	arg2     int
}

// ooo11111 11122222 22222222 22222222
//  op     arg1                   arg2
type PackedInstruction uint32

type Program []PackedInstruction

func (inst Instruction) Pack() PackedInstruction {
	return PackedInstruction((uint32(inst.op)&0x7)<<29 | (uint32(inst.arg1)&0xff)<<21 | uint32(inst.arg2)&0x1fffff)
}

func (packed PackedInstruction) Unpack() (inst Instruction) {
	inst.arg2 = int(packed & 0xfffff)
	if packed&0x100000 != 0 {
		inst.arg2 |= ^0xfffff
	}
	packed >>= 21
	inst.arg1 = uint8(packed & 0xff)
	packed >>= 8
	inst.op = uint8(packed)
	return
}

func (p *Program) Append(inst Instruction) {
	*p = append(*p, inst.Pack())
}

func (p *Program) Extend(p1 Program) {
	*p = append(*p, p1...)
}

func (p *Program) Len() int {
	return len(*p)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func VMRun(p Program) {
	var (
		memory [65536]uint8
		cursor uint16
		pc     int
	)
	for pc = 0; pc < len(p); pc++ {
		switch inst := p[pc].Unpack(); inst.op {
		case OINC:
			memory[cursor] += inst.arg1
		case OSHIFT:
			cursor += uint16(inst.arg2)
		case OINPUT:
			_, err := os.Stdin.Read(memory[cursor : cursor+1])
			checkError(err)
		case OOUTPUT:
			_, err := os.Stdout.Write(memory[cursor : cursor+1])
			checkError(err)
		case OJMP:
			if (memory[cursor] == 0) == (inst.arg1 == 0) {
				// should jump
				pc += inst.arg2
			}
		case OCLEAR:
			memory[cursor] = 0
		case OMULADD:
			memory[cursor+uint16(inst.arg2)] += inst.arg1 * memory[cursor]
		default:
			panic("unreached")
		}
	}
}
