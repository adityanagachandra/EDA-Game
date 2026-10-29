# Transistor counts and series stacks

Area, input capacitance and worst-case delay of a static gate all follow from two numbers that can be read straight off its Boolean expression: how many transistors it has, and how tall its tallest series stack is. Cell-library planners use them to decide which complex gates are worth offering, and synthesis cost models use them to estimate area. They also explain why libraries stop at three or four series devices.

## Two transistors per literal

In a single-stage static CMOS gate each literal of the expression drives one nMOS in the pull-down network and one pMOS in the pull-up network. The transistor count is therefore twice the number of literal occurrences. A NAND$n$ or NOR$n$ has $2n$ transistors: a NAND4 has 8. A compound gate such as $\overline{AB + CD + E}$ has 5 literals and so 10 transistors in one stage, where a NAND-and-inverter implementation needs several stages. If a variable appears twice in the expression, it drives two transistors in each network. The count ignores any output inverter: a non-inverting cell such as AND2 is a NAND2 plus an inverter, 6 transistors.

## The tallest series stack

Delay and sizing depend on the longest chain of series devices from the output to a rail. Compute it recursively from the expression. In the pull-down network AND is series and OR is parallel: a series group's height is the sum of its children's heights, a parallel group's height is the largest child's height, and a literal has height 1. For $\overline{(A + B)(C + D)}$ the PDN is two parallel pairs in series, height $1 + 1 = 2$. The pull-up network is the dual, so the roles flip: OR terms are series. For $\overline{AB + C + D}$ the PDN has height 2 (the AB pair), while the PUN puts the three OR terms in series, with AB as a parallel pair, height $1 + 1 + 1 = 3$.

## Why stacks are kept short

A series stack of $n$ devices has $n$ times the resistance of one, so each device must be made $n$ times wider to keep the gate's drive, which multiplies its input capacitance and its internal diffusion. Delay through a stack also grows faster than linearly once internal node capacitances are counted, and the device nearest the output suffers body effect and a reduced gate overdrive. pMOS stacks are worst because holes are slower, so a NOR4 needs pMOS devices about $4\gamma$ wide where $\gamma$ is the pMOS-to-nMOS width ratio of a unit inverter. This is why libraries offer NAND4 but rarely NOR4, and why compound gates are usually limited to stacks of three or four.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), sections 1.4 and 9.2.
