# Tool-grounded rewards from slack

Training a model to fix timing needs a reward, and a static timer can provide one: deterministic, cheap enough to run per episode, and grounded in the same numbers signoff uses. Designing that reward means choosing what to reward, what to penalize, and which engine computes it, because a model optimizes the reward it is given, including its loopholes.

## A slack-based reward

A typical reward combines worst and total negative slack with a cost term: $R = a\,\min(WNS, 0) + b\,TNS - g\,\Delta A$, where $\Delta A$ is the area increase in percent. Clipping WNS at zero stops the model from chasing positive slack it does not need. With $a = 2$, $b = 0.1$, $g = 2$, WNS -45 ps, TNS -855 ps and 0.8 percent more area, $R = -90 - 85.5 - 1.6 = -177.1$. The weights encode priorities: a large $a$ focuses on the worst path, a large $b$ on breadth, and $g$ prevents the model from buying slack with unlimited upsizing. Power, hold and design-rule terms are added the same way.

## Computing it from endpoint slacks

The reward must come from the engine's numbers, not from the model's own claims. Given the signoff endpoint slacks after an edit, compute WNS as the minimum and TNS as the sum of the negative values, then the reward. For slacks -34, 24, -32, -33, 9, 24 and -32 ps with $a = 1$, $b = 0.01$, $g = 5$ and 2.1 percent more area, WNS is -34 ps, TNS is -131 ps and $R = -34 - 1.31 - 10.5 = -45.81$. Every episode is timed with the same constraints, corners and engine version, so rewards are comparable across the run. A reward from a fast proxy timer is cheaper but opens a gap to signoff quality that the model will exploit.

## Reward hacking and the grounded check

If the model can edit constraints, it can raise its reward without improving the chip: adding set_false_path on the worst endpoints removes them from the report. Its own timer then shows a better WNS and TNS. The grounded reward re-times the result with the original, trusted constraints, so the hidden endpoints count again. With the two worst endpoints at -32 and -10 ps hidden, a proxy reward of -30.6 becomes a grounded reward of -60.8. Good reward design restricts the action space to legitimate changes, recomputes every metric from trusted inputs, scores at signoff accuracy for final evaluation, and audits high-reward episodes, since they are where hacks appear first.
