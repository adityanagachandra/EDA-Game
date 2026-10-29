# Line buffers and data reuse

Stencil computations, such as convolutions and filters over images, read each input many times. Fetching every window from external memory would need $K^2$ reads per output. Line buffers keep the most recent rows on chip, so each pixel is read from external memory once and reused from local storage. Sizing them is a small computation with large consequences for on-chip memory.

## How many pixels to keep

A $K \times K$ stencil over an image streamed in raster order, one pixel per cycle, needs the current pixel and the $K - 1$ rows above it. Keeping those rows on chip, $(K - 1) W$ pixels for an image of width $W$, lets each new pixel complete a window column with no external re-reads; the newest row streams in directly. A $5 \times 5$ convolution on a 2048-pixel-wide image keeps $4 \times 2048 = 8192$ pixels. A small window register of $K \times K$ pixels, shifted every cycle, holds the current window itself. External traffic drops from $K^2$ reads per output to one, which turns a bandwidth-bound kernel into a compute-bound one. The line buffers are an array accessed with a cyclic pattern, so each row is usually its own memory bank to give one read per row per cycle.

## Block RAMs for the buffers

Each buffered row is typically mapped to its own block RAM so that all $K - 1$ rows can be read in the same cycle. A row of $W$ pixels of $b$ bits needs $W b$ bits. Treating a block RAM as 18,432 bits of capacity, a row needs $\lceil W b / 18432 \rceil$ blocks, and the line buffers need $K - 1$ times that. A 3x3 stencil on 1024-pixel rows of 32-bit pixels needs $\lceil 32768 / 18432 \rceil = 2$ blocks per row, 4 in all. Real block RAMs come in fixed aspect ratios, so a narrow, deep row may waste capacity and the true count can be higher than this estimate. Wide pixels, wide images and large kernels multiply quickly, and line buffers can dominate the on-chip memory budget of an image pipeline.

## Latency to the first output

A line-buffered stencil cannot produce its first valid output until a full window is available. If pixel $(r, c)$ arrives in cycle $r W + c + 1$, the first window whose bottom-right pixel is $(K - 1, K - 1)$ is complete when that pixel arrives, in cycle $(K - 1) W + K$. For a $5 \times 5$ stencil on a 1024-wide image that is cycle 4101. This start-up latency is a fixed cost per frame, small relative to a full frame but important when frames are small or when the stencil sits inside a feedback loop. Border handling adds policy choices, such as padding, replicating edges or skipping invalid outputs, that change the output count and must match the reference model exactly, or co-simulation will flag differences at the image edges.

## Going deeper

Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018), chapter on image processing; Bailey, "Design for Embedded Image Processing on FPGAs" (2011).
