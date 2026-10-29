# Designing an extension-point ABI

An extension point is only stable if its binary interface is. Plug-ins compiled last month must still load after a core release, plug-ins must not slow the hot loop, and plug-in code must not carry hidden state that breaks determinism or threading. This concept covers the design rules for such an ABI and the arithmetic of the call boundary, which decides whether a plug-in is called per cell or per batch.

## What a frozen extension ABI promises

A frozen ABI promises that a binary built against version 1 keeps loading and working with every later core release of the same major version. That rules out reordering struct fields, changing argument types and removing members. Many systems expose a C ABI even when core and plug-ins are written in C++, because C calling conventions and struct layouts are stable across compilers, while C++ name mangling, vtable layouts and standard-library types vary between compilers and versions. Two habits make evolution possible: begin each extension struct with its size and a version number, so the core knows which fields a plug-in understands, and only append new optional functions or fields at the end. Exceptions should not cross the boundary, since unwinding across modules built by different toolchains is not portable; return error codes.

## Paying for the call boundary

Every call across a seam costs overhead: an indirect call, argument marshalling and, above all, the loss of inlining and vectorization. Suppose a pass evaluates $N = 700{,}000$ cells, the cost work is $w = 26$ ns per cell and each call costs $c = 10$ ns. Calling once per cell pays $N c = 7$ ms of overhead on top of $N w = 18.2$ ms of work. Calling once per batch of $B = 1024$ cells makes $\lceil N/B \rceil = 684$ calls, 6.84 microseconds of overhead. The pass speedup is $\frac{N(w+c)}{Nw + \lceil N/B\rceil c} = \frac{25.2}{18.2068} \approx 1.38$. The real gain is usually larger, because inside a batch the plug-in's own loop can be vectorized, which the per-cell design forbids. The rule follows: coarse calls, per batch, not per cell.

## Plain data and no hidden state

A plug-in should be a pure function of its explicit arguments: read-only input arrays of plain structs, an output array it fills, and no global or static state. Passing a std::vector of internal objects by reference couples the plug-in to layouts that are not ABI-stable; a pointer plus a length to plain structs is. A plug-in that caches the last design in a static variable makes results depend on call history, which breaks determinism, makes golden tests flaky, and races when the core calls it from several threads. A simple diagnostic: call the plug-in on the same batches in a different order, or on several threads. Pure plug-ins give identical results; any difference reveals hidden state. Purity is also what lets the core run a plug-in in a sandboxed process and compare its output to a golden result.

## Going deeper

The agent-first engine notes in this curriculum on seams and extension points; documentation of stable plug-in ABIs such as the C interfaces used by many EDA and compiler projects.
