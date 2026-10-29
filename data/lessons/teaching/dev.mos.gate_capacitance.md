# Gate capacitance by operating region

The gate capacitance of a transistor is the load every driver sees and the charge every switching event moves. It changes with operating region: the channel charge is shared between source and drain differently in the linear region, in saturation and in cutoff. Cell characterization, Liberty input-capacitance values and circuit-simulator charge models all encode this region dependence.

## Gate capacitance in the linear region

In strong inversion the channel is a conducting sheet under the whole gate, and the gate-to-channel capacitance is $WLC_{ox}$. The Meyer model assigns this capacitance to the source and drain terminals. With $V_{DS}$ near zero the channel is symmetric, so each gets half:
$$C_{gs} = C_{gd} = \tfrac12 WLC_{ox} + C_{ov}W,$$
where $C_{ov}$ is the overlap capacitance per unit width at each edge, from the gate overlapping the source and drain extensions plus fringing fields.

Example: $W = 0.5$ µm, $L = 50$ nm, $C_{ox} = 20$ fF/µm$^2$, $C_{ov} = 0.2$ fF/µm. $WLC_{ox} = 0.5$ fF, so $C_{gs} = 0.25 + 0.1 = 0.35$ fF.

At short channel lengths the overlap and fringe terms are a large fraction of the total, which is why extraction of fringe capacitance matters as much as the intrinsic term.

## Gate capacitance in saturation

In saturation the channel is pinched off near the drain, and the drain voltage no longer changes the channel charge. The Meyer model gives
$$C_{gs} = \tfrac23 WLC_{ox} + C_{ov}W, \qquad C_{gd} = C_{ov}W.$$
The two-thirds comes from integrating the channel charge profile of the square law.

Example: same device as before. $C_{gs} = \tfrac23 \times 0.5 + 0.1 = 0.433$ fF, and $C_{gd} = 0.1$ fF.

In cutoff the channel disappears and the gate couples to the body through the depletion region, a series combination of $C_{ox}$ and $C_{dep}$. The Meyer model is a set of capacitances, not derivatives of a single charge function, so it can violate charge conservation in transient simulation. Modern compact models, BSIM included, compute terminal charges first and derive capacitances from them, which guarantees conservation.

## The Miller effect

The gate-to-drain capacitance connects input and output. If the drain swings by $-A_v$ times the gate swing, the charge through $C_{gd}$ is $(1 + |A_v|)$ times what it would be with a fixed drain. Seen from the input,
$$C_{in} = C_{gs} + (1 + |A_v|)C_{gd}.$$

Example: $C_{gs} = 0.433$ fF, $C_{gd} = 0.1$ fF and $|A_v| = 10$. $C_{in} = 0.433 + 1.1 = 1.53$ fF, more than three times $C_{gs}$.

In digital switching the output moves opposite to the input over a full swing, so $C_{gd}$ counts roughly double in the effective switching load. This is one reason a Liberty file's input capacitance is a characterized value rather than a geometric sum, often given separately for rising and falling transitions. In amplifiers the Miller multiplication sets the dominant pole and limits bandwidth.

## Going deeper

Tsividis and McAndrew, "Operation and Modeling of the MOS Transistor" (3rd ed.), chapter 8 on charge modeling; Ward and Dutton, "A charge-oriented model for MOS transistor capacitances", IEEE Journal of Solid-State Circuits (1978).
