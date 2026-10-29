# SRAM organization

Static RAM is the on-chip memory of caches, buffers and register files, and in many chips it occupies more area than logic. It is never built as a logical $W \times B$ array directly: the bits are folded into a physical array whose shape is chosen for speed and area, with a row decoder, column multiplexers and sense amplifiers. Memory compilers generate these macros, and floorplanning, timing and power tools treat them as hard blocks whose geometry follows from the organization.

## The 6T cell and the array

The standard SRAM cell is six transistors: two cross-coupled inverters hold the bit, and two access transistors connect the storage nodes to a complementary pair of bitlines when the row's wordline is raised. A $W$-word by $B$-bit memory stores $W B$ bits, so $6 W B$ transistors in the cell array alone; a 32 KB cache data array (262,144 bits) holds about 1.57 million cell transistors.

The cells sit in a grid of rows and columns. One wordline per row selects a row; every cell in that row drives its column's bitlines. The row decoder (an address decoder) turns part of the address into the one-hot wordline.

## Column multiplexing

A $1024 \times 32$ memory laid out as 1024 rows of 32 cells is tall and thin: long bitlines (slow, high capacitance) and a narrow array. Instead, put $m$ words in each physical row and select among them with a column multiplexer. The array then has

$$\text{rows} = \frac{W}{m}, \qquad \text{columns} = B\, m.$$

With $m = 4$: 256 rows by 128 columns. The address splits accordingly: $\log_2(W/m)$ bits drive the row decoder and $\log_2 m$ bits drive the column mux. Here 8 row-address bits and 2 column-address bits make up the 10-bit word address.

## Choosing the shape

The column-mux factor $m$ is chosen to make the array close to square, which balances wordline length (RC delay across a row) against bitline length (the load each cell must discharge). The array is square when $W/m = B m$, that is $m = \sqrt{W/B}$. For $4096 \times 16$, $m = \sqrt{256} = 16$, giving 256 rows by 256 columns.

Large memories are further split into banks or sub-arrays, each with local decoders and sense amplifiers, because a single huge array would have unacceptably long lines. Memory compilers expose words, bits, mux factor and banking as parameters and emit the layout, a timing model and an abstract view; choosing among their configurations for a given timing and area target is routine design-space exploration.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed., 2010), chapter 12; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003), chapter 12.
