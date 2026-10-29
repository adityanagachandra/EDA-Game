# System in package and multi-chip modules

A system in package (SiP) places several dies, often from different processes, and passive components in one package. A multi-chip module (MCM) is the older term for several bare dies on a shared substrate. Both trade the integration of one large die for assembly of several smaller ones. The central economic and engineering quantity is the yield of the assembled module, which multiplies the yields of every die and every assembly step.

## Compound module yield

If a module holds $k$ dies with yields $y_1, \ldots, y_k$ and each die attach succeeds with probability $y_a$, and these events are independent, the module yield is $Y = y_a^k \prod_i y_i$. Three dies at 0.95, 0.90 and 0.98 with a 0.99 attach yield give $0.99^3 \times 0.95 \times 0.90 \times 0.98 = 0.8131$, or 81.31%. Yield multiplies, so a module of many parts loses a lot even when each part is good, and one bad die scraps the whole module including the good dies already attached. This is the multiplication the whole field fights, and it is why dies are tested before assembly.

## Known-good die

Testing before assembly replaces each raw die yield with the probability that a die passing test is really good. If wafer test lets a fraction $e$ of bad dies escape among those shipped as good, each tested die is good with probability $1 - e$, and the module yield becomes $Y = \prod (1-e_i) \times y_a^k$. With four dies at 2% escape each and 0.995 attach, $Y = 0.98^4 \times 0.995^4 = 0.9041$. Known-good die (KGD) testing is hard because bare dies are probed at wafer level, often with limited speed and temperature coverage, so test escapes rather than raw defect yield often dominate module loss.

## Cost per good module

The cost of a good module is the cost of everything that went into all modules divided by the yield: $C_{good} = (\sum_i C_i + C_{asm})/Y$. With dies costing 20, 12 and 6 units, assembly 5 units and a module yield of 0.85, the cost per good module is $43/0.85 = 50.59$ units. The formula shows two levers: raise $Y$ by better KGD test, or make failures cheaper by assembling expensive dies last or by allowing rework. It also shows why SiP pays off: small dies from the best process for each function, a radio die and a logic die say, can cost less in total than one large die on the most expensive process, even after assembly loss.

## Going deeper

Tummala, "Fundamentals of Microsystems Packaging", chapter on MCMs and SiP; the IEEE Heterogeneous Integration Roadmap chapter on SiP.
