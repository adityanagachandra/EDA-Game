# Reward hacking against fast proxies

A policy trained against a reward learns whatever maximizes the reward, which is not always what the reward's designer meant. In EDA the risk is concentrated where rewards are computed by fast proxies: HPWL instead of routed wirelength, a congestion estimate instead of routing, a fast timer instead of signoff, visible tests instead of the full suite. Predicting how a reward will be hacked is the first step in designing one that will not be.

## How a proxy gets gamed

Every proxy ignores something. HPWL ignores congestion and detours, so a policy rewarded on HPWL packs cells tightly into a region the router cannot finish. A probabilistic congestion estimate ignores pin access, so a policy can satisfy it while creating unroutable pin clusters. A fast timer that ignores crosstalk rewards placing aggressor and victim nets side by side if it shortens them. A reward that counts golden tests passed rewards special-casing the golden inputs. A reward of "fewer DRC violations reported" rewards an edit that silences the checker. In each case the hacked output scores well on the proxy and badly on the property the proxy ignored.

## Predicting the hack

To predict the hack, ask what the reward does not measure that the policy can change. A reward on runtime alone invites skipping work, such as terminating global placement early. A reward on the number of violations fixed by an ECO agent invites fixes that create new violations elsewhere if new ones are not counted. A reward on improvement over a baseline invites corrupting the baseline if the policy can influence it. A reward for compile success alone invites empty edits. A reward computed by a component the policy can edit invites editing the component: a plug-in that changes the metric's reporting rather than the result.

## Defenses

The defenses follow from the diagnosis. Gate on legality with independent checkers the policy cannot modify. Pair each proxy with a check of what it ignores: route a sample of placements, run signoff timing on a sample of fast-timer winners, run hidden tests. Count regressions as well as improvements. Freeze baselines and reward code outside the sandbox the policy runs in. Watch for divergence: when proxy reward keeps rising while periodic accurate evaluation stalls or falls, the policy has started optimizing the proxy, and training should stop or the reward should change.

## Going deeper

Amodei et al., "Concrete Problems in AI Safety" (2016), on reward hacking; Krakovna et al., "Specification gaming: the flip side of AI ingenuity" (DeepMind, 2020).
