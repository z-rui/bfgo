package main

func (p *Program) Compile(ast interface{}) {
	switch inst := ast.(type) {
	case IncInst:
		p.Append(Instruction{op: OINC, arg1: uint8(inst)})
	case ShiftInst:
		p.Append(Instruction{op: OSHIFT, arg2: int(inst)})
	case InInst:
		p.Append(Instruction{op: OINPUT})
	case OutInst:
		p.Append(Instruction{op: OOUTPUT})
	case Block:
		for _, stmt := range inst {
			p.Compile(stmt)
		}
	case Loop:
		if canDiscardLoop(*p) {
			break
		}
		if loop := specialLoop(inst); loop != nil {
			p.Extend(loop)
		} else {
			loop.Compile(Block(inst))
			offset := len(loop) + 1
			p.Append(Instruction{op: OJMP, arg1: 0, arg2: offset})
			p.Extend(loop)
			p.Append(Instruction{op: OJMP, arg1: 1, arg2: -offset})
		}
	default:
		panic("unreached")
	}
	return
}

func canDiscardLoop(p Program) bool {
	if len(p) == 0 {
		return false
	}
	switch inst := p[len(p)-1].Unpack(); inst.op {
	case OCLEAR:
		return true
	case OJMP:
		return inst.arg1 != 0 // JNZ
	default:
		return false
	}
}

func specialLoop(l Loop) Program {
	if len(l) == 0 {
		return nil
	}
	if delta, ok := l[len(l)-1].(IncInst); ok && delta == -1 {
		var loop Program
		offsetSum := 0
		for i := 0; i+1 < len(l); i += 2 {
			if offset, ok := l[i].(ShiftInst); ok {
				offsetSum += int(offset)
				if i+2 == len(l) {
					break
				}
				ratio, ok := l[i+1].(IncInst)
				if !ok || offsetSum == 0 {
					return nil
				}
				loop.Append(Instruction{op: OMULADD, arg1: uint8(ratio), arg2: offsetSum})
			} else {
				return nil
			}
		}
		if offsetSum != 0 {
			return nil
		}
		loop.Append(Instruction{op: OCLEAR})
		return loop
	}
	return nil
}
