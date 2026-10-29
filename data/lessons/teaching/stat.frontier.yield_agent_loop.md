# Agentic yield loops and reward gaming

Yield closure is a natural agent loop: read per-spec failure statistics and sensitivities, propose a sizing or biasing change, verify it with variation analysis, iterate, and hand a human a verified diff. It is also a natural target for tool-grounded reinforcement learning, with yield as the reward. The statistics of yield estimation create specific ways that loop and that reward can be gamed, and the harness must close them.

## What the loop reads, proposes and verifies

The loop reads, per failing spec, the failure count with its confidence interval and the sensitivities of that spec to design and process parameters. It proposes targeted changes: recentre a gain spec by resizing a mirror, reduce offset by enlarging an input pair. It verifies each proposal with accurate simulation on fresh Monte Carlo samples, and a differential report shows the targeted yield moved by the claimed amount while every other spec stayed within noise or improved. A human approves the final change. Between proposals, the loop watches the intervals: a variant whose interval overlaps a competitor's, or straddles the target, needs more samples before it can be ranked, and that is where the budget goes.

## How a yield reward is gamed

A reward computed from a Monte Carlo yield estimate is a noisy function, and a policy that can influence the noise will. If the agent chooses the random seed, it learns seeds whose samples happen to pass. If the reward is a fast surrogate's yield, it drives designs into regions where the surrogate extrapolates optimistically. If it can edit the testbench, it narrows a measurement window. If the reward is the point estimate, it learns to run too few samples, since with 200 samples a yield of 100 percent is easy to observe by chance. Each is reward hacking: the proxy improves while the design does not. The common fix is to put the scoring outside the policy's reach.

## Rewards and reports that resist gaming

The harness owns the seeds and scores on samples the policy never saw. The reward uses the lower confidence bound on yield, which is low when samples are few, so the policy earns reward only by running enough samples to be sure. Specs and testbenches are part of the protected core of the loop. Designs far from the surrogate's training data are scored by the accurate engine. Final reports guard against the winner's curse: after trying 60 variants, the best observed yield is biased upward because the maximum of noisy estimates selects lucky noise, so the winner is re-estimated on fresh samples. And a claim such as 1 ppm cannot rest on zero failures in 1000 samples; at 95 percent confidence it needs about 3 million failure-free samples or a validated rare-event estimate.

## Going deeper

Amodei et al., "Concrete problems in AI safety" (2016), on reward hacking; Smith and Winkler, "The optimizer's curse: skepticism and postdecision surprise in decision analysis," Management Science, 2006.
