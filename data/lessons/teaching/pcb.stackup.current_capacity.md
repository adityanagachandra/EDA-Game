# Trace current capacity rules of thumb

How wide must a trace be to carry a given current? The trace heats from $I^2 R$ loss and cools through the laminate and the air, so the answer is a temperature-rise limit, not a hard current limit. Board tools still use empirical rules, the oldest being the IPC-2221 chart fit, which a constraint system turns into minimum widths per net class. Knowing the rule, its inversion and its limits lets a tool designer implement width checks honestly.

## The IPC-2221 chart fit

IPC-2221 gives a curve fit to old test data: $I = k\, \Delta T^{0.44} A^{0.725}$, with $I$ in amperes, the temperature rise $\Delta T$ in C, the cross-section $A$ in square mils, $k = 0.048$ for outer layers and $k = 0.024$ for inner layers. Example: a 10 mil outer trace in 1 oz copper has $A = 10 \cdot 1.37 = 13.7$ square mils. For a 10 C rise, $I = 0.048 \cdot 10^{0.44} \cdot 13.7^{0.725} \approx 0.048 \cdot 2.754 \cdot 6.67 \approx 0.88$ A. Current grows slower than linearly with area (exponent 0.725), because a wider trace has proportionally less surface per unit cross-section through which to shed heat. Doubling the allowed rise raises the current only by $2^{0.44} \approx 1.36$.

## Inverting the rule for width

Constraint checks need the inverse: the minimum width for a current. Solve for area, $A = (I / (k\, \Delta T^{0.44}))^{1/0.725}$, then divide by thickness. Example: 2 A with a 10 C rise on an outer 1 oz layer gives $A = (2 / (0.048 \cdot 2.754))^{1.379} \approx 15.13^{1.379} \approx 42.4$ square mils, so $w = 42.4 / 1.37 \approx 31$ mil. On an inner layer $k$ halves, and the area must grow by $2^{1/0.725} \approx 2.6$, to about 80 mil. Heavier copper helps directly: 2 oz halves the width for the same area. A width rule per net class, computed this way, is what the DRC enforces on power nets.

## Limits of the rule

The IPC-2221 curves come from limited tests decades ago, and the inner-layer constant is a conservative halving rather than a measurement. **IPC-2152** replaced the chart with newer data showing that board thickness, nearby copper planes and the laminate's thermal conductivity matter as much as the layer position; planes spread heat well, so internal traces next to planes often run cooler than the old rule predicts. Neither rule handles short neck-downs, many parallel hot traces, or vias, which have their own current limits. A tool that reports one number should say which rule produced it. For high-current designs the honest method is a thermal simulation, using the rule of thumb only to size the first draft.

## Going deeper

IPC-2152, "Standard for Determining Current Carrying Capacity in Printed Board Design". IPC-2221, section on conductor sizing.
