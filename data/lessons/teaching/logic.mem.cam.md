# Content-addressable memories

A content-addressable memory (CAM) inverts the usual interface: give it data and it returns the addresses where that data is stored, comparing every entry in parallel. Ternary CAMs (TCAMs) add a don't-care state per bit. They implement translation lookaside buffers, fully associative cache tags, and the longest-prefix-match tables of network routers, at a high price in area and power.

## Parallel match

Every CAM entry has a comparator on every bit. A search drives the key onto search lines running through all entries; each entry's match line stays high only if every bit matches. A binary CAM with $E$ entries of $W$ bits performs $E W$ bit comparisons per search, in one cycle. That parallelism is the point, and also the cost: every search charges every match line, so CAM search energy grows with the whole array.

A ternary entry stores 0, 1 or X (don't care) per bit, using two storage bits per cell, so a TCAM with $E$ entries of $W$ ternary bits stores $2 E W$ bits. An entry matches a key when every non-X bit equals the key bit. For key 1011, entry 1X11 matches, entry 10X0 does not (its last bit is 0), and entry XXXX matches anything.

## Priority among multiple matches

In a TCAM several entries can match. A priority encoder on the match lines picks one, conventionally the lowest-indexed match, and the entries are stored in priority order. For entries 0: 0XX1, 1: 10X1, 2: 1XX1, 3: XXXX with key 1011, entries 1, 2 and 3 match and the result is entry 1.

Routers use exactly this for longest-prefix match: store routes sorted by prefix length, longest first, with the unspecified low bits as X, so the first match is the most specific route, and a catch-all default route at the end matches everything. Inserting a route must preserve the order, which is why TCAM table management is a software problem in its own right.

## Cost and alternatives

A CAM cell needs storage plus comparison logic, so it is several times larger than an SRAM cell, and a TCAM cell larger again; search energy scales with entries times width. That limits CAMs to small, latency-critical structures: TLBs with tens to a few thousand entries, associative tag arrays, and router tables.

The alternative is hashing into an SRAM, which is denser and lower power but has variable latency and collisions. Set-associative caches are a middle point: a few ways compared in parallel, selected by an index. Choosing among CAM, hashed SRAM and set-associative structures is a classic microarchitecture trade-off, and memory compilers offer CAM macros alongside SRAM for that reason.

## Going deeper

Pagiamtzis and Sheikholeslami, "Content-Addressable Memory (CAM) Circuits and Architectures: A Tutorial and Survey" (IEEE JSSC, 2006).
