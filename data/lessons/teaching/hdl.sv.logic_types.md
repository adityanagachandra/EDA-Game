# The logic type, packed arrays and structs

SystemVerilog's data types give a front end most of its bookkeeping: the four-state logic type, packed arrays that form one contiguous vector, packed structs with named fields and unpacked arrays of elements. Bit counts and select widths follow from a few rules, and getting them right is the first thing elaboration and lint must do before any logic is built.

## Packed arrays

Dimensions written before the name are packed: logic [3:0][7:0] bus; is one 32-bit vector, organized as four 8-bit bytes. Packed dimensions multiply, so logic [2:0][1:0][7:0] holds $3 \times 2 \times 8 = 48$ bits. A packed array can be used anywhere a vector can: in arithmetic, in comparisons and in assignments to other vectors of the same width.

Indexing selects along the leftmost dimension first. bus[2] is the third byte, 8 bits, and bus[2][3] is one bit. The logic type is four-state and can be written by procedural blocks or by one continuous assignment; the two-state bit type models 0 and 1 only and is used where x makes no sense.

## Packed structs

A typedef struct packed groups named fields into one vector whose width is the sum of the fields. With fields logic [3:0] tag, logic valid and logic parity, the struct is $4 + 1 + 1 = 6$ bits, and the first field occupies the most significant bits. A packed array of such structs, pkt_t [3:0] q, holds $4 \times 6 = 24$ bits.

Because a packed struct is just a vector with names, it can be assigned to or from a plain vector of the same width, compared and passed through ports. Packages typically hold these typedefs, so that every module that uses a bus agrees on the field layout and total width, which is exactly the consistency an interface needs.

## Unpacked arrays and selects

Dimensions written after the name are unpacked: logic [3:0][7:0] mem [128]; is 128 separate elements, each a 32-bit packed vector. Total storage is $128 \times 4 \times 8 = 4096$ bits, but the array is not one vector: it cannot be used in arithmetic as a whole, and synthesis may map it to a memory.

Selects peel dimensions from the outside in. mem[5] is the whole packed element, 32 bits; mem[5][1] then selects along the leftmost packed dimension, an 8-bit row; mem[5][1][0] is one bit. A front end computes these widths during elaboration and type checking, and a mismatch between a select's width and its destination is one of the most common lint warnings.

## Going deeper

IEEE Std 1800-2017, clauses 7.2 (structures) and 7.4 (packed and unpacked arrays); Sutherland, "SystemVerilog for Design", chapters 3 and 5.
