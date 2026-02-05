package assembler

import (
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
