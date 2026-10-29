# Agentic signoff and ECO loops

An agentic signoff loop reads the engines' reports, timing and SI violations, IR and EM hot spots, DRC markers, proposes ECOs, runs them through extraction and analysis, and iterates until the checklist is clean, with a human approving the final change. The engines stay the source of truth; the agent's job is to search and to assemble evidence. Designing that loop safely means deciding what the agent reads, what it may change, how it uses fast and accurate engines, and how it avoids gaming its own metrics.

## Reports in, verified ECOs out

The agent reads typed reports, not raw layout: lists of violating paths with slacks and delta delays, victims with their aggressors and coupling, IR-drop maps, EM violations, DRC markers. It proposes a fix, for example moving a victim to a free track or upsizing a holding driver, and the harness applies it, re-extracts the touched region and re-runs timing, SI, EM/IR and DRC there. Only a change that passes these checks is shown to a human, as a reviewable diff with its evidence. The human approves the final ECO; the agent never writes to the tape-out database directly. Engines exposed as tools with typed inputs and outputs, for example over MCP, make this loop possible without the agent reading engine source.

## Guarding against bad loops

Two failure modes dominate. Oscillation: a setup fix breaks hold, the hold fix breaks setup, and the loop cycles. The loop must detect repeated states and either propose joint fixes, such as delaying the short path while speeding the long one, or escalate. Reward hacking: an agent rewarded for fewer violations may learn to add false-path exceptions or relax margins, which removes violations without fixing anything. Constraints and signoff margins define what correct means, so changes to them must be blocked or require explicit human approval. A differential re-extraction and SI report, showing that only the intended nets and their neighbours changed, is the check that a physical fix did not create new victims elsewhere.

## Fast engines, accurate signoff, and the flywheel

Inside the loop the agent uses fast extraction and timing to rank candidates, and recurses to accurate extraction when a net's slack is within its fast-mode error bound: 2 ps of slack with a 5 ps bound is not a pass. Final signoff always runs the accurate engines on the whole design. The loop runs in a sandbox on the customer's site, with CPU, memory and time limits and no network, because layouts and parasitics must not leave the IP boundary. Batches are kept small enough to verify and roll back. When a fix pattern keeps working across projects, such as spacing-based coupling fixes on clock nets, it is validated with golden and differential tests and promoted into a library skill; the know-how flows back, the customer's design does not.

## Going deeper

The agent-first EDA domain of this curriculum, on harness loops, change tiers and the skill flywheel; the OpenROAD documentation on its Tcl and Python interfaces, which expose engines as callable tools.
