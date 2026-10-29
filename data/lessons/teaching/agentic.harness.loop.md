# The edit-build-test harness loop

The harness is the program that lets a model change an engine safely. It automates the cycle engine developers already run by hand: edit, compile, test, read the failure, try again. The model proposes; the harness writes, builds, tests and reports; a human approves. Every agent-driven engine change in this architecture passes through this loop, so its steps and its stopping rules are the foundation for everything that follows.

## The six steps

A request arrives in plain language: "abut the regular cells of this datapath". The model maps it to the smallest safe change and proposes a small code edit, typically a plug-in at an extension point or a configuration change. The harness writes the file and runs an incremental build, which takes seconds with a compiler cache because only the touched unit recompiles. The harness then runs the skill's golden tests and a QoR comparison against the baseline. If the build or a test fails, the compiler, sanitizer or test output goes back to the model as feedback and the model tries again. When everything passes, the change is shown to a human as a reviewable diff with its evidence: build log, tests passed and the QoR delta.

## Who does what

The division of labor is strict. The model only proposes text: an edit, a configuration, an explanation. The harness owns every side effect: it writes files inside a sandboxed workspace, invokes the pinned compiler, runs tests with resource limits and decides whether the evidence meets the bar. The human owns acceptance. This split means the model never needs, and never gets, the power to merge its own change, push to a shared library or skip a test. A model that claims "the tests pass" proves nothing; the harness's own test run is the only evidence that counts.

## Termination and escalation

A loop needs stopping rules. The harness gives each request a budget of attempts and of wall time. It detects oscillation, when the model alternates between two edits that each fix one failure and cause another, and stops with both failures reported together. It stops early when the model's edit leaves the allowed tier, for example by touching core legalization code while working on a Tier 2 request; that becomes a Tier 3 proposal for a human rather than an automatic retry. When the budget runs out, the harness reports the best attempt and its failures instead of silently giving up, so an engineer can pick it up.

## Going deeper

The agent-first EDA notes in this curriculum, section on the harness loop; the CMake and ccache documentation on incremental builds.
