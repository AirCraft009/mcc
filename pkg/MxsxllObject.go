package pkg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sort"
)

func SaveObjectFile(obj *ObjectFile, w io.Writer) error {
	// header
	if _, err := w.Write([]byte(MagicObject)); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, uint16(len(obj.Code))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint16(len(obj.Symbols))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint16(len(obj.Relocs))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint16(len(obj.BssSections))); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint16(len(obj.InitData))); err != nil {
		return err
	}

	// code
	if _, err := w.Write(obj.Code); err != nil {
		return err
	}

	// symbols (sorted)
	names, addrs := getSortedKeyVal(obj.Symbols)
	for i := range names {
		if err := writeKV(w, names[i], addrs[i]); err != nil {
			return err
		}
		var g byte
		if obj.Globals[addrs[i]] {
			g = 1
		}
		if err := binary.Write(w, binary.LittleEndian, g); err != nil {
			return err
		}
	}

	// relocs
	for _, r := range obj.Relocs {
		if err := writeKV(w, r.Lbl, r.Offset); err != nil {
			return err
		}
		var d byte
		if r.Data {
			d = 1
		}
		if err := binary.Write(w, binary.LittleEndian, d); err != nil {
			return err
		}
	}

	// BSS
	bssNames, bssSizes := getSortedKeyVal(obj.BssSections)
	for i := range bssNames {
		if err := writeKV(w, bssNames[i], bssSizes[i]); err != nil {
			return err
		}
	}

	// init data
	initNames, initVals := getSortedKeyVal(obj.InitData)
	for i := range initNames {
		name := initNames[i]
		val := initVals[i]

		if len(name) > 255 {
			return errors.New("init label too long")
		}
		if err := binary.Write(w, binary.LittleEndian, uint8(len(name))); err != nil {
			return err
		}
		if _, err := w.Write([]byte(name)); err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, uint16(len(val))); err != nil {
			return err
		}
		if _, err := w.Write(val); err != nil {
			return err
		}
	}

	return nil
}

func FormatObjectFile(data []byte) (*ObjectFile, error) {
	r := bytes.NewReader(data)

	header := make([]byte, 4)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}
	if string(header) != MagicObject {
		return nil, errors.New("invalid object file header")
	}

	var (
		codeLen       uint16
		symbolCount   uint16
		relocCount    uint16
		bssCount      uint16
		initDataCount uint16
	)

	binary.Read(r, binary.LittleEndian, &codeLen)
	binary.Read(r, binary.LittleEndian, &symbolCount)
	binary.Read(r, binary.LittleEndian, &relocCount)
	binary.Read(r, binary.LittleEndian, &bssCount)
	binary.Read(r, binary.LittleEndian, &initDataCount)

	code := make([]byte, codeLen)
	io.ReadFull(r, code)

	symbols := make(map[string]uint16, symbolCount)
	globals := make(map[uint16]bool)

	for i := 0; i < int(symbolCount); i++ {
		name, addr, err := readKV[uint16](r)
		if err != nil {
			return nil, err
		}
		var g byte
		binary.Read(r, binary.LittleEndian, &g)
		globals[addr] = g == 1
		symbols[name] = addr
	}

	relocs := make([]RelocationEntry, relocCount)
	for i := 0; i < int(relocCount); i++ {
		name, off, err := readKV[uint16](r)
		if err != nil {
			return nil, err
		}
		var d byte
		binary.Read(r, binary.LittleEndian, &d)
		relocs[i] = RelocationEntry{
			Lbl:    name,
			Offset: off,
			Data:   d == 1,
		}
	}

	bss := make(map[string]uint16, bssCount)
	for i := 0; i < int(bssCount); i++ {
		name, size, err := readKV[uint16](r)
		if err != nil {
			return nil, err
		}
		bss[name] = size
	}

	initData := make(map[string][]byte, initDataCount)
	for i := 0; i < int(initDataCount); i++ {
		var n uint8
		binary.Read(r, binary.LittleEndian, &n)

		name := make([]byte, n)
		io.ReadFull(r, name)

		var sz uint16
		binary.Read(r, binary.LittleEndian, &sz)

		val := make([]byte, sz)
		io.ReadFull(r, val)

		initData[string(name)] = val
	}

	return &ObjectFile{
		Code:        code,
		Symbols:     symbols,
		Relocs:      relocs,
		Globals:     globals,
		BssSections: bss,
		InitData:    initData,
	}, nil
}

func ReadObjectFile(path string) (*ObjectFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FormatObjectFile(data)
}

// Helpers - For writing binary data

func getSortedKeyVal[V any](m map[string]V) ([]string, []V) {
	sortedKeys := make([]string, 0, len(m))
	sortedValues := make([]V, 0, len(m))

	for k, _ := range m {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		sortedValues = append(sortedValues, m[k])
	}
	return sortedKeys, sortedValues
}

func writeKV[V any](w io.Writer, key string, val V) error {
	if len(key) > 255 {
		return errors.New("key too long")
	}
	if err := binary.Write(w, binary.LittleEndian, uint8(len(key))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(key)); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, val)
}

func readKV[V any](r io.Reader) (string, V, error) {
	var n uint8
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		var zero V
		return "", zero, err
	}

	name := make([]byte, n)
	if _, err := io.ReadFull(r, name); err != nil {
		var zero V
		return "", zero, err
	}

	var val V
	if err := binary.Read(r, binary.LittleEndian, &val); err != nil {
		return "", val, err
	}

	return string(name), val, nil
}
