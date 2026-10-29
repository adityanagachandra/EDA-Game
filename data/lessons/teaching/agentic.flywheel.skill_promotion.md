# The skill flywheel

The skill flywheel is how an agent-first EDA system gets better with use. A skill discovered while working on one design is iterated quickly as a plug-in, validated, and promoted into the compiled library, where every agent on every design can select it. Each turn of the wheel adds a reusable skill and teaches the selection model when to use it. The customer's design stays where it was.

## Identify, iterate, validate, promote

The cycle has four steps. **Identify**: someone notices a behavior that would help: the model watching the fast loop, an agent, a field engineer or a customer. For example, abutting the cells of each bit slice in a datapath reduces wirelength and improves timing on regular designs. **Iterate fast**: the behavior is written and tested as a plug-in loaded at runtime (through `dlopen` or the platform equivalent), so each attempt rebuilds only the plug-in in seconds. **Validate**: the plug-in passes golden and differential tests and its QoR gain holds across seeds. **Promote**: it is compiled into a static library skill with a manifest, golden tests and a version, and becomes a first-class member of the library.

## From plug-in to library skill

Promotion changes how the skill is built and who may use it. While iterating, a dynamically loaded plug-in is convenient: it can be swapped without restarting, and a broken one is simply unloaded. Once promoted, the skill is compiled into the engine's library, which lets the compiler optimize across it and lets every deployment select it at full speed with no runtime loading. Promotion also requires what iteration did not: a manifest that states the skill's applicability and limits, golden tests that can ship (built on designs the vendor may hold), and review by the engine's owners. A skill that works on one customer design but cannot be validated elsewhere stays a local plug-in.

## What flows back

A promoted skill lands in two places. The shared library makes it selectable by agents on every design. The domain model learns the skill and when to pick it: which design traits, objectives and loop metrics predicted that it would help. What flows back is the skill's code, its manifest, its tests and the know-how about its use. The customer's design, constraints and process data never flow back. This is what makes the flywheel acceptable to customers and valuable to the vendor: each site's discoveries improve the library without exposing anyone's design.

## Going deeper

The agent-first EDA notes in this curriculum, section on the skill flywheel and the brain model; the POSIX `dlopen` documentation.
