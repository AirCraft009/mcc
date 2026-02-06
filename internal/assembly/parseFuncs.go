package assembly

import (
	"fmt"
	"strconv"

	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"
)

func checkMalformedInstruction(parameters []string, currpc uint16) error {
	parameterlen := pkg.OffsetMap[parameters[0]]
	if int(parameterlen) > (len(parameters)) {
		return fmt.Errorf("Malformed instruction %s; at %d: \nexpecting at least: %d parameters; got %d.", parameters[0], currpc, parameterlen, len(parameters)-1)
	}
	return nil
}

func parseFormatOP(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	return currPC + 1, []byte{pkg.OpCodes[parameters[OpCLoc]]}, nil
}

func parseFormatOPRegAddr(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 3, 5)
	code[OpCLoc] = pkg.OpCodes[parameters[OpCLoc]]
	rx = pkg.RegMap[parameters[RegsLoc1]]
	ry, ok := pkg.RegMap[parameters[RegsLoc2]]
	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, !ok)
	var addr int
	if !ok {
		addr, syntax = strconv.Atoi(parameters[AddrLoc2])

	} else {
		addr, syntax = strconv.Atoi(parameters[StrLoc])
		//fmt.Printf("adding addr: %d at PC: %d\n", uint16(addr), currPC)
	}
	if syntax != nil {
		return currPC, code, syntax
	}
	hi, lo := helper.EncodeAddr(uint16(addr))

	code = append(code, hi, lo)
	currPC += uint16(len(code))
	return currPC, code, syntax
}

func parseFormatOPRegReg(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	err := checkMalformedInstruction(parameters, currPC)
	if err != nil {
		return 0, nil, err
	}
	code = make([]byte, 3)
	code[OpCLoc] = pkg.OpCodes[parameters[OpCLoc]]
	rx, ok := pkg.RegMap[parameters[RegsLoc1]]
	if !ok {
		return 0, nil, fmt.Errorf("Register %s not found in register map\n", parameters[RegsLoc1])
	}
	ry, ok = pkg.RegMap[parameters[RegsLoc2]]
	if !ok {
		return 0, nil, fmt.Errorf("Register %s not found in register map\n", parameters[RegsLoc2])
	}

	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, false)
	currPC += uint16(len(code))
	return currPC, code, nil
}

func parseFormatOPAddr(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 5)
	code[OpCLoc] = pkg.OpCodes[parameters[OpCLoc]]
	addr, syntax := strconv.Atoi(parameters[AddrLoc1])
	if syntax != nil {
		return 0, nil, fmt.Errorf("syntax error: %s\n", syntax.Error())
	}
	hi, lo := helper.EncodeAddr(uint16(addr))
	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, true)
	code[AddrOutLocHi] = hi
	code[AddrOutLocLo] = lo
	currPC += uint16(len(code))
	return currPC, code, syntax
}

// parseFormatOPLbl
// responsible for parsing
// JMP LBL & similar
func parseFormatOPLbl(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 5)
	code[OpCLoc] = pkg.OpCodes[parameters[OpCLoc]]
	/**
	Lbladdr, ok := parser.Labels[parameters[AddrLoc1]]
	if !ok {
		panic("label not found: " + parameters[AddrLoc1])
	}

	hi, lo := helper.EncodeAddr(uint16(addr))

	code[AddrOutLocHi] = hi
	code[AddrOutLocLo] = lo
	*/
	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, true)
	currPC += AddrOutLocHi
	parser.ObjFile.Relocs = append(parser.ObjFile.Relocs, pkg.RelocationEntry{
		InFileOffset: currPC,
		Lbl:          parameters[AddrLoc1],
	})
	code[AddrOutLocHi], code[AddrOutLocLo] = 0x00, 0x00
	currPC += uint16(len(code)) - AddrOutLocHi
	return currPC, code, syntax
}

func parseFormatOPRegLbl(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 5)
	code[OpCLoc] = pkg.OpCodes[parameters[OpCLoc]]
	rx = pkg.RegMap[parameters[RegsLoc1]]

	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, true)
	currPC += AddrOutLocHi
	parser.ObjFile.Relocs = append(parser.ObjFile.Relocs, pkg.RelocationEntry{
		InFileOffset: currPC,
		Lbl:          parameters[AddrLoc2],
	})
	code[AddrOutLocHi], code[AddrOutLocLo] = 0x00, 0x00
	currPC += uint16(len(code)) - AddrOutLocHi
	return currPC, code, syntax
}

func parseFormatOPReg(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 3)
	code[OpCLoc] = pkg.OpCodes[parameters[OpCLoc]]
	rx = pkg.RegMap[parameters[RegsLoc1]]
	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, false)
	currPC += uint16(len(code))
	return currPC, code, nil
}
