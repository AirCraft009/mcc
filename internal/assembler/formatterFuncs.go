package assembler

import (
	"errors"
	"strconv"

	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"
)

// StoreLoadFormatter
//
// formats Instructions of the STORE Type
//
// - STOREB
// - STOREW
// - STORE
//
// - LOADB
// - LOADW
// - LOAD
func StoreLoadFormatter(parameters []string, activeLabel string, currPC uint16, parser *Parser) (newParams []string, affectsPC bool) {

	//STOREB R1 start
	// loc0 RegsLoc1 RegsLoc2

	// currPC has to point to exactly where the reloc is
	// currently at STOREB
	// + AddrOutLocHi it points to the label start

	currPC += AddrOutLocHi
	label, err := checkLabel(parameters[RegsLoc2])
	if err != nil {
		// Label is a number or Register
		return parameters, true
	}
	parser.ObjFile.Relocs = append(parser.ObjFile.Relocs, pkg.RelocationEntry{
		Offset: currPC,
		Lbl:    label,
	})
	// set 0 STORE instructions can handle real number they will later get replaced by the relocation

	parameters[RegsLoc2] = "0"
	return parameters, true
}

func ZeroFormatter(parameters []string, activeLabel string, currPC uint16, parser *Parser) (newParams []string, affectsPC bool) {
	// how many zero-bytes
	ammount, err := strconv.Atoi(parameters[RegsLoc1])

	if err != nil {
		panic(errors.New(".ZERO takes an integer. Got " + parameters[RegsLoc1] + "\n err: " + err.Error()))
	}

	parser.ObjFile.Relocs = append(parser.ObjFile.Relocs, pkg.RelocationEntry{
		Offset: parser.BssPtr,
		Lbl:    activeLabel,
	})

	parser.BssPtr += uint16(ammount)
	return parameters, false
}

func WordFormatter(parameters []string, activeLabel string, currPC uint16, parser *Parser) (newParams []string, affectsPC bool) {
	// how many zero-bytes
	val, err := strconv.Atoi(parameters[RegsLoc1])

	if err != nil {
		panic(errors.New(".ZERO takes an integer. Got " + parameters[RegsLoc1] + "\n err: " + err.Error()))
	}

	parser.ObjFile.Relocs = append(parser.ObjFile.Relocs, pkg.RelocationEntry{
		Offset: parser.DataPtr,
		Lbl:    activeLabel,
	})
	hi, lo := helper.EncodeAddr(uint16(val))
	parser.InitData[activeLabel] = []byte{hi, lo}

	// Word == 2Bytes
	parser.DataPtr += 2
	return []string{}, false
}

func ByteFormatter(parameters []string, activeLabel string, currPC uint16, parser *Parser) (newParams []string, affectsPC bool) {
	// how many zero-bytes
	val, err := strconv.Atoi(parameters[RegsLoc1])

	if err != nil {
		panic(errors.New(".ZERO takes an integer. Got " + parameters[RegsLoc1] + "\n err: " + err.Error()))
	}

	parser.ObjFile.Relocs = append(parser.ObjFile.Relocs, pkg.RelocationEntry{
		Offset: parser.DataPtr,
		Lbl:    activeLabel,
	})
	parser.InitData[activeLabel] = []byte{byte(val)}

	// 1Byte
	parser.DataPtr += 1
	return []string{}, false
}
