# Energy per operation

Power tells you how fast energy is spent; energy per operation tells you what each computation costs. For battery-powered and throughput-per-watt designs it is the figure that matters, and it reveals a trade-off that power alone hides: slowing down saves dynamic power but lets leakage accumulate per operation.

## Dynamic energy per operation

If each operation switches capacitance $C$ with activity $\alpha$, its dynamic energy is
$$E_{dyn} = \alpha C V_{DD}^2,$$
independent of how long the operation takes. Example: 200 pF with $\alpha = 0.15$ at 0.8 V gives $0.15 \times 200 \times 0.64 = 19.2$ pJ, since pF times volts squared is pJ. Frequency does not appear: running the same operation twice as fast doubles dynamic power but not dynamic energy per operation. Only voltage, capacitance and activity change it, and voltage enters squared, which is why lowering the supply is the main lever for energy efficiency.

## Adding leakage energy

Leakage flows all the time, so each operation also pays leakage power times the cycle time:
$$E_{op} = \alpha C V_{DD}^2 + P_{leak} / f.$$
With the example above, 50 microwatts of leakage at 100 MHz adds $50/100 = 0.5$ pJ, for 19.7 pJ total (microwatts divided by megahertz is picojoules). At 10 MHz the leakage share grows to 5 pJ. The slower the clock, the more leakage energy each operation carries, because the circuit sits idle but powered for longer. This is why low-activity systems power-gate idle blocks instead of merely slowing their clocks.

## The minimum-energy trade-off

Setting the two terms equal gives the frequency below which leakage dominates energy:
$$f^* = \frac{P_{leak}}{\alpha C V_{DD}^2}.$$
For 19.2 pJ and 50 microwatts, $f^* = 2.6$ MHz. Lowering $V_{DD}$ cuts dynamic energy quadratically but slows the circuit, so each operation takes longer and accumulates more leakage energy. Near and below threshold the delay grows exponentially, so leakage energy per operation eventually rises faster than dynamic energy falls: there is a minimum-energy supply voltage, often near or somewhat below the threshold voltage. Energy-optimal design therefore picks voltage, frequency and Vt flavour together, and an energy-aware tool must model both terms rather than optimizing power at a fixed frequency.

## Going deeper

Calhoun, Wang and Chandrakasan, "Modeling and sizing for minimum energy operation in subthreshold circuits" (IEEE JSSC, 2005); Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), chapter 5.
