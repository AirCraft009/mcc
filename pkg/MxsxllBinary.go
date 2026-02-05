package pkg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"sort"
)

const MagicObject = "MXBO"
const MagicBinary = "MXBI"
const BinHeaderLen = len(MagicBinary) + 3

func FormatMxBinary(code []byte, debugSymbols map[uint16]string, debug bool) (data []byte) {
	if len(code) != MemorySize {
		panic(fmt.Errorf("invalid code length: %d", len(code)))
	}
	var buf bytes.Buffer = *bytes.NewBuffer(make([]byte, 0, len(code)+MemorySize+BinHeaderLen))

	// 4Bytes Header
	buf.Write([]byte(MagicBinary))
	// 2 Bytes code lenght
	// - 1 so it dosen't overflow 64 kb
	binary.Write(&buf, binary.LittleEndian, uint16(len(code)-1))
	binary.Write(&buf, binary.LittleEndian, uint16(len(debugSymbols)))

	if !debug {
		binary.Write(&buf, binary.LittleEndian, byte(0))
		buf.Write(code)
		return buf.Bytes()
	}

	binary.Write(&buf, binary.LittleEndian, byte(1))
	buf.Write(code)
	names := make([]string, 0, len(debugSymbols))
	addresses := make([]int, 0, len(debugSymbols))

	for addr, _ := range debugSymbols {
		addresses = append(addresses, int(addr))
	}

	sort.Ints(addresses)

	for _, addr := range addresses {
		names = append(names, debugSymbols[uint16(addr)])
	}

	for i := range len(names) {
		name := names[i]
		address := addresses[i]

		binary.Write(&buf, binary.LittleEndian, uint16(len(name)))
		buf.Write([]byte(name))
		binary.Write(&buf, binary.LittleEndian, uint16(address))
	}
	return buf.Bytes()
}

func WriteMxBinary(outputFile string, code []byte, debugSymbols map[uint16]string, debug bool) (err error) {
	data := FormatMxBinary(code, debugSymbols, debug)
	err = os.WriteFile(outputFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetDataFromMxBinary(data []byte) (code []byte, debugSymbols map[uint16]string, debug bool, err error) {
	reader := bytes.NewReader(data)
	header := make([]byte, 4)
	code = make([]byte, MemorySize)

	binary.Read(reader, binary.LittleEndian, &header)
	if string(header) != MagicBinary {
		return nil, nil, false, errors.New("invalid header: " + string(header) + "instead of" + MagicBinary)
	}

	var codeLen uint16
	var symbolCount uint16
	var Isdebug byte

	binary.Read(reader, binary.LittleEndian, &codeLen)

	if codeLen != MemorySize-1 {
		return nil, nil, false, fmt.Errorf("invalid code length: %d expected %d\nCode needs to fill the full memory", codeLen, MemorySize-1)
	}

	binary.Read(reader, binary.LittleEndian, &symbolCount)
	binary.Read(reader, binary.LittleEndian, &Isdebug)

	binary.Read(reader, binary.LittleEndian, &code)

	if Isdebug == 0 {
		return code, nil, false, nil
	}

	debugSymbols = make(map[uint16]string)
	var nameLen uint16
	var nameBytes []byte
	var addr uint16

	for _ = range symbolCount {
		binary.Read(reader, binary.LittleEndian, &nameLen)
		nameBytes = make([]byte, nameLen)
		binary.Read(reader, binary.LittleEndian, &nameBytes)
		binary.Read(reader, binary.LittleEndian, &addr)
		debugSymbols[addr] = string(nameBytes)
	}

	return code, debugSymbols, true, nil
}

func ReadMxBinary(inputFile string) (code []byte, debugSymbols map[uint16]string, debug bool, err error) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, nil, false, err
	}

	return GetDataFromMxBinary(data)
}
