# Multiple patterning

When a layer's pitch is below what one exposure can resolve, the pattern is split across several masks, each printing a sparser subset, and the results are combined on the wafer. Litho-etch-litho-etch (LELE) and its triple and quadruple variants split by assigning shapes to masks; self-aligned methods split by spacer deposition. For a tool designer, multiple patterning means coloured layouts, colour-aware rules and mask counts that DTCO has to price.

## Pitch splitting

In LELE, the first exposure and etch print every other line into a hard mask, and the second prints the lines in between. Each mask sees twice the final pitch. With $n$ litho-etch steps and lines assigned in rotation, each mask prints a pitch of $n p$. A 40 nm pitch split three ways asks each mask for 120 nm. The final pattern depends on overlay between the exposures: an overlay error shifts one set of lines relative to the other and changes alternate spaces, which is why same-mask and different-mask spacings are separate rules and why critical dimensions on multiple-patterned layers carry extra variation.

## How many exposures a pitch needs

A single exposure resolves a minimum pitch $P_{min} = 2 k_1 \lambda/NA$. To print pitch $p$ by splitting into $n$ masks, each mask's pitch $n p$ must be at least $P_{min}$, so $n = \lceil P_{min}/p \rceil$. For 193 nm immersion at NA 1.35 and $k_1 = 0.28$, $P_{min} = 80.1$ nm. A 64 nm pitch gives $80.1/64 = 1.25$ and needs 2 exposures; a 40 nm pitch gives $80.1/40 = 2.003$ and needs 3. Near such boundaries a small improvement in $k_1$ saves a whole mask, which is the sort of trade a DTCO study makes.

## Counting masks

Mask count drives cost and cycle time. Single exposure costs one mask per layer; LELE two; LELELE three. Self-aligned double or quadruple patterning needs one mandrel mask but usually one or more cut or block masks to remove unwanted line segments. A module with M1 by SADP plus two cut masks (3), M2 by LELE (2), V1 by LELELE (3) and M3 by single exposure (1) needs 9 masks for four layers. Counting masks per layer and per module lets a designer see why foundries pushed EUV, why unidirectional, gridded metal with cuts became standard, and why removing a single routing layer can matter in cost.

## Going deeper

Lin, "Optical Lithography: Here Is Why" (2010); Mack, "Fundamental Principles of Optical Lithography", on double patterning.
