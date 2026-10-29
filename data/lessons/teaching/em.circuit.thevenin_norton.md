# Thevenin and Norton equivalents

Any linear network seen from two terminals behaves like one voltage source in series with one resistor (Thevenin) or one current source in parallel with one resistor (Norton). That equivalence is how timers model a driving gate, how extractors summarize a power grid seen from a pin, and how designers reason about loading without solving the whole network again.

## Open-circuit voltage

The Thevenin voltage $V_{th}$ is the voltage at the terminals with nothing connected. With the terminals open, no current flows into them, so any resistor in series with a terminal drops no voltage and can be ignored for this step.

Example: a 3 V source feeds $R_1 = 1$ kohm into node a, $R_2 = 2$ kohm runs from a to ground, and $R_3 = 5$ kohm runs from a to the output terminal. With the output open, $R_3$ carries no current, so the output sits at the divider voltage $3 \cdot 2/3 = 2$ V. In general $V_{th}$ comes from any method that solves the network with the terminals open: a divider, nodal analysis, or superposition over several sources.

## Thevenin resistance

The Thevenin resistance $R_{th}$ is the resistance seen looking into the terminals with every independent source zeroed: voltage sources replaced by shorts and current sources by opens. Dependent sources must stay, and then $R_{th}$ is found by applying a test source.

Continuing the example: shorting the 3 V source puts $R_1$ in parallel with $R_2$, $1 \parallel 2 = 0.667$ kohm, and $R_3$ is in series with that, so $R_{th} = 5.667$ kohm. A timer's driver model is this resistance in spirit: an effective resistance that charges the load. A strong driver has small $R_{th}$; a long series wire adds its resistance in series, which is why buffering a long wire helps.

## Norton current and maximum power transfer

The Norton equivalent replaces the source and series resistor by a current source $I_N = V_{th}/R_{th}$ in parallel with $R_{th}$; $I_N$ is the current that flows when the terminals are shorted. The two forms are interchangeable, and simulators convert between them when they stamp sources.

A resistive load $R_L$ draws power $V_{th}^2 R_L/(R_{th} + R_L)^2$, which peaks at $R_L = R_{th}$ with value $P_{max} = V_{th}^2/(4R_{th})$. Example: with $V_{th} = 2$ V and $R_{th} = 5.667$ kohm, $I_N = 0.353$ mA and $P_{max} = 4/22.67 = 0.176$ mW. Maximum power transfer is not maximum efficiency: at the matched point half the power is lost in $R_{th}$. The same matching idea, with impedances, returns in transmission-line termination.

## Going deeper

Agarwal and Lang, "Foundations of Analog and Digital Electronic Circuits", chapter 3; Alexander and Sadiku, "Fundamentals of Electric Circuits", chapter 4.
