# Coarse plug-in boundaries and the hot kernel

A seam costs something every time the core crosses it. A virtual call, a call into a separately compiled shared library or a call across a language boundary carries overhead, and it prevents the compiler from inlining and vectorizing across the boundary. Placing a seam per element of a hot loop can make an engine several times slower. Batching the call is the fix, and choosing the batch size is a short computation.

## Overhead per call

Let the core's own work cost $k$ ns per element and a crossing of the plug-in boundary cost $c$ ns. If the core calls the plug-in once per element over $N$ elements, the overhead is $N c$ and the fraction of runtime spent crossing is $c / (k + c)$. A derate lookup with 6 ns of kernel work per arc and a 20 ns crossing spends $20 / 26 = 77\%$ of its time on the boundary. Worse, the per-element call blocks vectorization of the kernel, so the real penalty is often larger than this count suggests.

## Batching the call

If the core passes $B$ elements per call, the number of crossings falls to $\lceil N / B \rceil$, and the overhead becomes $\lceil N/B \rceil c$. For $N = 10^6$ arcs, $c = 20$ ns and $k = 6$ ns, a batch of 256 gives $3907 \times 20 \approx 78{,}000$ ns of overhead against $6 \times 10^6$ ns of kernel work, 1.3%. The plug-in receives an array of contexts and fills an array of results, and inside it the loop over the batch is ordinary, inlinable, vectorizable code. This is what "per batch, not per cell" means in the code-structure rules.

## Choosing a batch size

To keep the overhead below a fraction $x$ of the kernel time, require $c / B \le x k$, so
$$B \ge \frac{c}{x\,k}.$$
With $c = 20$ ns, $k = 6$ ns and $x = 1\%$, $B \ge 20 / 0.06 = 333.3$, so the smallest batch is 334. Batches cannot grow without limit: a larger batch needs more memory for contexts and results, and in an incremental timer the natural batch is the set of arcs at one level of the levelized graph, which may be small. The engine designer chooses a boundary whose natural unit, such as a level, a net or a bin, is large enough.

## Going deeper

Drepper, "What Every Programmer Should Know About Memory" (2007), on the cost of indirection; the OpenSTA source on levelized delay calculation.
