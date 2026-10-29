# Tool-feedback repair loops

An LLM harness rarely accepts the first output. It runs a tool, returns the errors to the model and tries again, up to a budget. Simple probability shows what such loops buy and what they cost, and it shows why compile success overstates progress when the real target is functional correctness. The independence assumptions behind the numbers are optimistic, which is part of the lesson.

## Independent retries

If each attempt compiles with probability $p$ and attempts are independent, the probability that at least one of $r$ attempts compiles is $1 - (1 - p)^r$. With $p = 0.4$ and $r = 5$ that is $1 - 0.6^5 = 0.922$. Real attempts are correlated: a model that misunderstands the task tends to fail the same way repeatedly, and error feedback helps on some failures and not on others. The independent formula is therefore an upper bound on what retries alone deliver, and measured success curves, success against attempt number, are the honest description of a loop.

## Expected tool calls

A loop that stops at the first success, or after $r$ attempts, makes on average $E[N] = \sum_{i=0}^{r-1}(1 - p)^i = (1 - (1 - p)^r)/p$ attempts, each one model call plus one tool run. With $p = 0.4$ and $r = 5$ that is $(1 - 0.0778)/0.4 = 2.31$ attempts; without a cap it would be $1/p = 2.5$. The tool run dominates the cost when it is a simulation or a synthesis, so the loop's budget is set by tool time, and a cheap first gate (a fast lint before a long simulation) saves most of it.

## Compiling is not correct

If a compiling module passes the functional testbench with probability $q$, a round succeeds with probability $pq$, and $r$ rounds give $1 - (1 - pq)^r$. With $p = 0.4$, $q = 0.5$ and $r = 5$, $pq = 0.2$ and the loop succeeds with probability $1 - 0.8^5 = 0.672$, well below the 0.922 compile rate. A harness that reports compile success as progress overstates what it delivers. And the testbench is itself a filter with false passes, so even the functional rate is an upper estimate of correctness unless the final check is an equivalence proof against a reference.

## Going deeper

Chen et al., "Evaluating Large Language Models Trained on Code" (2021), on sampling and selection; Ross, "A First Course in Probability," chapter 4, on the geometric distribution.
