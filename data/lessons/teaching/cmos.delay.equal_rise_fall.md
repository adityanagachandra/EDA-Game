# Sizing for equal rise and fall

Libraries size each gate so that its worst-case pull-up and pull-down match the drive of a reference inverter. This makes rise and fall delays comparable, gives every cell a well-defined drive strength, and is the sizing that logical effort is defined against. A cell generator or a transistor-level sizer implements exactly these rules.

## The reference inverter

Let a unit nMOS have resistance $R$ and a pMOS of the same width $\gamma R$. A reference inverter with equal rise and fall resistance uses an nMOS of width 1 and a pMOS of width $\gamma$: both then have resistance $R$. Its input capacitance, in units where a width-1 transistor's gate is 1, is $1 + \gamma$, so 3 for $\gamma = 2$. Every other gate is sized so that, along its worst-case conducting path, the resistance is also $R$ for both edges. Then its rise and fall delays into a given load match the inverter's, and it can be called a 1x gate. Larger drive strengths scale every width by the same factor.

## Series stacks get wider

Devices in series add their resistances, so $n$ devices in series must each be $n$ times wider to total one unit resistance. Devices in parallel are sized for the worst case, which is one device conducting alone, so each keeps the unit width. For a NAND$n$ the nMOS are in series and each gets width $n$; the pMOS are in parallel and each keeps width $\gamma$. For a NOR$n$ the nMOS keep width 1 and each pMOS gets $n\gamma$. With $\gamma = 2$: a NAND3 has nMOS width 3 and pMOS width 2, so each input sees $3 + 2 = 5$; a NOR3 has nMOS width 1 and pMOS width 6, so each input sees $1 + 6 = 7$. The NOR's larger input capacitance is the cost of its pMOS stack.

## Compound gates by longest path

For a general series-parallel gate, size each transistor by the length of the longest series path through it: width equals that length times the unit width (1 for nMOS, $\gamma$ for pMOS). In $\overline{AB + C}$ the pull-down has A and B in series (width 2 each) in parallel with C (width 1). The pull-up has A and B in parallel, in series with C, so every pMOS lies on a path of length 2 and gets width $2\gamma$. With $\gamma = 2$, input A sees $2 + 4 = 6$ and input C sees $1 + 4 = 5$. Inputs of one gate can therefore have different input capacitances and different logical efforts, and Liberty records a separate capacitance for each pin. Sizing can be refined further (for example, unequal widths along a path), but this rule is the textbook standard.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), sections 4.3 and 4.4; Sutherland, Sproull and Harris, "Logical Effort: Designing Fast CMOS Circuits" (1999), chapter 4.
