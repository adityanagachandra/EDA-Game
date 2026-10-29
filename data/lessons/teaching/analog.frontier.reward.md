# Tool-grounded rewards for analog sizing

Reinforcement learning for analog sizing scores a policy's proposals with simulation: the tool provides the reward. Designing that reward decides what the policy learns. A good reward is shaped so that partial progress counts, honest about feasibility, and computed by the accurate engine at the end, because a policy trained against a fast surrogate will learn the surrogate's errors.

## Shaping the infeasible region

A pass-or-fail reward gives a policy no gradient: almost every early proposal fails and scores the same. A shaped reward measures how far a design is from feasibility. One common form sums normalized shortfalls over failing specs and negates it: $r = -\sum_i \max(0, (s_i - m_i)/s_i)$ for lower-bound specs $s_i$ and measurements $m_i$. With specs gain 70 dB, GBW 50 MHz and PM 60 degrees, a design measuring 66 dB, 45 MHz and 62 degrees scores $-(4/70 + 5/50) = -0.1571$. Normalization puts different units on one scale, and ignoring power while infeasible keeps the policy focused on meeting specs first.

## Rewarding the objective once feasible

Once every spec passes, the reward should reflect the objective. With power as the objective and a reference $P_{ref}$, $r = 1 - P/P_{ref}$ gives a positive reward that grows as power falls: 0.6 mW against a 1.5 mW reference scores 0.6. Because every feasible design scores above every infeasible one (positive versus negative), the policy is never tempted to trade a spec violation for power. The choice of $P_{ref}$ and the normalization constants are reward hyperparameters; they change what the policy prefers, so they belong in the configuration and should be fixed before training and comparison, not tuned on the test circuits.

## Reward hacking and the surrogate gap

Training against SPICE is expensive, so policies are often scored by a fast surrogate. The policy then optimizes the surrogate, and its favourite designs cluster where the surrogate is too optimistic. If the surrogate reports a passing design with reward 0.53 and SPICE finds GBW 20% short with reward -0.23, the gap of 0.76 is reward the policy earned by exploiting model error. Other hacks are mechanical: proposals that make the simulator fail to converge, if failures are scored leniently, or measurements that are wrong for unusual waveforms. Defences are to score non-convergence as the worst reward, retrain or re-validate the surrogate on the policy's own proposals, check a sample of rewards against SPICE throughout training, and report final results only from SPICE across corners.

## Going deeper

The agent_first_eda lesson on tool-grounded reinforcement learning; Settaluri et al., "AutoCkt: deep reinforcement learning of analog circuit designs" (DATE, 2020).
