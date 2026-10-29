# Selecting a skill from manifests

Once engines are packaged as skills, choosing one becomes a cheap text decision. The agent reads the manifests, matches them against what it knows about the design and the current objective, and calls the chosen compiled binary. It never opens the source to make the choice. Selection done well puts the right algorithm on each design; selection done badly wastes runs on a skill that does not apply.

## Matching design traits to manifests

The agent first summarizes the design in the terms manifests use: cell regularity (is there a datapath with repeated bit slices?), utilization, macro count, number of clocks and corners, current congestion and timing state, and the objective of this step (wirelength, timing, congestion, runtime). It then reads each candidate manifest's applicability and limits. A datapath placer that advertises "regular bit-sliced logic" fits a design with a 64-bit ALU; a congestion-driven placer fits a design whose fast congestion estimate shows hot spots. Hard requirements come first: a skill that requires a single row height cannot be used on a design with multi-height cells, however well it scores otherwise.

## Scoring the candidates

Among the skills that pass the hard requirements, the agent ranks by fit to the objective. A simple scheme scores each eligible skill as $\sum_i w_i s_i$, where $s_i$ is the skill's advertised or measured strength on objective $i$ and $w_i$ is the current weight of that objective. Suppose the objective weights are 0.5 timing, 0.3 congestion and 0.2 runtime, skill A scores (8, 4, 6) and skill B scores (5, 9, 7). Then A scores $4 + 1.2 + 1.2 = 6.4$ and B scores $2.5 + 2.7 + 1.4 = 6.6$, so B is chosen even though A is stronger on the most heavily weighted objective. A skill that would win but fails a hard requirement is never a candidate.

## When selection should trigger customization

Selection can also conclude that no skill fits. If every eligible candidate scores poorly, or the fast modeling loop shows the chosen skill repeatedly failing on the same metric, the right next step is to customize: tune a parameter, or write a plug-in at the skill's extension point. That decision comes from the measured metrics and their variance in the live fast-skill loop, not from a static rule. Measured strengths should also replace advertised ones as evidence accumulates, because a manifest states intent and the loop measures fact.

## Going deeper

The agent-first EDA notes in this curriculum, sections on skills and on fast and accurate engines.
