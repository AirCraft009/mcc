package helper

const (
	MemorySize           = 1024 * 16
	ProgramStdLibStart   = 0x1800
	InterrupttableLoc    = 23965
	IncludeSignifier     = "#include"
	StdLibLocation       = "stdlib"
	StdlibLocationUse    = "lib/stdlib"
	includeLocationUse   = "lib/include"
	GlobalHeaderLocation = "lib/globalHeaders"
	StdLibLocationUse    = StdlibLocationUse + "/obj"
	IncludeLocationUse   = includeLocationUse + "/obj"
)
