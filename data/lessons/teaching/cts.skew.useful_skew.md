# Useful skew

Zero skew is not the goal; meeting timing is. If a flop's incoming path is too long and its outgoing path has slack, delaying that flop's clock lends time from one path to the other. This deliberate skew is called useful skew. Each flop's clock arrival then lives in a feasible window bounded by the setup and hold constraints of its neighbouring paths. Computing that window is the local version of clock scheduling, and modern CTS engines implement the result as latency targets.

## Lending time to a failing path

Let flop F capture from A and launch to B, with the clock reaching A and B at time 0 and F at time $s$. The path A to F meets setup when $t_{cq} + d_{max} + t_{su} \le s + T$. With period 800 ps, $t_{cq} = 45$, $d_{max} = 736$ and $t_{su} = 33$ ps, setup fails at $s = 0$ by 14 ps, and it is met for $s \ge 45 + 736 + 33 - 800 = 14$ ps. A later capture clock moves the required time out by exactly $s$. Nothing about the logic changes: the time comes from the next stage.

## The upper bounds

The same $s$ delays every launch from F. Setup on F to B requires $s + t_{cq} + d_{max,FB} + t_{su} \le T$, so $s \le T - t_{cq} - d_{max,FB} - t_{su}$. Hold on A to F requires the earliest new data to arrive after F's clock plus hold: $t_{cq} + d_{min,AF} \ge s + t_h$, so $s \le t_{cq} + d_{min,AF} - t_h$. With $T = 550$, $t_{cq} = 51$, $t_{su} = 30$, $t_h = 26$, $d_{max,FB} = 354$ and $d_{min,AF} = 187$, the bounds are 115 and 212 ps, and the tighter, 115 ps, wins. Hold is the constraint people forget: pushing a capture clock later also makes short paths into it fail.

## The feasible window

A fourth constraint, hold on F to B, gives another lower bound: $s + t_{cq} + d_{min,FB} \ge t_h$. The window is then $[\max(\text{lower bounds}), \min(\text{upper bounds})]$. With lower bounds 16 and -46 and upper bounds 175 and 180 ps, the window is $[16, 175]$, 159 ps wide. An empty window means skew cannot fix this flop alone and the logic must change. A wide window gives CTS freedom: it can implement any value inside it, and robust schedules aim for the middle so that variation cannot push the arrival outside. Choosing all flops' arrivals together is the scheduling problem.

## Going deeper

Fishburn, "Clock skew optimization", IEEE Transactions on Computers (1990); Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the timing closure chapter.
