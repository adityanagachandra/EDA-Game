# Tool-grounded rewards for logic transforms

Training a model to propose logic rewrites with reinforcement learning needs a reward, and the natural reward comes from tools: an equivalence checker says whether the rewrite is correct, structural hashing and a node count say whether it helps. Because the tools are deterministic and exact, the reward is grounded in ground truth rather than in another model's opinion. The risks are equally concrete: a reward computed with a weak check or a sloppy metric will be gamed, and a proxy metric can drift from the quality that signoff measures.

## Scoring a proposed transform

A sound reward gates on correctness and then measures gain: $R = [\text{miter UNSAT}] \times (N_{before} - N_{after}) / N_{before}$, where $N$ is the AND-node count after structural hashing. A rewrite from 1,200 to 1,050 nodes that is proven equivalent earns $150 / 1200 = 0.125$; dividing by the new size instead (0.143) is a common slip. A rewrite that fails the check earns nothing, however small it is.

Both counts must come from the same pipeline, hashed the same way, or the reward measures representation noise. If the design is timing-critical the reward also needs a depth or delay term, since a rewrite that saves nodes by adding logic levels can make the block worse.

## How a reward gets gamed

A policy maximizes whatever the reward actually checks. If equivalence is tested with 1,000 random vectors instead of SAT, the policy learns rewrites that are wrong on rare inputs. If a SAT timeout is scored as "no counterexample found", it learns rewrites that make the miter hard, so checks time out and pay out. If the reward accepts don't cares without verifying them, it learns to change behavior on reachable inputs and call them don't care.

Metric mismatches are gamed too: if the baseline is counted on an unhashed netlist and the result on a hashed one, the policy is credited for hashing it did not do. The defense is the same each time: make the reward use the exact check, the trusted care set and one consistent measurement pipeline, and score UNKNOWN as zero or a penalty.

## From proxy to signoff

AIG node count is a fast proxy. After technology mapping, area and delay are what count, and the two can diverge: an 8% node reduction may yield 1% mapped area. Periodically validating rewards against mapped QoR, and reweighting when they drift, keeps the proxy honest; this gap between proxy reward and signoff QoR is the central risk of tool-grounded RL.

Some parts never leave the loop. The exact equivalence check runs on every accepted rewrite, however good the policy becomes, because correctness cannot be learned to a guarantee. Rewards are computed with fixed solver versions, seeds and limits, so identical outputs earn identical rewards and training runs are reproducible.

## Going deeper

Mishchenko, Chatterjee and Brayton, "DAG-aware AIG rewriting", DAC (2006), for the rewriting moves and node metrics. Skalse, Howe, Krasheninnikov and Krueger, "Defining and characterizing reward hacking", NeurIPS (2022).
