#!/bin/sh
set -eu

BASEDIR="$(cd "$(dirname "$0")" && pwd)"
MCC="$BASEDIR/bin/mcc"

# ---- stdlib ----
for f in "$BASEDIR/stdlib/sources/"*.asm; do
    [ -e "$f" ] || continue
    echo "Compiling stdlib $(basename "$f")"
    "$MCC" "$f" \
        --o "$BASEDIR/stdlib/obj/$(basename "${f%.asm}.obj")" \
        --no_link=true
done

# ---- include ----
for f in "$BASEDIR/include/sources/"*.asm; do
    [ -e "$f" ] || continue
    echo "Compiling include $(basename "$f")"
    "$MCC" "$f" \
        --o "$BASEDIR/include/obj/$(basename "${f%.asm}.obj")" \
        --no_link=true
done
