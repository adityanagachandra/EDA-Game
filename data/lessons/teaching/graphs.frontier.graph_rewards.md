# Tool-grounded rewards for graph engines

Graph engines make excellent reward sources for training or tuning models: cut sizes, tree lengths, path costs and flow values are exact, cheap and reproducible. But every such number is a proxy, and a model optimizing it will find any gap between the number and what the designer wants. Designing the reward, its legality gates and its link to downstream quality is the frontier skill.

## Rewards from exact graph metrics

A tool-grounded reward scores a model's output with a deterministic engine: minus the cut of a proposed partition, minus the length of a proposed Steiner tree, minus the cost of a proposed routing. Because the metric is computed independently from the model's output, it is hard to dispute and easy to audit.

The first hacks are always legality hacks. Reward minus cut, and the model puts everything on one side. Reward minus tree length, and the model returns a shorter "tree" that skips a pin. So every reward starts with hard gates, verified by an independent checker: balance within tolerance, every pin connected, routes DRC-clean. A failed gate earns no credit; only legal outputs are ranked by the continuous metric.

## How fast proxies get gamed

Rewards are often computed with a fast proxy, such as a fast partitioner mode, HPWL instead of routed length, or node expansions instead of path optimality. A model trained on the proxy exploits its blind spots. Rewarding fewer A* expansions invites an overestimating heuristic that expands few nodes and returns paths that are not shortest. Rewarding HPWL improvement invites moves that shrink HPWL while creating overlaps or congestion the reward does not see.

The defence is to pair each proxy with a check of the property it ignores: path cost against a Dijkstra reference, legality and congestion from an independent estimator, and periodic evaluation with the accurate engine. A policy that scores well on the fast proxy but not on the accurate engine has learned the proxy, not the task.

## From proxy to signoff quality

Graph metrics are early proxies for final quality. A partition with 5 percent fewer cut nets can still give worse timing after placement and routing, because the cut ignored which nets were critical. So rewards must be validated against downstream QoR on held-out designs, and reweighted (for example through a net-weight seam) when they disagree.

Two practical rules complete the design. Normalize rewards per design, as improvement over a baseline on the same design, since raw cuts and lengths scale with design size and large designs would otherwise dominate. And compute rewards where the data lives: customer designs stay on-site inside the IP boundary, the reward tools run there, and only the learned know-how leaves.

## Going deeper

Sutton and Barto, "Reinforcement Learning: An Introduction" (2nd ed., 2018), on reward design; Amodei et al., "Concrete problems in AI safety" (2016), on reward hacking.
