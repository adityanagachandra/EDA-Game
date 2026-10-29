# Wafer probe, final test and burn-in

Patterns become quality only on the test floor. A die is tested on the wafer by a prober, tested again after packaging, and, for high-reliability parts, stressed in burn-in. Each step costs tester time, and tester time is money, so the economics of test, seconds per die, sites per touchdown, units per hour, drive decisions about compression, pattern counts and BIST as much as coverage does.

## Cost per die at probe

At **wafer probe**, a probe card contacts the pads of one or several dies at a time (one **touchdown**), and the automatic test equipment (ATE) applies parametric tests (contact, leakage, supply current), scan patterns and functional tests. Touchdown time is the pattern time plus fixed overhead. With 11,836,145 scan cycles at 100 MHz (118.36 ms) and 258 ms of overhead, a touchdown takes 376.36 ms; at 8 cents per second of tester and prober time, that is 3.01 cents, and with 4 dies tested in parallel, 0.753 cents per die. Testing several dies per touchdown, **multi-site** test, divides the cost, which is why scan compression's reduction in pins per die matters: fewer channels per die means more sites per tester.

## Final test throughput

After packaging, **final test** runs on a handler that indexes parts into sockets. Throughput in units per hour is
$$UPH = \frac{3600 \cdot S \cdot \eta}{t_{test} + t_{index}},$$
for $S$ sites, multi-site efficiency $\eta$ (the fraction of site time doing useful work, below 1 because sites wait for each other and for the handler) and times in seconds. With 2.684 s of test, 0.637 s of index time, 8 sites and $\eta = 0.85$, the cell tests about 7371 units per hour. Final test also screens package defects and performs speed binning.

## Burn-in

Some defects pass every time-zero test and fail early in life: thin oxide, marginal vias. **Burn-in** stresses parts at high temperature and voltage so that these infant-mortality failures happen before shipment. Thermally activated mechanisms follow the Arrhenius law, with acceleration factor
$$AF = \exp\left(\frac{E_a}{k}\left(\frac{1}{T_{use}} - \frac{1}{T_{stress}}\right)\right),$$
temperatures in kelvin, $k = 8.617 \times 10^{-5}$ eV/K. For $E_a = 0.5$ eV, 55 degrees C use and 125 degrees C stress, $AF \approx 22.4$: one burn-in hour ages a part like 22 hours of use. Burn-in is expensive, so many products replace it with voltage stress tests and IDDQ-style screens, keeping full burn-in for automotive and other high-reliability parts.

## Going deeper

Bushnell and Agrawal (2000), chapters 2 and 3 on test economics and ATE; JEDEC JESD94 and related reliability standards for acceleration models.
