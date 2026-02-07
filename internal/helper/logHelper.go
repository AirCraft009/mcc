package helper

import (
	"fmt"
	"io"
	"log"
	"os"
)

func InitLogger(supress, verbose bool) *log.Logger {
	var myLogger *log.Logger

	if !supress {
		if !verbose {
			create, err := os.Create("Mcc-Logger.log")
			if err != nil {
				return log.New(io.Discard, "", 0)
			}
			fmt.Printf("writing logs to: %s\n", create.Name())
			myLogger = log.New(create, "Mcc-assmbler:", log.LstdFlags|log.Lshortfile)

		} else {
			myLogger = log.New(os.Stderr, "Mcc-assembler:", log.LstdFlags)
		}

	} else {
		myLogger = log.New(io.Discard, "", 0)
	}
	return myLogger
}

func FatalWrapper(logger *log.Logger, msg string) {
	fmt.Println(msg)
	fmt.Println("-----")
	logger.Fatal(msg)
}

func FatalFWrapper(logger *log.Logger, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	FatalWrapper(logger, msg)
}

func ConfigureGlobalLogger() {
	// init the logger for warnings (always to stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("MCC-WARN: ")
}
