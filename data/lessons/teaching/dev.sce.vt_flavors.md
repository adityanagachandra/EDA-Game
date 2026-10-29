# Leakage across threshold-voltage flavors

A modern library offers the same cells in several threshold flavors: ultra-low, low, standard and high threshold. Lower threshold is faster and leakier, by a factor set by the subthreshold swing. Leakage-recovery and multi-threshold optimization in implementation tools swap flavors cell by cell, so the arithmetic of threshold steps and leakage ratios is what those tools optimize.

## Leakage ratio between flavors

Off-current scales as $10^{-V_t/S}$. Two flavors that differ by $\Delta V_t$ with everything else equal have a leakage ratio
$$\frac{I_{off,low}}{I_{off,high}} = 10^{\Delta V_t/S}.$$

Example: LVT is 80 mV below SVT and $S = 75$ mV/decade. The ratio is $10^{1.067} = 11.7$: an LVT cell leaks about twelve times as much as its SVT twin.

Drive current changes far less. In saturation it depends on overdrive to a power between 1 and 2, so 80 mV less threshold at a 0.75 V supply might add 15 to 25% current. That asymmetry, an order of magnitude in leakage for a fraction of speed, is why low-threshold cells are rationed to critical paths. Real flavors also differ slightly in DIBL and capacitance, which characterized libraries capture and simple ratios do not.

## Threshold step for a leakage target

Process and library teams often work backwards: a new flavor should cut leakage by a factor $R$. The threshold step needed is
$$\Delta V_t = S \log_{10} R.$$

Example: a factor of 10 at $S = 70$ mV/decade needs 70 mV; a factor of 30 needs $70 \times 1.477 = 103$ mV.

In high-k metal-gate processes, flavors are made mostly with different work-function metal stacks and sometimes with channel or gate-length changes, rather than doping alone. The available steps are therefore discrete and process-specific. For a power-optimization engine, the flavors are a small menu of cells with known delay and leakage tables. A leakage-recovery pass walks paths with positive slack and swaps cells to higher threshold while the slack allows, a greedy or Lagrangian optimization on exactly these trade-offs.

## Leakage of a mixed-flavor block

A block's leakage is the sum over cells. With counts $N_{LVT}$, $N_{SVT}$, $N_{HVT}$ and SVT leakage $I_0$ per cell,
$$I_{block} = I_0\left(N_{LVT}10^{\Delta_1/S} + N_{SVT} + N_{HVT}10^{-\Delta_2/S}\right).$$

Example: 5000 LVT, 40000 SVT and 30000 HVT cells, $I_0 = 1$ nA, LVT 70 mV below SVT, HVT 70 mV above, $S = 70$ mV/decade. LVT cells leak 10 nA each and HVT 0.1 nA. The total is $50 + 40 + 3 = 93$ µA. The 7% of cells that are LVT contribute over half the leakage.

This is the standard pattern: leakage concentrates in the low-threshold minority. An optimization engine therefore gets the most leakage recovery per unit of slack by targeting LVT cells first. The same formula, evaluated at each corner with corner-specific swing and leakage, is what a power report sums from Liberty data.

## Going deeper

Wei, Roy and others, "Design and optimization of dual-threshold circuits for low-voltage low-power applications", IEEE Transactions on VLSI Systems (1999); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.), chapter on timing closure.
