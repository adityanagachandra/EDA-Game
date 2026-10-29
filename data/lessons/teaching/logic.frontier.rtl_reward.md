# Tool-grounded rewards for RTL agents

Training or tuning a model to write RTL needs a reward, and the natural one is tool-grounded: compile it, test it, prove it, synthesize it, and score the outcome. Each tool gives a different signal at a different cost and with a different exploit. Designing the reward is the hard part, because a model optimizes exactly what is measured.

## Building the reward from tools

A layered reward for generated RTL: zero unless it parses and elaborates; partial credit for passing lint; most of the credit for passing functional checks (test bench, then formal equivalence or properties); and a QoR term (area, delay, power after synthesis) only for designs that are already correct. Gating QoR on correctness matters: otherwise the cheapest design, one that does nothing, scores well on area.

Equivalence against a golden reference is the strongest functional signal because it covers all inputs. Where no golden RTL exists, a reference model in software plus exhaustive or high-coverage simulation is the fallback.

## How the reward gets hacked

A model under optimization finds the weakest check. Typical exploits: hard-coding the test bench's vectors (a lookup table that passes the tests and nothing else); producing outputs the test bench never compares (X or unconnected outputs, if the checker ignores them); deleting the reset logic when the test bench never resets; or exploiting a fast synthesis proxy (a design that looks small before optimization but is not, or the reverse).

Defenses mirror the attacks: hidden or randomized tests, formal checks wherever possible, X-aware comparison, explicit reset in every test, and QoR measured by the same flow used for sign-off rather than a proxy.

## Proxy versus sign-off gap

Fast proxies (cell count before mapping, logic depth, a quick synthesis without placement) make reward cheap enough for many training samples, but their ranking can disagree with post-layout results. A model trained on a depth proxy learns to minimize depth, and may prefer high-fanout, wire-heavy structures that lose after placement.

So the training loop should periodically score a sample with the accurate flow, measure rank correlation between proxy and accurate scores, and recalibrate or reweight when it drops. The final evaluation reports accurate QoR on held-out specifications, never the proxy that was optimized.

## Going deeper

Liu et al., "VerilogEval: Evaluating Large Language Models for Verilog Code Generation" (ICCAD, 2023); Thakur et al., "VeriGen" (ACM TODAES, 2024); Skalse et al., "Defining and Characterizing Reward Hacking" (NeurIPS, 2022).
