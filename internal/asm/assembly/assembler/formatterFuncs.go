package assembler

import (
	"errors"
	"log"
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
func StoreLoadFormatter(parameters []string, activeLabel string, currPC uint16, parser *Parser) ([]string, bool) {

	//STOREB R1 start
	// loc0 RegsLoc1 RegsLoc2

	// currPC has to point to exactly where the reloc is
	// currently at STOREB
	// + AddrOutLocHi it points to the label start
	if len(parameters) != 3 {
		log.Printf("MCC-WARN: %s seems to be malformed as it only has %d slots instead of 3", parameters, len(parameters))
		// will be caught later in the parser
		return parameters, true
	}

	currPC += AddrOutLocHi
	label, offset := checkOffsetInstruction(parameters[RegsLoc2])
	//fmt.Println("label: ", label, offset)
	rawLabel, err := checkLabel(label)
	//fmt.Println("rawLabel: ", rawLabel, err)
	if err != nil {
		parameters[RegsLoc2] = label
		parameters = append(parameters, strconv.Itoa(int(offset)))
		// Label is a number or Register
		return parameters, true
	}

	parser.ObjFile.Relocs = append(parser.ObjFile.Relocs, pkg.RelocationEntry{
		InFileOffset: currPC,
		Offset:       offset,
		Lbl:          rawLabel,
		Data:         true,
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
	parser.Labels[activeLabel] = 0
	parser.ObjFile.BssSections[activeLabel] = uint16(ammount)
	return parameters, false
}

func WordFormatter(parameters []string, activeLabel string, currPC uint16, parser *Parser) (newParams []string, affectsPC bool) {
	// how many zero-bytes
	val, err := strconv.Atoi(parameters[RegsLoc1])

	if err != nil {
		panic(errors.New(".ZERO takes an integer. Got " + parameters[RegsLoc1] + "\n err: " + err.Error()))
	}

	hi, lo := helper.EncodeAddr(uint16(val))

	parser.Labels[activeLabel] = 0
	data := parser.ObjFile.InitData[activeLabel]
	if data == nil {
		data = []byte{hi, lo}
	} else {
		data = append(data, hi, lo)
	}
	parser.ObjFile.InitData[activeLabel] = data
	return parameters, false
}

// ByteFormatter
//
// formats .BYTE instructions
// STOREB R0 LABEL
// STOREB R0 0
func ByteFormatter(parameters []string, activeLabel string, currPC uint16, parser *Parser) (newParams []string, affectsPC bool) {
	// how many zero-bytes
	val, err := strconv.Atoi(parameters[RegsLoc1])

	if err != nil {
		panic(errors.New(".ZERO takes an integer. Got " + parameters[RegsLoc1] + "\n err: " + err.Error()))
	}

	parser.Labels[activeLabel] = 0
	data := parser.ObjFile.InitData[activeLabel]
	if data == nil {
		data = []byte{byte(val)}
	} else {
		data = append(data, byte(val))
	}

	parser.ObjFile.InitData[activeLabel] = data
	//fmt.Println(data)
	return parameters, false
}
