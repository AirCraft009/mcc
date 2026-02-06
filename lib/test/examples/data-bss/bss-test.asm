a:
    .ZERO 2

b:
    .ZERO 1

start:
    LOADW R0 a
    PRINT R0
    # R0 == 0
    MOVI R1 129
    STOREW R1 a
    LOADW R0 a
    PRINT R0
    # R0 == 129

    STOREB