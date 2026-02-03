package assembler

import (
	"fmt"
	"os"
	"strings"

	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"
)

const (
	OpCLoc         = 0
	RegsLoc1       = 1
	RegsLoc2       = 2
	RegsLocOut     = 1
	AddrLoc1       = 1
	AddrLoc2       = 2
	AddrOutLocHi   = 3
	AddrOutLocLo   = 4
	StrLoc         = 3
	RegWidthOffset = 1
)

// FirstPass
//
// seperates out labels and their Addresses
func (parser *Parser) FirstPass(data [][]string) [][]string {
	var PC uint16
	formattedExtraCode := make(map[int][][]string)

	for i, line := range data {
		// a : signifies a label
		if len(line) == 1 && strings.Contains(line[0], ":") {
			parser.Labels[line[0][:len(line[0])-1]] = PC
			// an undersorce makes it global
			if strings.HasPrefix(line[0], "_") {
				parser.ObjFile.Globals[PC] = true
			}

			continue
			// check for formatters (smth that arranges code in other ways
		} else if formatter, ok := parser.Formatter[strings.ToUpper(line[0])]; ok {
			formatted := formatter(data[i])
			formattedExtraCode[i] = formatted
			for _, formatLine := range formatted {
				ad, _ := pkg.OffsetMap[strings.ToUpper(formatLine[0])]
				PC += uint16(ad)
			}
			continue
		}

		ad, ok := pkg.OffsetMap[strings.ToUpper(line[0])]
		if !ok {
			fmt.Println(line[0])
			fmt.Println(PC)
			panic("unknown Offset")
		}
		PC += uint16(ad)
	}
	for i, extraData := range formattedExtraCode {
		data = helper.DeleteMatrixRow(data, i)
		data = helper.InsertMatrixAtIndex(data, extraData, i)
	}
	parser.ObjFile.Symbols = parser.Labels
	return data
}

func (parser *Parser) SecondPass(data [][]string) (ObjFile *pkg.ObjectFile) {
	code := make([]byte, 0)

	PC := uint16(0)
	for _, line := range data {
		if parsfunc, ok := parser.Parsers[line[0]]; ok {
			codeSnippet := make([]byte, 2)
			var err error
			PC, codeSnippet, err = parsfunc(line, PC, parser)
			if err != nil {
				panic(err)
			}
			code = append(code, codeSnippet...)
		}
	}
	parser.ObjFile.Code = code
	return parser.ObjFile
}

func Assemble(data string) *pkg.ObjectFile {
	parsedData := ParseLines(data)
	parser := newParser()
	var formattedData [][]string
	formattedData = parser.FirstPass(parsedData)

	return parser.SecondPass(formattedData)
}

// AssembleAndWrite
// assembles string asm files\
//
// returns an Objectfile containing relocation information
// & Code without resolved labels(0x0, 0x0)
func AssembleAndWrite(data, path string, write bool) *pkg.ObjectFile {
	ObjFile := Assemble(data)

	if write {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = pkg.SaveObjectFile(ObjFile, f)
		if err != nil {
			panic(err)
		}
	}
	return ObjFile
}
