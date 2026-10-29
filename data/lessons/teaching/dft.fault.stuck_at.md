# The stuck-at fault model

The single stuck-at model is the workhorse of digital test: each line of a gate-level netlist may be permanently fixed at logic 0 or logic 1, one line at a time. It is simple enough that fault lists, fault simulation and ATPG scale to billions of gates, and good enough that stuck-at test sets catch most gross defects. Every other model in this domain is defined relative to it, so a tool designer must know exactly what a line is and how faults are counted.

## Lines, stems and branches

A **line** is a signal segment that can fail independently. Every primary input and every gate output is a **stem**. When a stem drives two or more gate inputs, each connection is a separate **fanout branch**, because a defect can break one branch without affecting the others. A stem that drives a single input has no separate branch: stem and input are one line. The single stuck-at model puts two faults on each line, stuck-at-0 and stuck-at-1. For a circuit with 4 inputs and 3 gates in which input a feeds two gates, the lines are 7 stems plus 2 branches of a, so 9 lines and 18 single stuck-at faults. Tools usually list faults on gate pins, which is the same count expressed per pin.

## Single versus multiple faults

The **single fault assumption** says that exactly one line is faulty. Without it, each line could independently be fault-free, stuck-at-0 or stuck-at-1, giving $3^L - 1$ multiple stuck-at faults for $L$ lines, which is astronomically large: 20 lines already give about $3.5 \times 10^9$. The single model gives only $2L$. Empirically, a test set that detects all single stuck-at faults also detects a very high fraction of multiple faults, because it is rare for two faults to mask each other on every vector. This is the main reason the single model has dominated for decades.

## Activation and the excitation condition

To detect a line $x$ stuck-at-$v$, the test must first **activate** (excite) the fault: the fault-free circuit must drive $x$ to $\overline{v}$, so that the faulty and fault-free circuits disagree at $x$. Counting the activating vectors is a truth-table exercise: if $x = $ AND(a, b) and the fault is stuck-at-1, three of the four input combinations give $x = 0$ and excite it. Activation alone is not detection; the difference must also propagate to an output. But the activation count already shows why some faults are easy and others hard: a stuck-at-0 on the output of a 10-input AND needs all inputs at 1, a single vector out of 1024.

## Going deeper

Bushnell and Agrawal, "Essentials of Electronic Testing" (2000), chapter 4; Abramovici, Breuer and Friedman, "Digital Systems Testing and Testable Design" (1990), chapter 4.
