# The cost of a respin

A respin is a second (or third) fabrication of a chip after first silicon fails. It costs a new mask set, more engineering, and months of delay during which the product earns nothing. Because those costs are so large, the expected cost of a respin justifies almost any amount of verification and signoff effort before tape-out. The arithmetic is simple and worth doing explicitly.

## The direct cost

The direct cost of a respin is the masks that must be remade plus the engineering time to debug, fix and re-verify. A full respin remakes all masks; a metal-only fix remakes only the masks of some interconnect layers, which is much cheaper. If a mask set costs 12 million dollars and the fix takes 30 engineers for 4 months at 20,000 dollars per engineer-month, the direct cost is $12 + 30 \cdot 4 \cdot 0.02 = 14.4$ million dollars. Advanced-node mask sets are a large fraction of that sum, which is why respins hurt most on the newest processes.

## The cost of delay

The larger cost is usually time. A respin adds the fix time plus another fabrication cycle, often several months in total. A product that would earn 5 million dollars a month loses that revenue for each month of delay, and in a market with a competitor it may lose share permanently. A 5-month delay at 5 million per month is 25 million dollars of lost revenue, so the total cost of the respin in the example is $14.4 + 25 = 39.4$ million dollars. Delay cost is harder to estimate than mask cost, but it is usually the larger term for any product with a competitive market window.

## Expected cost and the value of verification

Before tape-out the respin is uncertain. If the probability of a respin is $p$ and its cost is $C$, the expected cost is $pC$. Extra verification that lowers $p$ from 0.30 to 0.10 on a 39.4-million-dollar respin is worth $(0.30 - 0.10) \cdot 39.4 = 7.88$ million dollars in expectation. Any verification, signoff margin or tool improvement that costs less than that is a good investment. This expected-value reasoning is the business case for formal verification, emulation, and conservative signoff margins. The same reasoning sizes signoff margins: a margin that costs a little frequency but noticeably lowers $p$ is worth paying for.

## Going deeper

Hennessy and Patterson, "Computer Architecture: A Quantitative Approach", chapter 1 on the cost of integrated circuits and of masks.
