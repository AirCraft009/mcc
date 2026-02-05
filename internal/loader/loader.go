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
	DataRelocs  [][]pkg.RelocationEntry
	BssPtr      uint16
	DataPtr     uint16
	DataValues  []byte
	BssSections map[string]uint16
	InitData    map[string][]byte
}

func LoadData(objs []*pkg.ObjectFile) {
	DatLoader := parseObjs(objs)
	DatLoader.setRelocs()
	for i, obj := range objs {
		obj.Relocs = DatLoader.DataRelocs[i]
	}
}

func (l *loader) setRelocs() {
	for _, relocs := range l.DataRelocs {
		for _, reloc := range relocs {
			if !reloc.Data {
				continue
			}

			if bssVal, ok := l.BssSections[reloc.Lbl]; ok {
				reloc.Offset = l.BssPtr
				l.BssPtr += bssVal
				if l.BssPtr > pkg.BssSectionEnd {
					panic("BSS section to large, failed at: " + reloc.Lbl)
				}
			} else if DataVal, ok := l.InitData[reloc.Lbl]; ok {
				reloc.Offset = l.DataPtr
				l.DataValues = append(l.DataValues, DataVal...)
				l.DataPtr += uint16(len(DataVal))
				if l.DataPtr > pkg.DataEnd {
					panic("Data section to large, failed at: " + reloc.Lbl)
				}
			}
		}
	}

}

func newLoader(size int) *loader {
	return &loader{
		make([][]pkg.RelocationEntry, size),
		pkg.BssSectionStart,
		pkg.DataStart,
		make([]byte, 0, pkg.DataSize),
		make(map[string]uint16),
		make(map[string][]byte, size),
	}
}

func parseObjs(objs []*pkg.ObjectFile) *loader {
	DatLoader := newLoader(len(objs))
	for i, obj := range objs {
		DatLoader.DataRelocs[i] = obj.Relocs
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
