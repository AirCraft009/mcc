package helper

const (
	MemorySize              = 1024 * 64
	ProgramStdLibStart      = 0x1800
	InterrupttableLoc       = 23965
	IncludeSignifier        = "#include"
	StdLibLocationSignifier = "stdlib"
	StdlibLocation          = "lib/stdlib"
	includeLocation         = "lib/include"
	GlobalHeaderLocation    = "/lib/globalHeaders"
	StdLibLocationUse       = StdlibLocation + "/obj"
	IncludeLocationUse      = includeLocation + "/obj"

	DataStart = 0xE400
	DataEnd   = 0xE3FF

	BssSectionStart = 0xE400
	BssSectionEnd   = 0xE7FF
)
