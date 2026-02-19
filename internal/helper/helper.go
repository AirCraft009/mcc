package helper

import (
	"os"
	"path/filepath"
	"strings"
)

func EncodeRegs(rx byte, ry byte, addrnecs bool) (byte, byte) {
	var addrFlag byte
	if addrnecs {
		addrFlag = 0x1
	} else {
		addrFlag = 0x0
	}
	return rx, ry<<1 | addrFlag
}

func EncodeAddr(addr uint16) (byte, byte) {
	if addr <= 255 {
		return byte(0), byte(addr)
	}
	return byte(addr>>8) & 0xff, byte(addr & 0xff)
}

func InsertMatrixAtIndex(dest, insert [][]string, index int) [][]string {
	if index < 0 {
		index = 0
	} else if index > len(dest) {
		index = len(dest)
	}

	result := make([][]string, 0, len(dest)+len(insert))
	result = append(result, dest[:index]...)
	result = append(result, insert...)
	result = append(result, dest[index:]...)

	return result
}

func DeleteMatrixRow(matrix [][]string, index int) [][]string {
	if index < 0 || index >= len(matrix) {
		return matrix
	}

	return append(matrix[:index], matrix[index+1:]...)
}

func ConcactSliceAtIndex(dest, input []byte, index int) []byte {
	if len(dest)-int(index) < len(input) {
		return dest
	}
	for i := 0; i < len(input); i++ {
		dest[index+i] = input[i]
	}
	return dest
}

func SatSubU32(a, b uint32) uint32 {
	if a < b {
		return 0
	}
	return a - b
}

func GetRootPath() string {
	// get the exe path mcc.exe/bin/mcc.exe.exe
	exe, err := os.Executable()
	if err != nil {
		panic(err.Error())
	}

	// hopefully clean up any symlinks
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		panic(err.Error())
	}

	// go to mcc.exe/
	return filepath.Clean(filepath.Join(filepath.Dir(exe), "../"))
}

func IsLetter(r byte) bool {
	return r >= 'a' && r <= 'z'
}

func IsDigit(r byte) bool {
	return r >= '0' && r <= '9'
}

func Clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func StandardizeSpaceAmmount(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
