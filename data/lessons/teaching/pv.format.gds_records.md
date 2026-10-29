# GDSII records and stream structure

Physical verification starts from the layout stream that will become masks. GDSII, the oldest and most widely exchanged stream format, is a flat sequence of binary records that describe libraries, structures (cells) and the elements inside them. A verification tool reads nothing else: no nets, no device types, only layers, numbers and coordinates. Knowing exactly what the records hold tells a tool designer what must be reconstructed later (connectivity, devices) and what the file size and parse cost will be.

## Records and element sizes

Every GDSII record starts with a 4-byte header: a 2-byte total length, a 1-byte record type and a 1-byte data type. A polygon is a BOUNDARY element written as five records: BOUNDARY (header only, 4 bytes), LAYER and DATATYPE (one 2-byte integer each, 6 bytes each), XY (4-byte x and y integers per point) and ENDEL (4 bytes). The XY list closes the ring by repeating the first vertex, so a polygon with $n$ vertices stores $n + 1$ points and costs $4 + 6 + 6 + (4 + 8(n+1)) + 4 = 32 + 8n$ bytes. A rectangle ($n = 4$) is 64 bytes. A PATH element (a wire given by its centerline) adds a WIDTH record (8 bytes) and does not close its point list, so $m$ points cost $32 + 8m$ bytes. Layers and datatypes are only numbers: the meaning of layer 68, datatype 20 lives in the technology's layer map, not in the file.

## Counting records in a structure

A structure is framed by BGNSTR, STRNAME and ENDSTR: three records around its elements. A BOUNDARY contributes 5 records, an SREF (a placement of another structure) 4 records (SREF, SNAME, XY, ENDEL), and a TEXT label 6 (TEXT, LAYER, TEXTTYPE, XY, STRING, ENDEL), ignoring optional records such as STRANS for rotation or mirroring. A structure with 19 boundaries, 2 SREFs and 4 labels therefore has $3 + 5 \cdot 19 + 4 \cdot 2 + 6 \cdot 4 = 130$ records. A reader is a loop over records with a state machine that knows which records may follow which. Text labels matter more than their size suggests: they are the only way a stream carries net and pin names, and LVS depends on them.

## The byte size of a structure

BGNSTR carries twelve 2-byte integers (modification and access dates), so it is 28 bytes. STRNAME and SNAME carry an ASCII name padded with a null byte to an even length. A structure INV with 3 rectangles and 2 SREFs of ADC_BIAS costs $28 + (4 + 4) + 3 \cdot 64 + 2 \cdot (4 + (4 + 8) + 12 + 4) + 4 = 296$ bytes: the name INV pads to 4 bytes and ADC_BIAS is already even at 8. These fixed costs explain why flat GDSII files of large designs reach many gigabytes: every rectangle pays 64 bytes, and nothing is shared between records. They also make GDSII trivially seekable and easy to parse in a streaming pass, which is one reason it survives.

## Going deeper

The Calma "GDSII Stream Format Manual" (Release 6.0); the KLayout documentation on stream formats.
