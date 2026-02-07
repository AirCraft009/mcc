package assembler

import (
	"github.com/AirCraft009/mcc/pkg"
)

type formatter func(parameters []string, activeLabel string, currPC uint16, parser *Parser) (newParams []string, affectsPC bool)
type parser func(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error)

type LableT byte

const (
	undefined LableT = iota
	codeLabel
	dataLabel
)

type Parser struct {
	Parsers   map[string]parser
	Formatter map[string]formatter
	Labels    map[string]uint16
	ObjFile   *pkg.ObjectFile
}

func newParser() *Parser {
	parser := &Parser{
		Parsers:   make(map[string]parser),
		Formatter: make(map[string]formatter),
		Labels:    make(map[string]uint16),
		ObjFile:   pkg.NewObjectFile(),
	}

	parser.Parsers["NOP"] = parseFormatOP
	parser.Parsers["LOADB"] = parseFormatOPRegAddr
	parser.Parsers["LOADW"] = parseFormatOPRegAddr
	parser.Parsers["LOAD"] = parseFormatOPRegAddr
	parser.Parsers["STOREB"] = parseFormatOPRegAddr
	parser.Parsers["STOREW"] = parseFormatOPRegAddr
	parser.Parsers["STORE"] = parseFormatOPRegAddr
	parser.Parsers["MOVI"] = parseFormatOPRegAddr
	parser.Parsers["DIVI"] = parseFormatOPRegAddr
	parser.Parsers["MULI"] = parseFormatOPRegAddr
	parser.Parsers["SUBI"] = parseFormatOPRegAddr
	parser.Parsers["ADDI"] = parseFormatOPRegAddr
	parser.Parsers["ADD"] = parseFormatOPRegReg
	parser.Parsers["SUB"] = parseFormatOPRegReg
	parser.Parsers["DIV"] = parseFormatOPRegReg
	parser.Parsers["MUL"] = parseFormatOPRegReg
	parser.Parsers["JMP"] = parseFormatOPLbl
	parser.Parsers["JZ"] = parseFormatOPLbl
	parser.Parsers["JC"] = parseFormatOPLbl
	parser.Parsers["CALL"] = parseFormatOPLbl
	parser.Parsers["PUSH"] = parseFormatOPReg
	parser.Parsers["POP"] = parseFormatOPReg
	parser.Parsers["PRINT"] = parseFormatOPReg
	parser.Parsers["RET"] = parseFormatOP
	parser.Parsers["HALT"] = parseFormatOP
	parser.Parsers["PRINTSTR"] = parseFormatOPReg
	parser.Parsers["JNZ"] = parseFormatOPLbl
	parser.Parsers["JNC"] = parseFormatOPLbl
	parser.Parsers["CMP"] = parseFormatOPRegReg
	parser.Parsers["CMPI"] = parseFormatOPRegAddr
	parser.Parsers["TEST"] = parseFormatOPRegReg
	parser.Parsers["TSTI"] = parseFormatOPRegAddr
	parser.Parsers["JL"] = parseFormatOPLbl
	parser.Parsers["JLE"] = parseFormatOPLbl
	parser.Parsers["JG"] = parseFormatOPLbl
	parser.Parsers["JGE"] = parseFormatOPLbl
	parser.Parsers["STZ"] = parseFormatOP
	parser.Parsers["STC"] = parseFormatOP
	parser.Parsers["CLZ"] = parseFormatOP
	parser.Parsers["CLC"] = parseFormatOP
	parser.Parsers["MOD"] = parseFormatOPRegReg
	parser.Parsers["MODI"] = parseFormatOPRegAddr
	parser.Parsers["MOV"] = parseFormatOPRegReg
	parser.Parsers["RS"] = parseFormatOPRegReg
	parser.Parsers["LS"] = parseFormatOPRegReg
	parser.Parsers["OR"] = parseFormatOPRegReg
	parser.Parsers["AND"] = parseFormatOPRegReg
	parser.Parsers["MOVA"] = parseFormatOPRegLbl
	parser.Parsers["GPC"] = parseFormatOPReg
	parser.Parsers["SPC"] = parseFormatOPReg
	parser.Parsers["GSP"] = parseFormatOPReg
	parser.Parsers["SSP"] = parseFormatOPReg
	parser.Parsers["GRFN"] = parseFormatOPRegReg
	parser.Parsers["GF"] = parseFormatOPReg
	parser.Parsers["SF"] = parseFormatOPReg
	parser.Parsers["SRFN"] = parseFormatOPRegReg
	parser.Parsers["YIELD"] = parseFormatOP
	parser.Parsers["UNYIELD"] = parseFormatOP
	parser.Parsers["STINTI"] = parseFormatOPAddr
	parser.Parsers["STINT"] = parseFormatOPReg
	parser.Parsers["XOR"] = parseFormatOPRegReg
	parser.Parsers["DRAWPX"] = parseFormatOPRegReg
	parser.Parsers["STOREBLOCK"] = parseFormatOPRegReg

	parser.Formatter["STORE"] = StoreLoadFormatter
	parser.Formatter["STOREB"] = StoreLoadFormatter
	parser.Formatter["STOREW"] = StoreLoadFormatter
	parser.Formatter["LOAD"] = StoreLoadFormatter
	parser.Formatter["LOADB"] = StoreLoadFormatter
	parser.Formatter["LOADW"] = StoreLoadFormatter

	parser.Formatter[".ZERO"] = ZeroFormatter
	parser.Formatter[".WORD"] = WordFormatter
	parser.Formatter[".BYTE"] = ByteFormatter
	//	parser.Formatter[".STRING"]
	return parser
}
