# EDA problems as Markov decision processes

Reinforcement learning applies when a design task is a sequence of decisions whose quality shows up later: place one macro, then the next; apply one synthesis transform, then the next. Casting the task as a Markov decision process means choosing states, actions, rewards and a discount, and those choices decide whether the agent can learn anything at all.

## Returns and discounting

An agent at step $t$ receives reward $r_t$, and its return from step 0 is $G_0 = \sum_t \gamma^t r_t$ with discount $\gamma \le 1$. For a synthesis agent whose rewards after four commands are 1, -4, 2 and 2 (node reductions, in hundreds), $\gamma = 0.9$ gives $1 - 3.6 + 1.62 + 1.458 = 0.478$. The discount says how much later rewards count. In an EDA episode, later steps often matter as much as early ones, because the final QoR is what is measured, so a small discount can make the agent short-sighted and favor transforms that help now and hurt the final result.

## Sparse terminal rewards

Many EDA objectives are known only at the end: a placement's wirelength and congestion are evaluated after all macros are placed. With reward 0 at every step and $-c$ at the last of $M$ steps, the return from step 0 is $-c\gamma^{M-1}$. With $M = 100$ and $\gamma = 0.99$, the terminal reward is scaled by $0.99^{99} = 0.37$; with $\gamma = 0.95$ by 0.006, so early macros receive almost no signal. Placement formulations therefore often use $\gamma = 1$, and some add intermediate rewards, such as partial wirelength, which shapes learning but must not change which final placements are best.

## Value estimates and TD error

A critic learns $V(s)$, the expected return from state $s$, to reduce the variance of policy updates. After one transition the temporal-difference error is $\delta = r + \gamma V(s') - V(s)$. With $V(s) = 2$, an intermediate density penalty $r = -1$, $\gamma = 0.9$ and $V(s') = 4$, $\delta = -1 + 3.6 - 2 = 0.6$: the move went better than expected, and actor-critic methods raise the probability of that action. In EDA the state is large (a partial placement plus a netlist), so $V$ is a neural network, often on top of a GNN encoder, and its errors become noise in the policy gradient.

## Going deeper

Sutton and Barto, "Reinforcement Learning: An Introduction" (2nd ed., 2018), chapters 3 and 6.
