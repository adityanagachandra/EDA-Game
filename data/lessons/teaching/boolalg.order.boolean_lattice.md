# The Boolean lattice and implication order

Bit vectors of length $n$, ordered componentwise, form the Boolean lattice $B^n$: the hypercube with 00...0 at the bottom and 11...1 at the top. Boolean functions, ordered by implication, form a second lattice whose elements are subsets of $B^n$. These orders are not decoration. Incompletely specified functions are intervals in the function lattice, monotone functions are order-preserving maps, and fixpoint algorithms in model checking and dataflow analysis are iterations in a subset lattice. Thinking in the order turns many Boolean questions into counting.

## Bit vectors ordered componentwise

Say $u \le v$ when $u_i \le v_i$ for every position, so 0100 $\le$ 1101 but 0110 and 1001 are incomparable. Join is bitwise OR and meet is bitwise AND, so $(B^n, \le)$ is a lattice, in fact a Boolean algebra with complement as bitwise NOT. The rank of a vector is its weight, and rank $k$ holds $\binom{n}{k}$ vectors.

The interval $[u, v]$ with $u \le v$ fixes the positions where $u$ has 1 (they stay 1) and where $v$ has 0 (they stay 0), and frees the rest, so it has $2^{|v| - |u|}$ elements, where $|x|$ is the weight. For $u$ = 0100 and $v$ = 1101 that is $2^{3 - 1} = 4$. An interval of $B^n$ is exactly a cube.

## Functions ordered by implication

For functions, $f \le g$ means $f$ implies $g$: every ON row of $f$ is an ON row of $g$. This is the subset order on ON-sets, so the functions of $n$ variables form the lattice of all subsets of the $2^n$ rows.

The functions $h$ with $f \le h \le g$ are free on the rows in $g$ but not in $f$, so there are $2^{|g| - |f|}$ of them. With $f = m(1, 3)$ and $g = m(1, 3, 5, 7)$ there are $2^2 = 4$. That interval is an incompletely specified function in another guise: $f$ is the ON-set, $g$ is ON plus don't cares, and every $h$ in between is a legal implementation. Minimizers search this interval for the cheapest member.

## Join and meet of functions

In the function lattice the join of $f$ and $g$ is $f + g$ (union of ON-sets) and the meet is $f \cdot g$ (intersection). The bottom is the constant 0 and the top the constant 1. For $f = m(0, 1, 4, 5, 7)$ and $g = m(0, 1, 4, 6)$ over three variables, $f + g = m(0, 1, 4, 5, 6, 7)$ with 6 ON rows and $f \cdot g = m(0, 1, 4)$ with 3.

The lattice view explains why some operations are safe. Replacing a node's function by anything in the interval allowed by its don't cares keeps the circuit correct; over-approximating a reachable set by moving up the lattice keeps a safety proof sound but may add false alarms. Abstract interpretation and approximate reachability are built on exactly this ordering.

## Going deeper

Davey and Priestley, "Introduction to Lattices and Order" (2nd ed., 2002), chapters 1 and 2. Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms", on intervals of Boolean functions.
