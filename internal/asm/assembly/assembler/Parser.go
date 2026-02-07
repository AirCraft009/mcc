// Package assembly
//
// Handles Assembling the .asm files into object files
package assembler

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"
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

// checkOffsetInstruction
//
// returns the label and offset from the
// assembler directive [label + off]
//
// example:
// STOREB RO [a+16]
//
// returns "a", 16
func checkOffsetInstruction(rawInstruction string) (string, int32) {
	if !strings.HasPrefix(rawInstruction, "[") || !strings.HasSuffix(rawInstruction, "]") {
		return rawInstruction, 0
	}

	instr := rawInstruction[1 : len(rawInstruction)-1]
	parts := regexp.MustCompile("[+-]").Split(instr, 2)
	if len(parts) != 2 {
		log.Printf("Instruction with square brackets but no offset instruction: %s\n", rawInstruction)
		return rawInstruction, 0
	}

	offset, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		log.Fatalf("Illegal instruction: %s\nNot a valid number after the offset signifier '+'\nerror: %s", rawInstruction, err.Error())
	}
	if strings.Contains(rawInstruction, "-") {
		offset = offset * -1
	}

	return strings.TrimSpace(parts[0]), int32(offset)
}

// firstPass
//
// seperates out labels and their Addresses
func (parser *Parser) firstPass(data [][]string, logger *log.Logger) [][]string {
	var PC uint16
	var activeLabel string
	var labelType LableT

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
			labelType = undefined
			activeLabel = rawLabel
			continue
			// check for formatters (smth that arranges code in other ways
		} else if actFormatter, ok := parser.Formatter[line[0]]; ok {
			if labelType == codeLabel {
				logger.Println(activeLabel)
				logger.Println(PC)
				helper.FatalWrapper(logger, "Tried to write Data to Code-label")
			}
			formatted, affectsPC := actFormatter(data[i], activeLabel, PC, parser)
			line = formatted
			data[i] = formatted
			if !affectsPC {
				// if it doesn't affect the PC it's an assembler directive
				// This also means that the label associated with this is now a data label
				labelType = dataLabel
				continue
			}
		}

		ad, ok := pkg.OffsetMap[line[0]]
		if !ok {
			logger.Println(line[0])
			logger.Println(PC)
			helper.FatalWrapper(logger, "Unknown offset")
		}
		if labelType == dataLabel {
			logger.Println("label: ", activeLabel)
			logger.Println("Instruction: ", line)
			logger.Println("PC: ", PC)
			helper.FatalWrapper(logger, "Tried to write Code to Data-label")
		}
		PC += uint16(ad)
	}
	parser.ObjFile.Symbols = parser.Labels
	return data
}

func (parser *Parser) secondPass(data [][]string, logger *log.Logger) (ObjFile *pkg.ObjectFile) {
	code := make([]byte, 0)

	PC := uint16(0)
	for _, line := range data {
		if parsfunc, ok := parser.Parsers[line[0]]; ok {
			codeSnippet := make([]byte, 2)
			var err error
			PC, codeSnippet, err = parsfunc(line, PC, parser)
			if err != nil {
				helper.FatalWrapper(logger, err.Error())
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
func Assemble(data string, logger *log.Logger) *pkg.ObjectFile {
	parsedData := parseLines(data)
	parser := newParser()
	var formattedData [][]string
	formattedData = parser.firstPass(parsedData, logger)

	return parser.secondPass(formattedData, logger)
}

// AssembleAndWrite assembles string asm files and writes them if wanted
func AssembleAndWrite(data, path string, write bool, logger *log.Logger) *pkg.ObjectFile {
	ObjFile := Assemble(data, logger)

	if write {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			helper.FatalWrapper(logger, err.Error())
		}
		defer f.Close()

		err = pkg.SaveObjectFile(ObjFile, f)
		if err != nil {
			helper.FatalWrapper(logger, err.Error())
		}
	}
	return ObjFile
}
