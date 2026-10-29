# Generate loops and conditional generate

Generate constructs build hardware structure at elaboration time: a loop that creates one instance per bit, a nested loop that builds a systolic array, a conditional that picks an implementation from a parameter. Elaboration evaluates them completely, so the elaborated design contains only the instances they produce. Counting those instances is the same computation the elaborator performs.

## Generate loops

A generate loop uses a genvar, a variable that exists only during elaboration. for (i = 0; i < 8; i = i + 1) begin : g ... end creates one copy of its body for each value of i, named g[0] to g[7], and inside each copy i is a constant. The loop never becomes hardware; only its copies do.

Counting copies is counting iterations. With start s, bound e and step k, the loop with i < e runs $\lceil (e - s) / k \rceil$ times when $e > s$, and the loop with i <= e runs $\lfloor (e - s) / k \rfloor + 1$ times. So for (i = 0; i <= 10; i = i + 1) creates 11 instances, and for (i = 2; i < 17; i = i + 3) creates 5, at i = 2, 5, 8, 11 and 14. Named blocks matter: the label g gives each copy a stable hierarchical name such as g[3].u.

## Nested loops with conditions

Generate loops nest, and a generate-if inside them keeps only the copies whose condition holds. A two-dimensional processing array with for i < 3 and for j < 4 and the condition (i + j) % 2 == 0 enumerates 12 (i, j) pairs and keeps 6. A triangular condition such as j <= i over $n \times n$ pairs keeps $n(n+1)/2$.

The condition is evaluated per copy during elaboration, so it may use genvars and parameters but nothing that changes at run time. This is how designers build structures such as the upper triangle of a matrix unit or the boundary cells of a mesh, and it is why elaborated designs can have very different instance counts from what the source text suggests.

## Parameter-dependent structures

Bounds can depend on parameters and on outer genvars. A reduction tree of two-input adders over N inputs, with level l creating N >> (l + 1) adders for $l = 0, 1, \dots, L - 1$, builds $N/2 + N/4 + \dots$ adders. With N = 24 and L = 2 the levels contribute 12 and 6, so 18 adders. For N a power of two and $L = \log_2 N$ the total is $N - 1$, the familiar count for a binary tree.

Because the structure depends on parameters, the same module elaborates to different hardware for each instance, and each distinct parameter set is a separate specialization. Recursion through generate-if, where a module instantiates smaller copies of itself until a base case, is the other common pattern, and it terminates only because elaboration evaluates the condition with constants.

## Going deeper

IEEE Std 1800-2017, clause 27 (generate constructs); Sutherland, "SystemVerilog for Design", the chapter on generate.
