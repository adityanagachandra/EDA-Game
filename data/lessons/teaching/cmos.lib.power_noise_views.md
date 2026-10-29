# Liberty power and noise views

Liberty describes power and noise as well as timing. Leakage is given per cell and per input state, internal energy per transition is given in tables like delay, and noise views describe how a cell rejects and propagates glitches. Power-analysis and signal-integrity tools read these groups, so their exact meaning decides whether a power report is right.

## State-dependent leakage

A cell's leakage depends on its input state, because which transistors are off, and with what voltages across them, changes with the inputs. Liberty records this with leakage_power groups carrying a when condition, such as when : "A & !B", and a value, plus an optional default cell_leakage_power. A power tool computes the expected leakage as the probability-weighted sum over states. For a NAND2 with independent inputs, $P(A = 1) = 0.5$ and $P(B = 1) = 0.2$, and state leakages of 4, 10, 12 and 30 nW for AB = 00, 01, 10, 11, the state probabilities are 0.4, 0.1, 0.4 and 0.1, giving $1.6 + 1.0 + 4.8 + 3.0 = 10.4$ nW. Using state probabilities from simulation rather than independence gives a better estimate.

## Internal power tables

For each output pin and related input, an internal_power group stores tables of energy per transition, rise_power and fall_power, indexed like delay tables by input transition and output load. The energy covers what is dissipated inside the cell: short-circuit current and charging of internal nodes. It deliberately excludes charging the external load, which the power tool computes from the net's capacitance. A tool multiplies the looked-up energy by the transition rate: an output rising 200 million times per second at a grid point of 6 fJ per rise contributes $6 \times 200 \times 10^{6}$ fJ per second, which is 1.2 microwatts. Input pins can carry internal power too, for energy drawn when an input switches without changing the output.

## Total pin power and noise views

Total dynamic power of a toggling output is internal plus switching. If an output toggles $r$ times per second (both edges), each toggle costs its internal energy plus $\tfrac{1}{2} C V_{DD}^2$ for the net, since a full cycle of two toggles costs $C V_{DD}^2$. With 5 fJ internal, 12 fF at 0.8 V and 100 million toggles per second: switching energy is 3.84 fJ, total 8.84 fJ per toggle, or 0.884 microwatts. Noise views serve signal integrity: older noise immunity tables give the glitch height and width a pin tolerates, and CCS noise models describe a cell's response to and propagation of glitches, plus its holding resistance when quiet, which crosstalk analysis uses for the victim driver.

## Going deeper

The Liberty Reference Manual, power and noise sections; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), chapter 5.
