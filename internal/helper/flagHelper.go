package helper

import (
	"flag"
)

type FlagKeeper struct {
	Flags      *flag.FlagSet
	InputPath  string
	OutPath    string
	NoLink     bool
	Debug      bool
	Resolution bool
	Verbose    bool
	Supress    bool
}
