# Functional versus syntactic correctness

Generated RTL can be syntactically correct, compiling cleanly, and functionally wrong. Benchmarks judge function with testbenches, which are themselves imperfect: a test that does not exercise a bug lets it pass. The probability that a bug escapes random testing, and the inflation that escapes cause in a reported pass rate, are both computable, and both argue for equivalence checking where a reference exists.

## Escape probability

Suppose a bug changes the output on a fraction $\delta$ of input vectors, and a testbench applies $T$ independent uniformly random vectors. The bug escapes when every vector misses it, with probability $(1 - \delta)^T$. With $\delta = 0.01$ and $T = 100$, $(0.99)^{100} = 0.366$: more than a third of such buggy modules pass. Directed tests do better on the cases their author anticipated and no better on the rest. A combinational equivalence check between the generated module and the reference, using a miter and a SAT solver, finds the bug with certainty when it exists.

## Vectors for a target escape rate

To push the escape probability below $\epsilon$, solve $(1 - \delta)^T \le \epsilon$: $T \ge \ln\epsilon/\ln(1 - \delta)$. With $\delta = 0.01$ and $\epsilon = 0.01$, $T \ge 458.2$, so 459 vectors. The requirement grows as $1/\delta$: a bug that affects one input in a million needs millions of random vectors. Real bugs are rarely uniform over inputs; corner-case bugs (overflow, a specific state, a reset sequence) have tiny $\delta$ under random stimulus, which is why constrained-random testing, coverage and formal methods exist. Random testing never reaches certainty: it lowers the escape probability, and only a proof removes it.

## Inflated pass rates

If a fraction $c$ of a model's samples are truly correct and an incorrect sample passes a weak testbench with probability $e$, the benchmark reports a pass rate of $c + (1 - c)e$. With $c = 0.5$ and $e = 0.2$ it reports 60 percent, overstating true correctness by 10 points. The bias favors models whose errors are subtle over models whose errors are blatant, which can reorder a leaderboard. Stronger testbenches or equivalence checks against hidden reference solutions reduce $e$; syntax-only metrics such as compile rate set $e$ near one for every compiling sample and measure almost nothing about function.

## Going deeper

Liu et al., "VerilogEval: Evaluating Large Language Models for Verilog Code Generation" (ICCAD 2023); Kuehlmann and Krohm, "Equivalence Checking Using Cuts and Heaps" (DAC 1997), on combinational equivalence checking.
