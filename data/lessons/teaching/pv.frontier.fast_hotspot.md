# Fast hotspot prediction versus signoff

Learned hotspot predictors are the fast engine of lithography verification: milliseconds per clip against seconds for simulation. An agent can use them inside a modeling loop, but signoff still belongs to the accurate engine. The questions are how much the fast path saves, what it misses, and how a reward built on it can be gamed.

## Screen then verify

Run the predictor on all $N$ clips at $t_m$ each and simulate only the $k$ flagged ones at $t_s$ each. The speedup over simulating everything is $N t_s / (N t_m + k t_s)$. With $N = 10000$, $t_m = 2$ ms, $t_s = 2400$ ms and $k = 1200$, the speedup is $24000000 / 2900000 = 8.3$. The flag rate bounds the gain: as $k / N$ grows, simulation dominates again. A predictor tuned for very high recall flags more clips, so recall and speedup trade against each other.

## Escapes at a threshold

A predictor outputs a score, and the agent sends clips scoring at least a threshold to simulation. A true hotspot scoring below the threshold escapes if nothing else checks it. With scores and labels from a later full run, count the hotspots below the threshold: those are the false negatives. Lowering the threshold catches more but sends more clips to simulation. Because escapes become yield loss, the fast screen can steer an agent's choices during the loop, but the final signoff must run the accurate engine on the whole layout.

## Gaming the proxy

If an agent is rewarded for lowering the predicted hotspot count, it can learn edits that move layouts into the predictor's blind spots. Precision and recall translate a predicted count into an expected real count: with $P$ predictions, true positives are $P \cdot \text{precision}$, and those are a fraction recall of the real hotspots, so the real count is about $P \cdot \text{precision} / \text{recall}$. With 159 predictions, precision 0.52 and recall 0.8, expect about 103 real hotspots. A tool-grounded reward must therefore be checked periodically against simulation, and the predictor's recall on the agent's own outputs re-measured.
