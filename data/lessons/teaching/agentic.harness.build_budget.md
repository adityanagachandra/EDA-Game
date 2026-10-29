# Incremental builds and compiler-cache budgets

The harness loop is only as fast as its build. If every model edit costs ten minutes of compilation, an agent gets a handful of attempts a day; if it costs ten seconds, it gets hundreds. Build time is therefore a design parameter of an agent-first engine, controlled by how the code is split into compile units, how dependencies are structured and how well a compiler cache hits.

## Compile units an edit touches

A C++ build recompiles every translation unit whose source or included headers changed, then relinks. An edit to a `.cpp` file of a plug-in rebuilds one unit; an edit to a widely included header rebuilds every unit that includes it. Suppose an edit touches a header included by 14 units, each taking 6 s to compile, and linking takes 4 s: the build costs $14 \times 6 + 4 = 88$ s. The same change made in the plug-in's own source file costs $6 + 4 = 10$ s. This is why the code-structure rules say one concern per translation unit, each policy in its own small shared library, and a cap on how many compile units any one change can touch.

## Expected build time with a compiler cache

A fresh sandbox often starts with no object files, so the harness rebuilds everything, and a compiler cache such as ccache serves unchanged units from earlier builds. Units whose inputs changed always miss. If a build has $n$ units of which $c$ changed, the cache hit rate on the others is $h$, a hit costs $t_h$ and a miss costs $t_m$, the expected time is
$$E[T] = c\,t_m + (n - c)\,\big(h\,t_h + (1 - h)\,t_m\big) + t_{link}.$$
With $n = 40$, $c = 2$, $h = 0.9$, $t_h = 0.2$ s, $t_m = 6$ s and a 3 s link: $12 + 38 \times (0.18 + 0.6) + 3 = 44.64$ s. At $h = 0.5$ it rises to $12 + 38 \times 3.1 + 3 = 132.8$ s. The hit rate dominates.

## Attempts per hour

The loop iteration is model time plus build plus tests. With 20 s of model time, a 45 s build and 35 s of golden tests, one attempt takes 100 s and the harness can run 36 attempts an hour. Cutting the build to 10 s by isolating the plug-in raises this to 55. Cache hit rate falls when compiler flags, toolchain versions or include paths vary between builds, so a pinned toolchain helps speed as well as reproducibility.

## Going deeper

The ccache manual, on how cache keys are computed; the CMake documentation on target-based builds.
