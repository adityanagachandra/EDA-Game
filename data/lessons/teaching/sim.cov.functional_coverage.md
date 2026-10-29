# Functional coverage

Functional coverage measures whether the scenarios in the verification plan actually happened. Unlike code coverage it is written by engineers: a coverage model of the features, values and combinations that matter. In SystemVerilog the model is a covergroup made of coverpoints and bins, sampled at chosen events. The model, not the simulator, determines what the coverage number means, so its precise binning and weighting rules are worth knowing exactly.

## Covergroups, coverpoints and bins

A `covergroup` is sampled on an event, such as a clock edge or an explicit `sample()` call from a monitor. Each `coverpoint` watches an expression, and its **bins** partition the values of interest. `bins low = {[0:3]};` is one bin covering four values. `bins mid[] = {[4:7]};` makes one bin per value: 4 bins. `bins high[2] = {[8:15]};` splits the eight values into 2 bins. Together these declare $1 + 4 + 2 = 7$ bins. With no bins declared, the tool creates automatic bins: one per value for small ranges, capped by `option.auto_bin_max`, default 64, so an unsigned 8-bit coverpoint gets 64 bins of 4 values each. Choosing bins is a modeling decision: each bin should be a distinct situation the design must handle.

## Coverpoint coverage

A bin is covered when its hit count reaches `option.at_least`, default 1. A coverpoint's coverage is the number of covered bins divided by the number of bins. Example: 8 bins with hit counts 3, 0, 1, 5, 2, 0, 4, 1 and `at_least = 2`: the bins with 3, 5, 2 and 4 qualify, so coverage is $4/8 = 50.0\%$. With the default threshold of 1, six bins qualify and coverage is 75%. Raising `at_least` asks for repeated evidence, which is useful for scenarios that should be exercised under several surrounding conditions. Sampling is part of the model: a covergroup sampled on every clock counts idle cycles, while one sampled from a monitor counts transactions.

## Covergroup coverage and weights

A covergroup's coverage is the weighted average of its items (coverpoints and crosses), using `option.weight`, default 1. Example: cp1 at 80% with weight 2, cp2 at 50% with weight 1, cp3 at 100% with weight 1. The covergroup is at $(2 \cdot 80 + 50 + 100) / 4 = 77.5\%$. A weight of 0 removes an item from the average while still collecting it, which is how teams track an informational coverpoint without letting it move the sign-off number. Averaging hides detail: a group at 95% can still have one critical bin at zero. That is why closure reviews read bin-level holes, not just the percentage at the top of the report.

## Going deeper

IEEE Std 1800-2017, clause 19 (functional coverage); Piziali, "Functional Verification Coverage Measurement and Analysis" (2004).
