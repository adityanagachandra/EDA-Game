# FinFETs

A FinFET wraps the gate around a thin vertical fin of silicon, controlling the channel from both sidewalls and the top. The thin fin gives strong electrostatic control, lowering subthreshold swing and DIBL, and the vertical sidewalls give more channel width per unit of layout area. For EDA the most visible consequence is that transistor width became an integer number of fins, which reshaped standard cells, sizing and the compact models behind them.

## Effective width of a FinFET

Current flows along the fin's surfaces under the gate: two sidewalls of height $H_{fin}$ and the top of width $W_{fin}$. The effective electrical width of one fin is approximately
$$W_{fin,eff} = 2H_{fin} + W_{fin},$$
and a device with $N$ fins has $W_{eff} = N(2H_{fin} + W_{fin})$.

Example: fins 45 nm tall and 7 nm wide give 97 nm per fin. Three fins give $W_{eff} = 291$ nm.

Because the fins are tall and closely spaced, the effective width per unit of footprint can exceed what a planar device of the same footprint offers. The fin width is kept small, several nanometres, so that the body is fully depleted and the gate controls it from both sides. The fin is typically undoped or lightly doped, with threshold set mostly by the gate work-function metal, which reduces random dopant fluctuation.

## Drive current from effective width

Processes quote drive current per unit of effective width or per unit of footprint; the two differ, so the convention matters. With current per effective width $I'$,
$$I_{on} = W_{eff}\cdot I'.$$

Example: 3 fins of 97 nm effective width each (291 nm $= 0.291$ µm) at $I' = 900$ µA/µm give $I_{on} = 262$ µA.

Electrostatically, the double-gate-like control gives subthreshold swings near 65 to 70 mV/decade and DIBL of a few tens of mV/V, much better than late planar bulk. This let FinFET processes lower threshold and supply voltage while keeping leakage acceptable. On the other hand, fins add parasitics: gate-to-contact capacitance between closely spaced fins and gates, and series resistance in the narrow source and drain regions. In many FinFET cells these parasitics, not the intrinsic channel, dominate delay, which is why middle-of-line extraction became critical.

## Width quantization

A FinFET's width comes only in whole fins. A device that needs current $I$ when each fin delivers $I_{fin}$ uses
$$N = \lceil I / I_{fin} \rceil$$
fins, and overshoots the requirement.

Example: 25 µA per fin and a requirement of 140 µA gives $\lceil 5.6 \rceil = 6$ fins, delivering 150 µA.

Quantization reshaped design. Standard-cell drive strengths map to fin counts, and cell height is set in fin pitches (track heights). Analog designers lost fine width control and use fin count and finger count instead, with matching depending on identical fin arrays. Sizing algorithms in synthesis and optimization operate on discrete cell sizes anyway, but custom circuits and SRAM bitcells, where ratios between transistors matter, must work with integer ratios such as 1:1:2 fins. Compact models such as BSIM-CMG take fin count, fin height and fin width as instance parameters.

## Going deeper

Hisamoto and others, "FinFET: a self-aligned double-gate MOSFET scalable to 20 nm", IEEE Transactions on Electron Devices (2000); Colinge, "FinFETs and Other Multi-Gate Transistors" (2008); the BSIM-CMG technical manual.
