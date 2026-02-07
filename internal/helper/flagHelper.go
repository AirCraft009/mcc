package helper

import (
	"flag"
	"os"
)

type FlagKeeper struct {
	flags      *flag.FlagSet
	InputPath  string
	OutPath    string
	NoLink     bool
	Debug      bool
	Resolution bool
	Verbose    bool
	Supress    bool
}

func NewFlagKeeper() *FlagKeeper {
	return &FlagKeeper{
		flags: flag.NewFlagSet("initFlags", flag.ExitOnError),
	}
}

func (fk FlagKeeper) Parse() {
	fk.flags.StringVar(&fk.OutPath, "o", "a.bin", "output file")
	fk.flags.BoolVar(&fk.NoLink, "n", false, "do not use linker\n overrides debug and res because no full file is created")
	fk.flags.BoolVar(&fk.Debug, "debug", false, "creates debug symbols")
	fk.flags.BoolVar(&fk.Resolution, "res", false, "creates the object files at in the dir next to eachother")
	fk.flags.BoolVar(&fk.Verbose, "v", false, "verbose output")
	fk.flags.BoolVar(&fk.Supress, "s", false, "doesn't write to a log file at all")

	fk.flags.Parse(os.Args[1:])
	fk.InputPath = os.Args[0]
}
