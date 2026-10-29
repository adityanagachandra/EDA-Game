# Judging an AI autorouter by DRC and SI

Learned and agentic autorouters are judged the same way as any engine change: by deterministic checks against golden references. For board routing those checks are DRC, constraint reports (skew, length, impedance) and SI simulation. Acceptance gates turn those results into a pass or fail per candidate, and differential tests on reference designs prove that a change moved only what it was supposed to move. Both are computations a harness runs automatically on every candidate.

## Acceptance gates

A candidate routing is accepted only if every gate passes: zero DRC violations, worst intra-pair skew within its budget, every length-match group within tolerance, and so on. Gates are conjunctive, so failing one rejects the candidate regardless of how good the rest is. Example with a 5 ps skew gate: candidates with (0 violations, 4.5 ps), (0, 1.5 ps), (2, 1.0 ps) and (0, 6.5 ps) yield two acceptances; the third fails DRC despite the best skew. The gates are computed by tools, the DRC engine and the constraint reporter, not by the router's own estimate, so a router cannot claim success. A golden test suite fixes the reference designs and the expected gate results, so the same candidates always get the same verdict.

## Signal-integrity gates

Geometry checks are not enough for high-speed nets. A routing can be DRC-clean and still put a controlled net over a plane split or through a long via stub. So the harness extracts each critical net and simulates it, comparing the result with targets: impedance within a tolerance band, such as 10% of a 50 ohm target (45 to 55 ohms), insertion loss below budget, eye opening above the receiver mask. Example: simulated impedances of 46.5, 51.5 and 56 ohms against 50 ohms plus or minus 10% pass the first two and fail the third. Fast extraction and quasi-static solvers keep this affordable inside the loop; full-wave simulation is reserved for nets near their limits and for final signoff.

## Differential tests on reference designs

When an agent changes routing policy, a **differential** test on a fixed reference design shows what moved. Route the design with the baseline and with the change, then compare per-net delay, length, via count and violations. A change meant to retune only the DDR group passes if non-DDR nets moved by no more than a small tolerance. Example: with a 0.5 ps tolerance, a SIG net moving from 247.2 to 245.8 ps is an unintended change, while one moving from 239.3 to 239.1 ps is not. Unintended changes are not necessarily bad, but they must be explained before promotion. The same discipline as a timing derate change applies: the proof is a diff that shows only the intended nets moved, by the intended amount.

## Going deeper

The agent-first EDA notes in this curriculum (golden and differential testing). Bogatin, "Signal and Power Integrity Simplified" (3rd ed., 2018), on post-route SI verification.
