package helper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func FatalWrapper(logger *log.Logger, msg string) {
	fmt.Println(msg)
	fmt.Println("-----")
	logger.Fatal(msg)
}

func FatalFWrapper(logger *log.Logger, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	FatalWrapper(logger, msg)
}
