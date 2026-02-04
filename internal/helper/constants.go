package helper

const (
	MemorySize              = 1024 * 16
	ProgramStdLibStart      = 0x1800
	InterrupttableLoc       = 23965
	IncludeSignifier        = "#include"
	StdLibLocationSignifier = "stdlib"
	StdlibLocation          = "lib/stdlib"
	includeLocation         = "lib/include"
	GlobalHeaderLocation    = "/lib/globalHeaders"
	StdLibLocationUse       = StdlibLocation + "/obj"
	IncludeLocationUse      = includeLocation + "/obj"
)
