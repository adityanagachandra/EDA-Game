# Setup, hold and clock-to-q

A flip-flop samples correctly only if its input is stable for a short window around the clock edge, and its output changes only after a delay. These three numbers, setup time, hold time and clock-to-q delay, are characterized per cell in the Liberty library and are the parameters of every check static timing analysis performs. They connect the combinational delays of the logic between registers to the clock period.

## Three cell timing parameters

- **Clock-to-q** $t_{cq}$: the delay from the active clock edge to the new value appearing at $Q$.
- **Setup time** $t_{su}$: how long before the edge $D$ must be stable.
- **Hold time** $t_h$: how long after the edge $D$ must stay stable.

Inside a master-slave flip-flop, setup is the time the master needs to capture $D$ before it closes, and hold covers the delay of the clock reaching the master's gate after the edge. Liberty libraries store these as tables indexed by data and clock transition times, so a timer looks them up per instance. Violating either window can leave the flip-flop metastable, which is a separate topic.

## The setup and hold checks

Between a launching and a capturing flip-flop on the same clock with period $T$ and combinational delay between them in $[t_{\min}, t_{\max}]$:

$$\text{setup slack} = T - (t_{cq} + t_{\max} + t_{su}),$$
$$\text{hold slack} = (t_{cq,\min} + t_{\min}) - t_h.$$

With $T = 1000$ ps, $t_{cq} = 60$, $t_{\max} = 820$, $t_{su} = 40$, setup slack is $1000 - 920 = 80$ ps. With $t_{cq,\min} = 45$, $t_{\min} = 15$, $t_h = 70$, hold slack is $60 - 70 = -10$ ps, a violation. Note the period appears only in the setup check: hold compares the launch and capture of the same edge.

## Clock skew shifts both checks

If the capture flip-flop's clock arrives $s$ later than the launch flip-flop's (positive skew), the capture edge moves out by $s$:

$$\text{setup slack} = T + s - (t_{cq} + t_{\max} + t_{su}),\qquad \text{hold slack} = t_{cq,\min} + t_{\min} - (t_h + s).$$

Positive skew helps setup and hurts hold by the same amount. With $s = 30$ ps the example above has setup slack 110 ps and hold slack -40 ps. Useful-skew optimization deliberately schedules clock arrival times to borrow slack from one stage for another, while clock tree synthesis normally tries to make skew small. A hold violation cannot be fixed by slowing the clock, which is why it is fatal on silicon and why hold fixing (inserting delay) is part of every flow.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 3.5; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 7.
