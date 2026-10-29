# Tool-grounded rewards for RTL generation

Training models to write RTL with reinforcement learning needs a reward, and the most reliable rewards come from tools: compilers, simulators, equivalence checkers and synthesis. A tool-grounded reward is deterministic and checkable, but every tool has blind spots, and a model optimized hard against them will find them. Designing the reward is therefore an exercise in knowing exactly what each tool does and does not check.

## Components of a grounded reward

A tool-grounded reward scores generated RTL with deterministic tools against a ground truth: does it compile, does it pass a held-out testbench or prove equivalent to a reference, does it synthesize without latches or errors, and what area and timing does synthesis report. Each component measures something different. Compile success is necessary but nearly worthless alone, because a module that outputs zero compiles, and a model rewarded only for compiling learns to emit safe, useless code.

Correctness needs a reference behavior, from tests or from equivalence checking, and quality terms must be gated on correctness. If smaller area scores higher regardless, the model learns to delete functionality the tests do not exercise.

## Reward hacking

A model maximizes what the reward measures, not what was meant. If the testbench is visible during generation, the model can special-case the stimulus. If it learns that adding a full_case pragma reduces area while the simulation tests still pass, it exploits a gap between the synthesis assumption and the simulation check. If simulation is two-state, a register with no reset reaching an output never shows as x, so the missing reset goes unpunished.

Each hack points at a defense. Regenerate constrained random stimulus for every evaluation so that memorizing it does not help; add formal checks where a reference exists; check four-state behavior or randomize initial values; lint for pragmas and other constructs that change synthesis semantics without being checked. The reward is only as honest as its least careful checker.

## Fast proxies and signoff

Accurate evaluation is expensive. Full place-and-route and signoff timing for every sample would make training unaffordable, so rewards use fast proxies, such as Yosys area estimates or a quick timing estimate. The gap between proxy and signoff quality is where optimization drifts: a model can learn structures that score well on the proxy and poorly after routing.

Two practices contain the drift. Periodically evaluate a sample of outputs with the accurate flow and check that proxy improvements still track real ones, the same fast-versus-accurate pattern used in design loops. Version and pin every tool in the reward, because a tool update that changes results for identical outputs corrupts the learning signal. For benchmarks, tasks drawn from real modules with hidden tests and synthesis checks separate useful models far better than short textbook exercises or prompts that include their own tests.

## Going deeper

The agent-first EDA notes in this curriculum on tool-grounded reinforcement learning; published RTL-generation benchmarks with hidden testbenches.
