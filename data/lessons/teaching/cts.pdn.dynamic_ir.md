# Dynamic IR drop

Static analysis uses average currents. Real current arrives in bursts at clock edges, and for a short time the grid cannot deliver it: local charge from decoupling capacitance does, and the voltage dips. Dynamic IR analysis models the grid with its capacitance and solves it in time. The difference matters because timing failures follow the worst instantaneous droop, which can be several times the static drop.

## An RC model of a grid region

Lump a region as a resistance $R$ from the supply to a node with local decap $C$, and let the load current step from 0 to $I$ at $t = 0$. KCL at the node gives $C\,\dot d = I - d/R$ for the drop $d$, a first-order linear ODE with solution
$$d(t) = IR\left(1 - e^{-t/RC}\right).$$
With $R = 0.3\ \Omega$ and $C = 600$ pF, $RC = 180$ ps (ohms times picofarads give picoseconds). For $I = 330$ mA the static drop is 99 mV, and 380 ps after the step the drop is $99(1 - e^{-2.11}) = 87.0$ mV. Initially the decap supplies the current and the drop rises gradually.

## How fast the droop develops

Inverting the response gives the time to reach a given drop: $t = -RC\ln(1 - d/IR)$. With $R = 0.25\ \Omega$, $C = 1500$ pF and $I = 200$ mA, the final drop is 50 mV and $RC = 375$ ps, so the drop reaches 25 mV after $375 \ln 2 = 259.9$ ps. If the current burst is much shorter than $RC$, the decap rides through it and the drop stays far below $IR$. That observation is the basis of decap sizing: make $RC$ long compared with the burst.

## Why averages mislead

Suppose the current is a 460 mA burst lasting 420 ps every 4620 ps, with $R = 0.5\ \Omega$ and $C = 600$ pF. Static analysis sees an average of 41.8 mA and a drop of 20.9 mV. The first burst, from rest, drives the node to $230(1 - e^{-1.4}) = 173.3$ mV, about 8.3 times as much. Dynamic analysis therefore runs a transient solve of the grid, usually with backward Euler or the trapezoidal rule, using per-instance current waveforms from vectors or from vectorless worst-case switching scenarios. The grid matrix becomes $G + C/h$ at time step $h$, still sparse and positive definite.

## Going deeper

Popovich, Mezhiba and Friedman, "Power Distribution Networks with On-Chip Decoupling Capacitors" (2008); Najm, "Circuit Simulation" (2010), the chapters on transient analysis.
