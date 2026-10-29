# Thermal resistance networks

A thermal resistance network models a package as resistors between temperature nodes, with power as current and temperature as voltage. Resistances add in series when heat flows through layers in turn and combine reciprocally when heat has parallel paths. With Kirchhoff's laws, the network gives every node temperature, and it is the compact model that system-level thermal tools exchange.

## Layers in series

When all heat flows through layers one after another, their resistances add: $R = \sum R_i$. A die at 0.05 K/W, an interface material at 0.1 K/W, a lid at 0.02 K/W, a second interface at 0.08 K/W and a heat sink at 0.25 K/W total 0.5 K/W. With 120 W the junction sits 60 K above ambient. The series view identifies the dominant resistance immediately: here the heat sink, then the two interface layers. Improving a layer that contributes 4% of the total cannot help much. The analogy is exact for steady state: temperature differences add like voltages and the same heat current flows through each layer.

## Parallel paths

When heat can leave by two paths, the resistances combine as $R = (1/R_1 + 1/R_2)^{-1}$. A die can lose heat upward through the lid and sink at 0.4 K/W and downward through bumps, substrate and board at 2.0 K/W, giving $1/(2.5 + 0.5) = 0.333$ K/W. The parallel result is always below the smaller resistance. Paths that look negligible still matter: the board path here lowers the junction temperature by 6.7 K at 100 W. In a mobile device with no heat sink, the board and the case may be the main paths, so the network depends on the product, not just the package.

## Heat split between paths

Parallel paths share the heat in inverse proportion to their resistance, exactly like current dividers. With top resistance $R_t$ and bottom $R_b$ from the same junction node to ambient, the fraction leaving through the bottom is $R_t/(R_t + R_b)$. For 0.4 and 2.0 K/W, $0.4/2.4 = 16.7\%$ leaves through the board. This split matters beyond temperature: board heat raises the temperature of neighbouring parts and of the solder joints, which feeds fatigue models. It also shows why a single junction-to-ambient number is not portable between boards, and why standards define two-resistor and multi-node compact models that keep the paths separate.

## Going deeper

JEDEC JESD15 on compact thermal models; Incropera and others, "Fundamentals of Heat and Mass Transfer", chapter 3 on the thermal circuit.
