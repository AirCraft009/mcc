package pkg

// ObjectFile
//
// Code - The code with no label-addresses
// Symbols - maps the label to the relativ addr
// Relocs - Saves unkown labels
// Globals - maps labels to if they're global
// BssPtr - How much of Bss has been used
// BssSections - maps bss labels to the size in bss
// DataPtr - How much of Data has been used
// InitData - Maps data labels to values
type ObjectFile struct {
	Code        []byte
	Symbols     map[string]uint16
	Relocs      []RelocationEntry
	Globals     map[uint16]bool
	BssSections map[string]uint16
	InitData    map[string][]byte
}

// RelocationEntry
//
// Used for relocating labels to addr
// Offset - the offset in the file
// Lbl - The label-name
// Data - The label is a data label and belongs in bss/data section
type RelocationEntry struct {
	Offset uint16 // Where in Code the label is called/JMP'd to
	Lbl    string
	Data   bool
}

// NewObjectFile
//
// returns a base *ObjectFile
// Code - is not initialized
func NewObjectFile() *ObjectFile {
	return &ObjectFile{
		Code:        nil,
		Symbols:     make(map[string]uint16),
		Relocs:      make([]RelocationEntry, 0),
		Globals:     make(map[uint16]bool),
		BssSections: make(map[string]uint16),
		InitData:    make(map[string][]byte),
	}
}
