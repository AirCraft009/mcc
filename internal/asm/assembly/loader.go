package assembly

import (
	"github.com/AirCraft009/mcc/internal/helper"
	"github.com/AirCraft009/mcc/pkg"
)

/*
Plan

1. data von nicht data trennen
2. data relocs zuordnen (bss/data , für beide einen ptr führen
3. code bauen bss ist 0 aber data muss gefüllt werden
4. Adressen klar stellen
*/

type loader struct {
	BssPtr          uint16
	DataPtr         uint16
	DataValues      []byte
	BssSections     map[string]uint16
	InitData        map[string][]byte
	BSSRelocations  map[string]uint16
	DataRelocations map[string]uint16
}

func LoadData(objs []*pkg.ObjectFile) (data []byte) {
	DatLoader := parseObjs(objs)
	DatLoader.setRelocs()
	for _, obj := range objs {
		for lbl, addr := range DatLoader.BSSRelocations {
			obj.Symbols[lbl] = addr
			//fmt.Println("inside loader: ", lbl, addr)
		}
		for lbl, addr := range DatLoader.DataRelocations {
			obj.Symbols[lbl] = addr
			//fmt.Println("inside loader: ", lbl, addr)
		}
		//fmt.Println(obj.Symbols["a"])
	}
	return DatLoader.DataValues
}

func (l *loader) setRelocs() {

	bssLabels, ammounts := helper.GetSortedKeyVal(l.BssSections)
	dataLabels, dataPoints := helper.GetSortedKeyVal(l.InitData)

	for i := range len(ammounts) {
		bsslbl, ammount := bssLabels[i], ammounts[i]
		//fmt.Println("bsslbl", bsslbl)
		//fmt.Println("ammount", ammount)
		//fmt.Println("ptr", l.BssPtr)

		l.BSSRelocations[bsslbl] = l.BssPtr
		l.BssPtr += ammount
		if l.BssPtr > pkg.BssSectionEnd {
			panic("BSS section to large, failed at: " + bsslbl)
		}

	}

	for i := range len(dataLabels) {
		datalbl, data := dataLabels[i], dataPoints[i]
		//fmt.Println(datalbl, data, l.DataPtr)

		l.DataRelocations[datalbl] = l.DataPtr
		l.DataPtr += uint16(len(data))
		l.DataValues = append(l.DataValues, data...)
		if l.DataPtr > pkg.DataEnd {
			panic("Data section to large, failed at: " + datalbl)
		}
		l.DataValues = append(l.DataValues, data...)
	}

}

func newLoader(size int) *loader {
	return &loader{
		pkg.BssSectionStart,
		pkg.DataStart,
		make([]byte, 0, pkg.DataSize),
		make(map[string]uint16),
		make(map[string][]byte, size),
		make(map[string]uint16, size),
		make(map[string]uint16, size),
	}
}

func parseObjs(objs []*pkg.ObjectFile) *loader {
	DatLoader := newLoader(len(objs))

	for _, obj := range objs {
		DatLoader.BssSections = combineMaps(DatLoader.BssSections, obj.BssSections)
		DatLoader.InitData = combineMaps(DatLoader.InitData, obj.InitData)
		//fmt.Println("InitData", obj.InitData)
	}

	return DatLoader
}

func combineMaps[K comparable, V any](a, b map[K]V) map[K]V {
	for k, v := range b {
		a[k] = v
	}
	return a
}
