package loader

import (
	"github.com/AirCraft009/mcc/pkg"
)

/*
Plan

TODO:
1. data von nicht data trennen
2. data relocs zuordnen (bss/data , für beide einen ptr führen
3. code bauen bss ist 0 aber data muss gefüllt werden
4. Adressen klar stellen
*/

type loader struct {
	DataRelocs  []pkg.RelocationEntry
	BssPtr      uint16
	DataPtr     uint16
	DataValues  []byte
	BssSections map[string]uint16
	InitData    map[string][]byte
}

func (loader) LoadData() {

}

func newLoader(size int) *loader {
	return &loader{
		make([]pkg.RelocationEntry, 0, size),
		pkg.BssSectionStart,
		pkg.DataStart,
		make([]byte, 0, size),
		make(map[string]uint16),
		make(map[string][]byte, size),
	}
}

func ParseObjs(objs []*pkg.ObjectFile) *loader {
	DatLoader := newLoader(len(objs))
	for _, obj := range objs {
		DatLoader.DataRelocs = append(DatLoader.DataRelocs, obj.Relocs...)
		DatLoader.BssSections = combineMaps(DatLoader.BssSections, obj.BssSections)
		DatLoader.InitData = combineMaps(DatLoader.InitData, obj.InitData)
	}

	return DatLoader
}

func combineMaps[K comparable, V any](a, b map[K]V) map[K]V {
	for k, v := range b {
		a[k] = v
	}
	return a
}
