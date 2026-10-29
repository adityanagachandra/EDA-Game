# Expected tests to close coverage bins

How many random tests does it take to cover every bin? If tests hit bins at random, coverage closure is the coupon collector problem: the first bins come quickly and the last ones take most of the time. The analysis explains the familiar shape of a coverage curve, sets realistic regression budgets, and shows quantitatively why closing the last holes needs directed or biased stimulus rather than more seeds.

## Waiting for one bin

If each random test hits a particular bin with probability $p$, independently, the number of tests until the first hit is geometric with mean $1/p$. A bin hit by 4 out of every 1,000 tests on average needs $1000/4 = 250$ tests in expectation. The distribution has a long tail: the probability of still missing after $k$ tests is $(1 - p)^k$, so after 250 tests there is still a 37% chance of no hit. Rare bins dominate closure time. Estimating $p$ for the hardest bins, from early regression data or by sampling the constraint solver, is the most useful single number when planning a coverage effort.

## All n bins: n times the harmonic number

Suppose each test hits exactly one of $n$ equally likely bins. With $j$ bins still missing, a test hits a new bin with probability $j/n$, so the wait for the next new bin has mean $n/j$. Summing over $j = n, n-1, \ldots, 1$ gives the expected total $n \left(1 + \tfrac{1}{2} + \cdots + \tfrac{1}{n}\right) = n H_n \approx n (\ln n + 0.5772)$. For $n = 10$: $H_{10} \approx 2.9290$, so about 29.29 tests, nearly three times the number of bins. For $n = 100$: about 518.7 tests. Unequal bin probabilities make it worse, because the rarest bin governs the tail.

## The last bins dominate

The same argument gives the remaining effort from a partly covered state. With $k$ of $n$ bins covered, the expected number of further tests is $n H_{n-k}$. For $n = 40$ and $k = 36$: $40 \cdot H_4 = 40 \cdot 2.0833 \approx 83.33$ tests for the last 4 bins, while the first 36 took about $40(H_{40} - H_4) \approx 87.8$. Ten percent of the bins cost nearly half the effort. This is why coverage curves flatten into plateaus and why "run more seeds" stops working: the expected cost of each remaining bin is $n$ tests, and with non-uniform probabilities it can be far more. Biasing stimulus toward the holes raises their probabilities directly and shortens the tail.

## Going deeper

Feller, "An Introduction to Probability Theory and Its Applications", volume 1, the section on the coupon collector; Mitzenmacher and Upfal, "Probability and Computing" (2nd ed., 2017), chapter 2.
