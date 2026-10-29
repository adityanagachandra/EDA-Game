# GDSII and OASIS

OASIS is the newer stream format designed to replace GDSII for very large layouts, especially post-OPC mask data. Both describe the same thing: cells, placements, polygons on numbered layers. OASIS is typically several to many times smaller because it removes the redundancy that GDSII writes into every record. For a verification tool, the format decides parse cost, memory, and how much repetition the reader can recover without analysis.

## Variable-length integers

GDSII writes every coordinate as a fixed 4-byte integer. OASIS writes unsigned integers in 7-bit groups, one group per byte, with the top bit set when another byte follows. A value $v$ with $b$ significant bits takes $\lceil b / 7 \rceil$ bytes, so 100 takes one byte, 1266 (11 bits) two bytes, and 300000 (19 bits) three bytes. Signed integers spend one bit of the first byte on the sign. Because OASIS can also write a vertex as a delta from the previous one, most coordinates in a dense layout are small numbers and cost one or two bytes. The same idea, varint encoding, appears in many binary formats; here it is the single largest source of savings on uncompressed data.

## Storing a regular array

Repetition is the second lever. In GDSII, a $16 \times 22$ via array can be written as 352 flat rectangles at 64 bytes each (22528 bytes), as 352 SREFs of a one-rectangle cell, or as one AREF. The AREF costs AREF 4 + SNAME + COLROW 8 + XY 28 + ENDEL 4 bytes; its XY holds three points (origin, column end, row end). With the five-letter name FILL1 padded to 6, the cell costs $28 + 10 + 64 + 4 = 106$ bytes and the AREF $4 + 10 + 8 + 28 + 4 = 54$, so 160 bytes in total. GDSII can only express repetition of placements, and only as one regular two-dimensional array. OASIS attaches repetitions to placements and to individual shapes, with regular and irregular forms.

## Why OASIS files are smaller

Four mechanisms combine. Variable-length integers and relative coordinates shrink numbers. Modal variables let a record omit fields equal to the previous record's (layer, datatype, width, geometry size): an info byte says which fields are present. Repetitions turn arrays of identical shapes into one record. Name tables store each cell name and text string once and refer to it by index. OASIS also defines optional compressed blocks (CBLOCK, using DEFLATE). None of this changes the geometry, so a converter between the formats is lossless for the shared feature set. GDSII remains common because every tool reads it and its fixed records are simple; OASIS dominates where size hurts most, such as post-OPC data headed to the mask shop.

## Going deeper

SEMI P39, the OASIS specification; the KLayout documentation on OASIS reading and writing.
