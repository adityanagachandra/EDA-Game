# VHDL signals, variables and delta cycles

VHDL draws the line between immediate and deferred update with two kinds of object: variables, which change at once, and signals, which change only when the process suspends. Delta cycles order the resulting updates without advancing time. Verilog's blocking and nonblocking assignments are the same idea in another syntax, and every HDL simulator kernel, including SystemC's, is built on this evaluate-then-update loop.

## Signals update when the process suspends

A signal assignment a <= v + 4; inside a process schedules a new value; the signal keeps its old value until the process suspends. Every read of a signal within one activation therefore sees the value from before the activation, however many times the process assigned it. A variable assignment v := c + a; takes effect immediately, so later statements read the new value.

Take a = 3, b = 2, c = 8 and a clocked process that executes v := c + a; a <= v + 4; v := v + 4; c <= v + a; b <= b + c;. Then v = 11, a is scheduled to 15, v becomes 15, c is scheduled to $15 + 3 = 18$ (a still reads 3), and b is scheduled to $2 + 8 = 10$ (c still reads 8). After the edge, a = 15, b = 10 and c = 18.

## Variables keep their values between activations

Variables declared in a process persist from one activation to the next; they are initialized once, at elaboration. This matters when a variable is read before it is written: the read sees the value from the previous activation, and in a clocked process that variable becomes a register, exactly like a Verilog blocking temporary read before it is written.

Over several clock edges, simulate one activation at a time: run the statements with immediate variable updates and deferred signal updates, apply the signal updates, then start the next activation with the new signal values and the persisting variable values. A process that writes its variable before reading it on every path behaves as pure logic; one that does not has hidden state.

## Delta cycles

Concurrent signal assignments outside processes, such as s2 <= x or s1;, are processes sensitive to the signals they read. When input x changes in delta cycle 0, every assignment that reads x wakes and its target updates in delta 1; assignments reading those targets update in delta 2, and so on, all at the same simulation time.

A signal settles one delta after the latest of its inputs, so its final delta number is $\delta(s) = 1 + \max_{u \in \text{inputs}(s)} \delta(u)$ with $\delta(x) = 0$: the longest path from x in the dependency graph. With s1 <= not x; s2 <= x or s1; s4 <= s1 or s2;, s1 settles at delta 1, s2 at delta 2 and s4 at delta 3, whatever order the statements are written in. Delta cycles make the result independent of process order, which is the whole point of deferred update.

## Going deeper

Ashenden, "The Designer's Guide to VHDL" (3rd ed.), chapters 5 and 6; IEEE Std 1076-2019, clause 14.7 (simulation cycle).
