package assembler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"
)

func checkMalformedInstruction(parameters []string, currpc uint16) {
	parameterlen := pkg.OffsetMap[strings.ToUpper(parameters[0])]
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

func formatString(parameters []string) (formatted [][]string) {
	//STRING RegUse RegAddr "STRING"
	var rx, ry string
	//rx = part
	//ry = addr
	rx = parameters[RegsLoc1]
	ry = parameters[RegsLoc2]
	inputStringParts := parameters[StrLoc:]
	if len(inputStringParts) == 0 {
		return formatted
	}
	var inputString string
	for _, part := range inputStringParts {
		inputString += part + " "
	}
	inputString = strings.ReplaceAll(inputString, "\"", "")
	inputString = inputString[:len(inputString)-1]
	length := len(inputString)
	/**
	Jeder einzelne char wird mit diesen drei Op's dargestellt
	es ist lenght prefix based bedeutet das erste Word, welches gelesen wird ist die länge des Strings
	MOVI reg1 part
	ADDI reg2 0
	STOREB reg1 reg2
	*/
	formatted = [][]string{
		{"MOVI", rx, strconv.Itoa(length)},
		{"STOREW", rx, ry},
		{"ADDI", ry, "1"},
	}
	for _, part := range inputString {

		formatted = append(formatted,
			[]string{"MOVI", rx, strconv.Itoa(int(part))},
			[]string{"ADDI", ry, "1"},
			[]string{"STOREB", rx, ry},
		)
	}
	formatted = append(formatted,
		[]string{"SUBI", ry, strconv.Itoa(length + 1)})
	return formatted
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

func ParseLines(data string) [][]string {
	//turns into array removes comments
	stringLines := strings.Split(data, "\n")
	stringParts := make([][]string, len(stringLines))
	stringPartIndex := 0
	for index, line := range stringLines {
		line = strings.TrimSpace(line)
		commentIndex := strings.Index(line, "#")
		if commentIndex != -1 {
			line = line[:commentIndex]
			stringLines[index] = line
		}
		if line != "" {
			lineParts := strings.Fields(line)
			stringParts[stringPartIndex] = lineParts
			stringPartIndex++
		}
	}
	return stringParts[:stringPartIndex]
}
