package helper

const (
	MemorySize           = 1024 * 16
	ProgramStdLibStart   = 0x1800
	InterrupttableLoc    = 23965
	IncludeSignifier     = "#include"
	StdLibLocation       = "stdlib"
	includeLocation      = "/include"
	GlobalHeaderLocation = "/globalHeaders"
	StdLibLocationUse    = StdLibLocation + "/obj"
	IncludeLocationUse   = includeLocation + "/obj"
)
