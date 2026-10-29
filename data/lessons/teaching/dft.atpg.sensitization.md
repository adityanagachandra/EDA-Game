# Fault activation and path sensitization

Every test for a stuck-at fault does two jobs: it activates the fault by driving the line to the opposite of the stuck value, and it sensitizes a path so that the resulting error reaches an output. ATPG algorithms differ in how they search for such a vector, but they all solve these two subproblems. Choosing a test by hand on a small netlist is the clearest way to see the conflicts that make the search hard.

## Activation

To test line $x$ stuck-at-$v$, the fault-free value of $x$ must be $\overline{v}$. This is a **justification** problem: find input values that produce $\overline{v}$ at $x$. For $x = $ NAND(a, b) stuck-at-0, activation needs $x = 1$, which any vector with a = 0 or b = 0 provides. Activation alone is easy to check by simulating the fault-free circuit. Many vectors activate a fault but do not detect it, because the error dies before an output.

## Sensitizing a path

The error at $x$ must travel to an output. Along a chosen path, each gate's **side inputs** must be at the non-controlling value: 1 for AND and NAND, 0 for OR and NOR. Inverting gates flip D to $\overline{D}$ but still propagate it. Setting those side inputs is itself a justification problem, and it can conflict with activation: the side input needed to pass the error through an AND may require an input value that also forces $x$ to the stuck value. When a stem fans out and the paths reconverge, the error can arrive at a gate on two inputs, and D meeting $\overline{D}$ at an AND cancels. Some faults are detectable only by sensitizing several paths at once.

## Branch faults and choosing a test

A fault on a fanout **branch** affects only the gate that branch feeds; the stem and other branches keep their fault-free values. So a test for the branch of c into $g_1$ must propagate through $g_1$, even if c's other branches offer easier paths. To choose a test among candidate vectors, check both conditions for each: compute the fault-free value at the site (activation), then simulate with the fault and compare outputs (propagation). A vector that activates the fault but sets a side input to its controlling value is the most common wrong answer, and it is exactly the case ATPG's D-frontier bookkeeping is designed to avoid.

## Going deeper

Bushnell and Agrawal (2000), chapter 7; Abramovici, Breuer and Friedman (1990), chapter 6.
