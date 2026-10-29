# Serializing a design database

A tool must save its database between flow stages, ship it between machines and reload it quickly. Serialization turns the object tables into bytes, and its design decides file size and load time. Because databases are built on integer handles and coordinates in database units, compact integer encodings pay off directly: variable-length integers, delta encoding of sorted values and zigzag encoding of signed differences. These same encodings appear in many binary formats and checkpoint schemes.

## Variable-length integers

Most stored integers are small: pin indices, layer numbers, short offsets. A fixed 4-byte field wastes space on them. A variable-length integer (the LEB128 scheme) stores 7 payload bits per byte and uses the top bit as a continuation flag. A value below $2^7 = 128$ takes 1 byte, below $2^{14} = 16384$ takes 2, below $2^{21} = 2097152$ takes 3, and below $2^{28}$ takes 4. Example: the values 5, 300, 70000 and 3000000 take $1 + 2 + 3 + 4 = 10$ bytes, against 16 at fixed width. Decoding reads bytes until one has its top bit clear, shifting each 7-bit group into place. The format is self-delimiting, so no length field is needed, but it cannot be indexed randomly: to find the tenth value the reader must decode the first nine, which is why files add offset tables for random access.

## Delta encoding sorted coordinates

Coordinates are large numbers in database units, but sorted coordinates are close to each other. Delta encoding stores the first value and then each difference from the previous one. The x positions of cells in a row, 412000, 412380, 412760, 414280 and 416180, become 412000, 380, 380, 1520 and 1900. As variable-length integers these take $3 + 2 + 2 + 2 + 2 = 11$ bytes, against $5 \times 4 = 20$ bytes fixed, a saving of 9 bytes. Over millions of cells, routing segments or polygon vertices the savings are large, and general compressors do even better on the regular deltas that result. The method depends on order: the writer must sort first, or choose an order such as row by row in which neighbours are close, and the reader must decode sequentially to rebuild absolute values.

## Zigzag encoding signed deltas

Differences are not always positive: polygon vertices go back and forth, and wires turn. A negative number in two's complement has its high bits set, so as an unsigned variable-length integer it would take the maximum number of bytes. Zigzag encoding maps signed integers to unsigned ones so that small magnitudes stay small: $n \mapsto 2n$ for $n \ge 0$ and $n \mapsto -2n - 1$ for $n < 0$. So 0, -1, 1, -2, 2 map to 0, 1, 2, 3, 4. Example: deltas 63, -64, 700 and -8192 map to 126, 127, 1400 and 16383, taking $1 + 1 + 2 + 2 = 6$ bytes. Note the boundaries: 64 maps to 128, which needs 2 bytes, while -64 maps to 127 and needs one. Binary database formats and GDS-like streams use these ideas, alongside versioned headers so that a newer reader can load an older file.

## Going deeper

The Protocol Buffers encoding documentation on varints and zigzag encoding; the OASIS format specification (SEMI P39), whose integer encoding uses the same ideas.
