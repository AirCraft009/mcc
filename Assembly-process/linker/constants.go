package linker

const (
	MemorySize         = 1024 * 16
	ProgramStdLibStart = 0x1800
	InterrupttableLoc  = 23965
	includeSignifier   = "#include"
	stdLibLocation     = "stdlib"
	includeLocation    = "/include"
	stdLibLocationUse  = stdLibLocation + "/obj"
	includeLocationUse = includeLocation + "/obj"
)
