# Thermal-aware 3D floorplanning

In a 3D stack, dies far from the heat sink must push their heat through the dies nearer to it, and thinned dies spread heat poorly. Thermal-aware floorplanning decides which die goes where in the stack and where hot blocks sit within each die, so that peak temperature stays in budget. Simple resistance models make the trade-offs computable inside a floorplanner's inner loop.

## Peak temperature of a stack

Model a two-die stack with the heat sink above: ambient connects through $R_s$ to die A, next to the sink, and die A connects through $R_b$ to die B, farther away. All heat passes through $R_s$, and die B's heat also passes through $R_b$. So $T_A = T_a + (P_A + P_B) R_s$ and $T_B = T_A + P_B R_b$. With $P_A = 20$ W, $P_B = 60$ W, $R_s = 0.4$ K/W, $R_b = 0.3$ K/W and 40 °C ambient, $T_A = 72$ °C and $T_B = 90$ °C. The far die is the hottest whenever it dissipates power, and every watt in it heats the near die too.

## Ordering dies in the stack

Swapping the dies changes only the second term: the peak becomes $T_a + (P_A + P_B)R_s + P_{far} R_b$. Putting the 60 W die next to the sink and the 20 W die far gives a peak of $72 + 20 \times 0.3 = 78$ °C instead of 90, a 12 K reduction for the same total power. The rule is to put the high-power die nearest the sink, which often means logic near the sink and memory farther, but the electrical design may prefer the opposite order, for example to give the logic die direct access to package bumps. The floorplanner must weigh both, and thermal TSVs or dummy bumps can lower $R_b$ when the order is fixed.

## Lateral placement of hot blocks

Within a die, blocks heat each other. A compact model writes each block's temperature as $T_i = T_a + R_{ii} P_i + \sum_j R_{ij} P_j$, where $R_{ij}$ is the mutual thermal resistance, falling with distance. For two 10 W blocks with $R_{ii} = 2$ K/W, the mutual term is 0.8 K/W when adjacent and 0.2 K/W when separated, so the peak falls from $T_a + 28$ to $T_a + 22$, a 6 K gain from placement alone. Thermal-aware floorplanners add peak temperature from such models to wirelength and area in their cost function. In 3D, the mutual resistances between vertically aligned blocks on different dies are large, so hot blocks should not be stacked directly on each other.

## Going deeper

Cong, Wei and Zhang, "A thermal-driven floorplanning algorithm for 3D ICs", ICCAD (2004); Huang and others, "HotSpot: a compact thermal modeling methodology", IEEE Transactions on VLSI Systems (2006).
