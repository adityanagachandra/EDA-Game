# The cost of selection: manifests versus source

Selecting a skill from manifests is claimed to be nearly free. That claim is arithmetic, and a tool designer should be able to check it: count the tokens the model reads, convert them to seconds, and compare with the engine runs the selection steers. The same arithmetic shows when a skill library has grown too large to put every manifest in front of the model.

## Counting the tokens

A model's cost scales with the tokens it reads. Selecting among $S$ skills from manifests of about $m$ tokens each, plus a design summary of $d$ tokens, reads $S m + d$ tokens. Reading the source instead would cost $S s + d$, where $s$ is the size of one skill's source. With 30 skills, 400-token manifests, a 1,500-token design summary and 60,000 tokens of source per skill, manifest selection reads $30 \times 400 + 1500 = 13{,}500$ tokens, while source selection would read $1{,}801{,}500$, about 133 times more and far beyond any context window. Manifests make selection possible at all, not merely cheaper.

## Selection time against loop time

The time of a model call is roughly the input tokens divided by the prefill rate plus the output tokens divided by the decode rate. Reading 13,500 tokens at 5,000 tokens per second takes 2.7 s; writing a 300-token decision at 50 tokens per second takes 6 s; the call takes 8.7 s. If one iteration of the fast place-route-analyze loop takes 4 minutes, the selection adds $8.7 / 240 = 3.6\%$. Against an accurate signoff run of several hours, it is negligible. The cost is dominated by the engines, as intended. Selection would stop being cheap only if it ran far more often than the engine runs it chooses.

## What fits in the context window

A model has a fixed context window $W$. After reserving $R$ tokens for instructions, the design summary and the answer, at most $\lfloor (W - R) / m \rfloor$ manifests fit. With a 32,000-token window, 6,000 reserved and 400-token manifests, 65 fit. A library of 200 skills does not, so selection becomes two-stage: filter by structured tags (engine type, required inputs, supported row heights) with ordinary code, then let the model choose among the shortlist. Keeping manifests short and their tags structured is what lets a library grow.

## Going deeper

The agent-first EDA notes in this curriculum, section on skills, manifests and selection.
