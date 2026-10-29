# DSP and block RAM tiles

Modern FPGAs embed hard blocks among the logic: block RAMs with configurable aspect ratios and DSP blocks with hard multipliers and adders. They are far denser and faster than the same function in LUTs, but they come in fixed sizes and fixed columns. Mapping a memory or a multiplication onto them is an integer packing problem that synthesis and placement must both respect.

## Block RAM shapes

A block RAM tile has a fixed capacity but can be configured to several depth-by-width shapes, for example 32768×1, 16384×2, 8192×4, 4096×9, 2048×18, 1024×36 and 512×72; the 9-, 18-, 36- and 72-bit widths include parity bits. A $D \times W$ memory built from identically shaped tiles needs $\lceil D/d \rceil \cdot \lceil W/w \rceil$ tiles for shape $d \times w$, and the mapper picks the minimum. For $1500 \times 20$: the 2048×18 shape gives $1 \times 2 = 2$ tiles, and 1024×36 gives $2 \times 1 = 2$ as well, while 512×72 gives $3 \times 1 = 3$. A shape that wastes width or depth still consumes a whole tile, which is why memory shapes that match the tiles matter for utilization.

## DSP tiling of wide multipliers

A DSP block multiplies operands of fixed widths, for example 18 by 18 or 27 by 18 bits. A wider product is split into partial products: cut the $m$-bit operand into $\lceil m/p \rceil$ chunks and the $n$-bit operand into $\lceil n/q \rceil$ chunks, and each pair of chunks is one DSP multiply. The two operand orientations can differ. For a $40 \times 33$ unsigned product on $27 \times 17$ multipliers, assigning the 40-bit operand to the 27-bit port gives $2 \times 2 = 4$ blocks, while the other way gives $3 \times 2 = 6$. The partial products are then shifted and summed, often in the DSP blocks' own post-adders and cascade paths. Signed arithmetic loses a bit per chunk to the sign, which changes the counts; this count treats unsigned chunks only.

## Distributed RAM in LUTs

A 6-input LUT is a 64-bit memory, and many FPGAs let it be written at run time, giving distributed (LUT) RAM. A $D \times W$ distributed RAM needs $\lceil D/64 \rceil$ LUTs per bit column, $W\lceil D/64 \rceil$ in all, plus multiplexers to select among the column's LUTs when $D > 64$. Distributed RAM reads asynchronously and sits anywhere in the fabric, so it suits small register files and shallow FIFOs. A $448 \times 1$ RAM needs 7 LUTs; a $16384 \times 32$ RAM would need 8192 LUTs and belongs in block RAM. The trade-off between the two is a capacity and timing decision that synthesis makes from the memory's size, read behavior and the device's remaining resources.

## Going deeper

Kuon, Tessier and Rose (2008), chapter on heterogeneous blocks; the VTR documentation on memory and multiplier architecture modeling.
