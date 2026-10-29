# Build systems and incremental compilation

A modern EDA tool is millions of lines of C++ across thousands of files, and an engineer, or an agent in a harness loop, rebuilds it many times a day. A build system turns a source edit into the smallest set of recompilations and links that produce a correct binary. It does so with the same machinery as incremental computation: a dependency graph from sources and headers to object files to libraries, and invalidation that follows the graph. Header structure, more than compiler speed, decides how long an edit takes to test.

## Header changes recompile their dependents

In C++, each source file plus everything it includes is a translation unit, compiled to one object file. A translation unit depends on every header it includes directly or through other headers. Editing a header therefore invalidates every translation unit that reaches it in the include graph: a reverse reachability question.

Example: geom.h is included by db.h; net.h includes db.h; place.cpp includes net.h, route.cpp includes db.h, sta.cpp includes timer.h, and gui.cpp includes geom.h. An edit to geom.h recompiles place.cpp (through net.h and db.h), route.cpp (through db.h) and gui.cpp (directly): 3 translation units. An edit to timer.h recompiles only sta.cpp. Build tools such as Make with compiler-generated dependency files, or Ninja under CMake, record these edges after each compile, so the next build knows exactly which objects are stale.

## Parallel compile and the link step

Stale translation units are independent, so the build compiles them in parallel, up to $p$ jobs at once, and then links. If $k$ units each take $c$ seconds, the compile phase needs $\lceil k/p \rceil$ waves and the build takes $\lceil k/p \rceil\, c + L$, with $L$ the link time.

Example: an edit to a low-level header dirties 7 units of 20 s each, the build runs 4 jobs and the link takes 12 s: $\lceil 7/4 \rceil \times 20 + 12 = 52$ s. A one-line change in a .cpp file costs one compile and the link, 32 s. For small edits the link dominates, which is why large tools split into shared libraries or use fast linkers, and why compiler caches such as ccache, which return a stored object when the preprocessed input is unchanged, pay off when switching branches.

## Finding the most expensive header

The cost of touching a header is the number of translation units that reach it, so headers near the bottom of the include graph (core types, geometry, logging) are the expensive ones. A useful audit computes, for every header, its count of dependent translation units and looks at the largest.

In the example, geom.h has 3 dependents, db.h has 2, net.h and timer.h have 1 each: geom.h is the hot spot. Remedies are structural. Forward declarations let a header mention a type without including its definition. Splitting a large header separates rarely changing interfaces from frequently changing details. Keeping implementation out of headers, rather than in inline functions and templates, stops internal edits from rippling outward. These are the same narrow-interface rules that make an engine easy to extend: a plug-in that compiles against a small, stable header rebuilds in seconds when only the plug-in changes.

## Going deeper

Lakos, "Large-Scale C++ Software Design" (1996) and "Large-Scale C++, Volume I" (2019); the Ninja and CMake documentation; Mokhov, Mitchell and Peyton Jones, "Build Systems à la Carte" (ICFP 2018).
