# Linking

This linker operates on already-loaded `ObjectFile`s with fixed base locations.
All file handling, parsing, and assembling happens before this step.

Linking consists of **two passes**:
1. Build a global label table
2. Resolve relocations and emit final code

---

## Global Lookup Table

Global labels are identified by the object file itself.

For every object file at base address `location`:

- Iterate over `objFile.Symbols`
- If the symbol’s relative address is marked global in `objFile.Globals`
- Insert into the Global Lookup Table as:


Duplicate global labels cause a hard error.

---

## Relocation Resolution

Relocations are processed per object file.

For each relocation:

1. Try resolving the label locally:
    - If found:
        - If the relocation is **not data**, add the object’s base location
        - **not data** means that the label is in either the bss or data section
2. Otherwise resolve via the Global Lookup Table
3. If the label does not exist, linking fails

The resolved address is always absolute.

The address is encoded into two bytes and written directly into the object’s code
at the relocation offset.

---

## Code Emission

After all relocations are applied:

- The object’s code is copied into the final memory buffer at its base location

Objects may share the same base address; later objects are placed directly after
earlier ones.

---

## Debug Symbols (Optional)

If debug mode is enabled:

- Every resolved relocation records:


- Objects without relocations still export all their symbols for debugging.

---

## Determinism

The linker is fully deterministic:

- Objects are linked in a fixed order
- No concurrency is used
- Relocations are applied sequentially

Same inputs always produce identical output.
