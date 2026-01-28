## MxsxllBinary Fileformat

### Usage

- It is a Object data
- used to contain 
  - bytecode
  - relocation entries
  - symbols


### Format

| byte                 | size          | name          | description                                                         |
|----------------------|---------------|---------------|---------------------------------------------------------------------|
| 0                    | 4B            | Header        | MXBI(MxsxllBinary Header)                                           |
| 4                    | 4B            | code-len      | lenght of the code                                                  |
| 8                    | 4B            | symbols-len   | lenght of the symbol-section                                        |
| 12                   | 4B            | reloc-len     | lenght of the relocation table                                      |
| 16                   | 1B            | entry         | old padding byte(used to signify this data as the entrypoint)       |
| 17                   | code-len      | code          | contains the bytecode without resolved labels                       |
| 17+code-len          | symbols-len   | symbols       | contains the symbols for addr resolution                            |
| ├─ len               | 1B            | symbol-len    | lenght of the symbol to come                                        |
| ├─ name              | symbol-len    | symbol-name   | the name of the symbol                                              |
| └─ gobal             | 1B            | symbol-global | 1 for a global and 0 for a local symbol                             |
| -                    | -             | -             | the previous three fields are repeated for every symbol in the code |
| symbols + symbol-len | reloc-len     | reloc-table   | contains information necesarry to resolving addresses               |
| ├─ offset            | 4B            | rel-offset    | the offset in the local data                                        |
| ├─ len               | 1B            | rel-label-len | the lenght of the label in str form                                 |
| └─ label             | rel-label-len | rel-label     | the name of the relocation entry                                    |
| -                    | -             | -             | the previous three fields are repeated for every relocation entry   |