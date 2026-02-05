// Package assembler
//
// Handles Assembling the .asm files into object files
package assembler

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

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

// parseLines
//
// prepares the data into a format which is easily parseable
// it removes any whitespaces, comments and formats cleanly,
//
// it returns a string slice slice that has OPs and Params split
func parseLines(data string) [][]string {
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

func checkLabel(rawLabel string) (string, error) {
	_, err := strconv.Atoi(rawLabel)
	if err == nil {
		return "", errors.New("Label can't be a number: " + rawLabel)
	}

	if _, ok := pkg.RegMap[rawLabel]; ok {
		return "", errors.New("Label can't be a Register: " + rawLabel)
	}

	return rawLabel, nil
}

// firstPass
//
// seperates out labels and their Addresses
func (parser *Parser) firstPass(data [][]string) [][]string {
	var PC uint16
	var activeLabel string

	for i, line := range data {
		// a : signifies a label
		if len(line) == 1 && strings.Contains(line[0], ":") {
			rawLabel, err := checkLabel(line[0][:len(line[0])-1])
			if err != nil {
				panic(err.Error())
			}
			parser.Labels[rawLabel] = PC
			// an undersorce makes it global
			if strings.HasPrefix(line[0], "_") {
				parser.ObjFile.Globals[PC] = true
			}
			activeLabel = rawLabel
			continue
			// check for formatters (smth that arranges code in other ways
		} else if actFormatter, ok := parser.Formatter[line[0]]; ok {
			formatted, affectsPC := actFormatter(data[i], activeLabel, PC, parser)
			line = formatted
			if !affectsPC {
				continue
			}
		}

		ad, ok := pkg.OffsetMap[line[0]]
		if !ok {
			fmt.Println(line[0])
			fmt.Println(PC)
			panic("unknown Offset")
		}
		PC += uint16(ad)
	}
	parser.ObjFile.Symbols = parser.Labels
	return data
}

func (parser *Parser) secondPass(data [][]string) (ObjFile *pkg.ObjectFile) {
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

// Assemble
//
// assembles a .asm files (data) into an objectfile
// That conatins The bytecode with non relocated labels
// It first parses all lables to relocate,
// then finishes compiling into bytecode
func Assemble(data string) *pkg.ObjectFile {
	parsedData := parseLines(data)
	parser := newParser()
	var formattedData [][]string
	formattedData = parser.firstPass(parsedData)

	return parser.secondPass(formattedData)
}

// AssembleAndWrite assembles string asm files and writes them if wanted
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
