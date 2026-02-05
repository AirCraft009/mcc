package assembly

import (
	"fmt"
	"strconv"

	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"
)

func checkMalformedInstruction(parameters []string, currpc uint16) {
	parameterlen := pkg.OffsetMap[parameters[0]]
	if int(parameterlen) > (len(parameters)) {
		panic(fmt.Sprintf("Malformed instruction %s; at %d: \nexpecting at least: %d parameters; got %d.", parameters[0], currpc, parameterlen, len(parameters)-1))
	}
}

func parseFormatOP(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	return currPC + 1, []byte{pkg.OpCodes[parameters[OpCLoc]]}, nil
}

func parseFormatOPRegAddr(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 3)
	code[OpCLoc] = pkg.OpCodes[parameters[OpCLoc]]
	rx = pkg.RegMap[parameters[RegsLoc1]]
	ry, ok := pkg.RegMap[parameters[RegsLoc2]]
	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, !ok)
	if !ok {
		addr, syntax := strconv.Atoi(parameters[AddrLoc2])
		if syntax != nil {
			panic("syntax error: " + syntax.Error())
		}
		hi, lo := helper.EncodeAddr(uint16(addr))
		code = append(code, hi, lo)
	} else {
		hi, lo := helper.EncodeAddr(uint16(0))
		code = append(code, hi, lo)
	}
	currPC += uint16(len(code))
	return currPC, code, syntax
}

func parseFormatOPRegReg(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	checkMalformedInstruction(parameters, currPC)
	code = make([]byte, 3)
	code[OpCLoc] = pkg.OpCodes[parameters[OpCLoc]]
	rx = pkg.RegMap[parameters[RegsLoc1]]
	ry = pkg.RegMap[parameters[RegsLoc2]]
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
		panic("syntax error: " + syntax.Error())
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
		Offset: currPC,
		Lbl:    parameters[AddrLoc1],
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
		Offset: currPC,
		Lbl:    parameters[AddrLoc2],
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
