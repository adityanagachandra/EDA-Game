# Content-addressed result caching

Flows recompute the same things over and over: a regression reruns unchanged designs, a parameter sweep shares its synthesis step, a developer reruns placement after editing only a report script. A result cache stores each step's output under a key computed from its inputs, so an identical invocation is answered by a lookup. The idea is the same one build systems use for object files. It is worth doing only when lookups are cheap relative to the computation, and correct only when the key covers every input that affects the output.

## Expected time with a cache

Let a lookup cost $L$ (hashing the inputs and probing the store), the computation cost $C$, and the hit rate be $h$. Every invocation pays the lookup; only misses pay the computation. The expected time is
$$E[T] = L + (1 - h)\,C.$$
Writing the result back into the store on a miss adds a little more, which we ignore here.

Example: an extraction step takes $C = 300$ s, the lookup costs $L = 12$ s and the hit rate is 75%. Then $E[T] = 12 + 0.25 \times 300 = 87$ s, against 300 s without a cache. The expectation is a linearity-of-expectation statement over many invocations; any single call takes either 12 s or 312 s. Large flows care about the expectation because they run thousands of steps.

## The break-even hit rate

A cache can lose. A miss costs $L + C$, more than the uncached $C$, so the cache helps on average only when
$$L + (1 - h)\,C < C \quad\Longleftrightarrow\quad h > \frac{L}{C}.$$
Example: hashing a large netlist and probing remote storage costs $L = 20$ s, and the step runs in $C = 160$ s. The break-even hit rate is $20/160 = 12.5\%$. Above it the cache saves time; below it, the cache is overhead.

The ratio $L/C$ explains where caching belongs. Long steps (synthesis, global routing, extraction) have tiny break-even rates and are excellent candidates. Short steps whose inputs are large (a quick report over a full database) may never pay back the cost of hashing their inputs. Cheap keys help: hashing file contents once and reusing the digest is far faster than rehashing on every lookup.

## Keys must hash every input

A content-addressed cache names a result by a hash of everything that determines it: the input files, the step's parameters, the random seed, and the version of the tool binary. Identical keys must mean identical outputs. If any determining input is missing from the key, two different computations share a key and the cache returns a wrong result, silently, which is far worse than a slow run.

Example: runs arrive with (netlist, SDC, tool version, seed) equal to (1, 1, 1, 2), (1, 1, 2, 2), (1, 1, 1, 2), (2, 1, 1, 1). The first, second and fourth keys are new and miss; the third repeats the first and hits, one hit. A key that omitted the tool version would also report the second run as a hit and hand back a result from a different build. Environment variables, PDK files and hidden defaults are the usual omissions. Deterministic tools are a prerequisite: if the same inputs can give different outputs, a cache cannot be trusted at all.

## Going deeper

Mokhov, Mitchell and Peyton Jones, "Build Systems à la Carte" (ICFP 2018), on cloud build caches and constructive traces; the Bazel documentation on remote caching.
