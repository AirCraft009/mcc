# MxsxllBinary Fileformats

## MXBO - Object files

### Usage

- It is a Object data
- used to contain 
  - bytecode
  - relocation entries
  - symbols


### Format

| byte                 | size          | name          | description                                                         |
|----------------------|---------------|---------------|---------------------------------------------------------------------|
| 0                    | 4B            | Header        | MXBO(MxsxllBinary Header - object)                                  |
| 4                    | 2B            | code-len      | lenght of the code                                                  |
| 6                    | 2B            | symbols-len   | lenght of the symbol-section                                        |
| 8                    | 2B            | reloc-len     | lenght of the relocation table                                      |
| 10                   | 1B            | entry         | old padding byte(used to signify this data as the entrypoint)       |
| 11                   | code-len      | code          | contains the bytecode without resolved labels                       |
| 11+code-len          | symbols-len   | symbols       | contains the symbols for addr resolution                            |
| ├─ len               | 1B            | symbol-len    | lenght of the symbol to come                                        |
| ├─ name              | symbol-len    | symbol-name   | the name of the symbol                                              |
| └─ gobal             | 1B            | symbol-global | 1 for a global and 0 for a local symbol                             |
| -                    | -             | -             | the previous three fields are repeated for every symbol in the code |
| symbols + symbol-len | reloc-len     | reloc-table   | contains information necesarry to resolving addresses               |
| ├─ offset            | 2B            | rel-offset    | the offset in the local data                                        |
| ├─ len               | 1B            | rel-label-len | the lenght of the label in str form                                 |
| └─ label             | rel-label-len | rel-label     | the name of the relocation entry                                    |
| -                    | -             | -             | the previous three fields are repeated for every relocation entry   |

- All numbers are encoded in little endian

## MXBI - Binary files

### Usage

- used for storing the files that get executed by MxsxllBox
- Can contain the debug labels necesarry for reverse Address resolution


### Format

| Byte       | Size            | name            | description                                                                              |
|------------|-----------------|-----------------|------------------------------------------------------------------------------------------|
| 0          | 4B              | Header          | MXBI(MxsxllBox Binary Header - binary)                                                   |
| 4          | 2B              | Code-len        | The code lenght should be 2^16(date 31.01.2026) cause the code fills the whole mem space |
| 6          | 2B              | debug-label-len | The lenght of the debug label entries                                                    |
| 8          | 1B              | debug           | Does this file have debug symbols                                                        |
| 9          | code-len        | code            | The code for MxsxllBox                                                                   |
| 9+code-len | debug-label-len | debug-labels    | consists of debug-label-len debug entries                                                |
| ├─len      | 2B              | entry-len       | The lenght of the label name str                                                         |
| ├─name     | entry-len       | entry-name      | The name of an label                                                                     |
| └─ address | 2B              | entry-addr      | The address this symbol belongs to                                                       |
| -          | -               | -               | the previous three fields are repeated for every debug-label entry                       |

- All numbers are encoded in little endian



