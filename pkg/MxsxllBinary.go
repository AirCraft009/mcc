package pkg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

const MagicObject = "MXBO"
const MagicBinary = "MXBI"
const BinHeaderLen = len(MagicBinary) + 3

func SaveObjectFile(obj *ObjectFile, w io.Writer) error {
	w.Write([]byte(MagicObject)) //MxsxllBox-Object header
	binary.Write(w, binary.LittleEndian, uint16(len(obj.Code)))
	binary.Write(w, binary.LittleEndian, uint16(len(obj.Symbols)))
	binary.Write(w, binary.LittleEndian, uint16(len(obj.Relocs)))

	w.Write(obj.Code)

	for name, addr := range obj.Symbols {
		w.Write([]byte{byte(len(name))})
		w.Write([]byte(name))
		binary.Write(w, binary.LittleEndian, addr)
		if obj.Globals[addr] {
			w.Write([]byte{1})
		} else {
			w.Write([]byte{0})
		}
	}

	for _, rel := range obj.Relocs {
		binary.Write(w, binary.LittleEndian, rel.Offset)
		w.Write([]byte{byte(len(rel.Lbl))})
		w.Write([]byte(rel.Lbl))
	}

	return nil
}

func ReadObjectFile(path string) (*ObjectFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FormatObjectFile(data)
}

func FormatObjectFile(data []byte) (*ObjectFile, error) {
	buf := bytes.NewReader(data)

	header := make([]byte, 4)
	if _, err := buf.Read(header); err != nil {
		return nil, err
	}
	if string(header) != MagicObject {
		return nil, fmt.Errorf("invalid object file format: missing MXBI header")
	}

	var symbolCount uint16
	var codeLen uint16
	var relocCount uint16

	//.Read reads from r in this case buf into arg2 so &codelen
	binary.Read(buf, binary.LittleEndian, &codeLen)
	binary.Read(buf, binary.LittleEndian, &symbolCount)
	binary.Read(buf, binary.LittleEndian, &relocCount)

	code := make([]byte, codeLen)
	buf.Read(code)

	symbols := make(map[string]uint16)
	globals := make(map[uint16]bool)
	for i := 0; i < int(symbolCount); i++ {
		var nameLen uint8
		binary.Read(buf, binary.LittleEndian, &nameLen)
		nameBytes := make([]byte, nameLen)
		buf.Read(nameBytes)

		var addr uint16
		binary.Read(buf, binary.LittleEndian, &addr)

		var global byte
		binary.Read(buf, binary.LittleEndian, &global)
		if global == 1 {
			globals[addr] = true
		} else if global == 0 {
			globals[addr] = false
		} else {
			panic(fmt.Errorf("invalid global flag: %d", global))
		}

		symbols[string(nameBytes)] = addr
	}

	relocs := make([]RelocationEntry, relocCount)
	for i := 0; i < int(relocCount); i++ {
		binary.Read(buf, binary.LittleEndian, &relocs[i].Offset)
		var labelLen byte
		binary.Read(buf, binary.LittleEndian, &labelLen)
		name := make([]byte, labelLen)
		buf.Read(name)
		relocs[i].Lbl = string(name)
	}

	return &ObjectFile{Code: code, Symbols: symbols, Relocs: relocs, Globals: globals}, nil
}

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
