package pkg

const (
	Flash     = 16 * 1024 // the rom is 8 KB of storage I'm mapping to
	BitmapEnd = 256 * 8

	MemorySize = 1024 * 64 // 64 KB total memory

	// ProgramStart ───── Code Region (8 KB) ─────
	ProgramStart       = 0x0000
	ProgramUserEnd     = 0x17FF // 8 KB (User + StdLib)
	ProgramStdLibStart = 0x1800 // Last 2 KB for stdlib
	ProgramEnd         = 0x1FFF

	// HeapStart ───── Heap (16 KB) ─────
	HeapStart          = 0x2000
	HeapEnd            = 0x6000
	writeableHeapStart = 9628
	writeableHeapEnd   = 23964
	HeapSize           = writeableHeapEnd - writeableHeapStart
	Interrupttable     = 23965
	InterruptTableSIze = HeapEnd - Interrupttable
	BlockSize          = 0x10

	// StackStart ─────  (8 KB) ─────
	StackStart = 0x6000
	StackEnd   = 0x7FFF
	StackInit  = StackEnd

	// VideoStart ───── Video RAM / Framebuffer (16 KB) ─────
	VideoStart = 0x8000
	VideoEnd   = 0xBFFF

	// KeyboardStart ReservedStart ───── Reserved for IO / Buffers / MMIO (8 KB) ─────
	KeyboardStart   = 0xC000
	ReadPtr         = 0xC000
	WritePtr        = 0xC001
	RingBufferStart = 0xC002
	RingBufferEnd   = 0xC020 //N = 30
	RingBufferSize  = RingBufferEnd - RingBufferStart

	// ExtraStart ───── Unused / Future Expansion / Paging Tables / Filesystem etc (≈16KB KB) ─────
	ExtraStart = 0xC021

	DataStart = 0xE000
	DataEnd   = 0xE3FF
	DataSize  = DataEnd - DataStart

	BssSectionStart = 0xE400
	BssSectionEnd   = 0xE7FF
	BssSectionSize  = BssSectionEnd - BssSectionStart

	ExtraEnd  = 0xFFFF
	ExtraSize = ExtraEnd - ExtraStart
)
