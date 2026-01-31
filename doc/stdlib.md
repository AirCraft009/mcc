## Std. Lib.

- Labels that are saved into the library Region of the Program Space
- ProgramStdLibStart = 0x0C01
- ProgramEnd         = 0x0FFF
- In Syscalls O1 is always the return addr if there's only 1 return value

### String functions

- `_strcpy`: copies a String to O1 from O2
- `_strlen`: Loads len(O2) into O1
- `_strcmp`: compares two strings sets 0 flag if they are equal carry-flag if a byte is higher
- `_strcat`: concacts two strings O1 = O1+O2 `"a", "b" = "ab"`

### Utility funtctions

- `_memset`: sets a region of memory(lenght O3) starting at addr(O1) to val(O2)
- `_memcpy`: copies from addr(O1) to addr(O2) for ammount(O3) bytes

### sys functions

- `_alloc`: allocates ammount(O2) blocks(16 B) the start is returned in O1 (0 if OOM)
- `_free`: frees a block of memory O1 is the start of that Memory


### math functions

- basic arithmetic:
    - `_add`/`_sub`/`_mul`/`_div`/`_mod`/`_inc`/`_dec`
- `_max`: returns the larger val of O1/O2 in O1
- `_min`: return the smaller val of O1/O2 in O1
- `_pow`: returns O1**O2 in O1
- `_clamp`: if Val(O2) is between two values(low O1, high O3) the 0 flag is set

### io functions

- `_printstr`: prints a string from O2 is also in instruction set
- `_printchar`: prints a char from O2 O1 contains char
- `_readchar`: reads a char from the Keyboard Buffer into O1
