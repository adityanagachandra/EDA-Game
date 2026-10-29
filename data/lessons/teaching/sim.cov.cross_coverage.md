# Cross coverage

Many bugs appear only in combinations: a write to the last address while a flush is pending, a short packet on the highest-priority channel. Cross coverage records combinations of coverpoint bins. It is the most informative part of a coverage model and the most expensive, because the number of combinations multiplies. Counting cross bins exactly, including the ones excluded with ignore and illegal bins, is how a verification engineer sizes the closure effort.

## Crosses as Cartesian products

A `cross` of coverpoints creates one bin for every combination of their bins. If `cp_op` has 4 bins (read, write, atomic, flush) and `cp_size` has 3 bins, `cross cp_op, cp_size` has $4 \times 3 = 12$ bins. Adding `cp_mode` with 5 bins makes $4 \times 3 \times 5 = 60$. A cross is covered bin by bin, exactly like a coverpoint: a bin is hit when one sample falls in all of its component bins at the same sampling event. Because the components must coincide in one sample, a cross measures something neither coverpoint can: the two fields can each be fully covered while most of their combinations never occurred.

## Removing ignored and illegal combinations

Not every combination is meaningful. `ignore_bins` removes combinations from coverage, and `illegal_bins` removes them and also reports an error if they are ever hit. Selections use `binsof`: `ignore_bins i1 = binsof(cp_a.a0);` removes every cross bin whose cp_a component is a0. Count removals as a set union. Example: cp_a has 4 bins and cp_b has 5, so 20 combinations. Ignoring cp_a's a0 removes 5, ignoring cp_b's b2 removes 4, and the two overlap in (a0, b2), so 8 are removed. An illegal combination (a3, b4) outside both removes one more: 11 bins remain. An illegal combination that is already ignored removes nothing new.

## Coverage of a cross and its cost

Cross coverage is the number of remaining bins hit divided by the number of remaining bins. Samples that land in ignored bins and repeated hits add nothing. With 11 remaining bins, samples hitting 7 distinct remaining combinations give $7/11 \approx 63.6\%$. The cost grows multiplicatively: crossing three 16-bin coverpoints gives 4,096 bins, and with uniform random stimulus the number of tests to hit them all grows like $n \ln n$ in the number of bins, as the coupon-collector argument shows. Good models cross only fields whose interaction the specification makes meaningful, collapse values into a few behavioral bins before crossing, and prune impossible combinations explicitly so the denominator is honest.

## Going deeper

IEEE Std 1800-2017, clause 19.6 (cross coverage); Piziali, "Functional Verification Coverage Measurement and Analysis" (2004).
