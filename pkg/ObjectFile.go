package pkg

type ObjectFile struct {
	Code     []byte
	Symbols  map[string]uint16
	Relocs   []RelocationEntry
	Globals  map[uint16]bool
	BssPtr   uint16
	DataPtr  uint16
	InitData map[string][]byte
}

type RelocationEntry struct {
	Offset uint16 // Where in Code the label is called/JMP'd to
	Lbl    string
}

func NewObjectFile() *ObjectFile {
	return &ObjectFile{
		Code:     nil,
		Symbols:  make(map[string]uint16),
		Relocs:   make([]RelocationEntry, 0),
		Globals:  make(map[uint16]bool),
		BssPtr:   BssSectionStart,
		DataPtr:  DataStart,
		InitData: make(map[string][]byte),
	}
}
