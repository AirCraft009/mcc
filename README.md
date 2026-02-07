# MCC - Mxsxll Compiler Collection
Neither a compiler nor a collection 

## Setup

- download a Release
- place it in a dir like C:/
- add mcc/bin to PATH

## Usage

Usage of mcc:
    -debug creates debug symbols
    -n do not use linker
        overrides debug and res because no full file is created
    -o string
        output file (default "a.bin")
    -res
        creates the object files at in the dir next to eachother
    -s doesn't write to a log file at all
    -v verbose output - log output to stderr

## Build

- run ```go build -o bin/mcc.exe -v ./Assembler-main``` to build the project
- run ```./buildHelper/startup.(bat/sh)``` to build the project and the stdlib

## Writing Code

- The Abi for MCC/MxsxllBox can be found [here](https://github.com/AirCraft009/MxsxllBox/blob/master/docs/abi.md)
- All instructions (OP codes) can be found [here](https://github.com/AirCraft009/MxsxllBox/blob/master/docs/instruction-set.md)
- Stdlib documentation can be found [here](https://github.com/AirCraft009/mcc/blob/master/doc/stdlib.md)
