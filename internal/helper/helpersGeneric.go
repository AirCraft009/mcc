package helper

import (
	"encoding/binary"
	"errors"
	"io"
	"sort"
)

// Helpers - For writing binary data

func GetSortedKeyVal[V any](m map[string]V) ([]string, []V) {
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

func WriteKV[V any](w io.Writer, key string, val V) error {
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

func ReadKV[V any](r io.Reader) (string, V, error) {
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
