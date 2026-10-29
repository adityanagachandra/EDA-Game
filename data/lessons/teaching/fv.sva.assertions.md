# SystemVerilog assertions

SystemVerilog assertions (SVA) are how hardware teams write temporal properties that both simulators and formal tools check. A concurrent assertion is sampled on a clock edge and, for an implication, starts a new attempt every cycle its antecedent holds. Counting failing and pending attempts on a waveform by hand is the quickest way to learn exactly what an assertion means, which is the first thing to get right before any proof.

## Overlapping and non-overlapping implication

assert property (@(posedge clk) a |-> b) samples the signals at each rising edge. At every cycle where the antecedent $a$ is 1, an attempt starts and checks the consequent: with the overlapping implication |->, $b$ must be 1 in the same cycle; with the non-overlapping implication |=>, $b$ must be 1 in the next cycle. Cycles where $a$ is 0 start no attempt and cannot fail. Take cycles 0 to 6 with $a$ = 0110100 and $b$ = 0011001, leftmost bit first. Attempts start at cycles 1, 2 and 4. For a |-> b, cycle 1 fails ($b = 0$), cycle 2 passes and cycle 4 fails: 2 failures. For a |=> b, cycle 1 checks $b$ at 2 (passes), cycle 2 checks cycle 3 (passes) and cycle 4 checks cycle 5 (fails): 1 failure. An attempt whose check cycle lies beyond the end of a finite trace is pending, neither passed nor failed.

## Fixed delays

The cycle delay ##n moves the check $n$ cycles forward: a |-> ##2 b requires $b$ exactly two cycles after each cycle where $a$ holds, and a |=> b is the same as a |-> ##1 b. With the trace above, $a$ = 0110100 and $b$ = 0011001, the attempts at cycles 1, 2 and 4 check $b$ at cycles 3, 4 and 6: 1, 0 and 1. So one attempt fails, the one started at cycle 2. If $a$ were also 1 at cycle 5, that attempt would need $b$ at cycle 7, beyond a 7-cycle trace, and would be pending. In simulation, pending attempts at the end of a test are usually reported separately; in formal verification there is no end of trace, and BMC only reports an assertion as failing when a complete failing attempt fits within its bound.

## Delay ranges

A range ##[m:n] lets the consequent happen at any cycle from $m$ to $n$ cycles after the antecedent: a |-> ##[1:3] b passes as soon as $b$ is seen in that window and fails only when the whole window passes without $b$. Take cycles 0 to 7 with $a$ = 10111000 and $b$ = 00011000 and the property a |-> ##[1:2] b. The attempt at cycle 0 checks cycles 1 and 2, both 0: fail. The attempt at 2 sees $b$ at 3: pass. The attempt at 3 sees $b$ at 4: pass. The attempt at 4 checks cycles 5 and 6, both 0: fail. Two failures. An attempt whose window is cut off by the end of the trace passes if $b$ already appeared and is pending otherwise. Unbounded ranges, whose upper end is written with a dollar sign, describe liveness, and their formal treatment differs.

## Going deeper

IEEE Std 1800 (SystemVerilog), the chapter on assertions; Cerny, Dudani, Havlicek and Korchemny, "SVA: The Power of Assertions in SystemVerilog" (2nd ed., 2015).
