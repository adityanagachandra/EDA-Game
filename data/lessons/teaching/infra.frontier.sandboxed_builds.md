# Sandboxed plug-in builds

In an agent harness, model-written plug-in code is compiled and run many times an hour. That code is untrusted, the designs it touches are confidential, and the loop must be fast and reproducible. A sandboxed build answers all three: isolated builds with resource limits and no network, modular builds with a compiler cache, and pinned offline toolchains, so every iteration is safe, quick and repeatable.

## The sandbox around a plug-in build

The harness builds and runs each candidate plug-in inside a sandbox: a container or similar isolated environment with no outbound network, so code cannot exfiltrate design data or pull in unreviewed dependencies, and with limits on CPU, memory and wall-clock time, enforced on Linux by mechanisms such as control groups. An infinite loop becomes a time-limit failure and a runaway allocation becomes a memory-limit failure, both reported back to the model instead of harming the host. The core is visible only as read-only headers and prebuilt libraries, so a plug-in edit cannot silently turn into a core edit. The build is modular: the plug-in compiles as its own unit against the stable interface, so an edit rebuilds only the plug-in, not the engine.

## Iteration time with a compiler cache

A compiler cache such as ccache returns stored object files when a unit's preprocessed input and flags are unchanged. With $U$ units, hit rate $h$, hit cost $t_h$, miss cost $t_m$, link time $L$ and test time $T$, one iteration takes $U(h t_h + (1-h) t_m) + L + T$. Example: $U = 20$, $h = 0.85$, $t_h = 0.2$ s, $t_m = 12$ s, $L = 6$ s, $T = 30$ s gives $20(0.17 + 1.8) + 36 = 75.4$ s, so $\lfloor 3600/75.4 \rfloor = 47$ iterations per hour. A cold build costs $20 \times 12 + 36 = 276$ s, so the cache gives a speedup of $276/75.4 = 3.66$. The fixed link and test time caps the gain; past a point, faster golden tests matter more than a better hit rate. The hit rate itself depends on the build's structure: a plug-in that includes a header the core changes often will miss on every unit after each core update, while a plug-in that sees only a small, stable interface header keeps hitting. Narrow interfaces therefore speed up the agent loop as well as making it safer, and modular builds keep the number of units an edit touches small.

## Pinned offline toolchains

A build that fetches the latest compiler can pass on Monday and fail on Tuesday with no code change. The sandbox therefore uses a pinned toolchain: exact compiler, library and build-tool versions, often in a container image referenced by an immutable content digest rather than a moving tag such as latest. Dependency sources come from a local mirror, so builds stay offline and unchanged even if an upstream project disappears. Reproducible builds go one step further: identical source and toolchain should give bit-identical binaries, which requires stripping embedded timestamps and normalizing absolute paths. A promoted skill records the toolchain digest and flags that built it, so the same binary can be rebuilt and audited later, and a regression can be traced to a code change rather than an environment change.

## Going deeper

The agent-first notes in this curriculum on the harness loop and sandboxing; the ccache documentation; the Reproducible Builds project documentation.
