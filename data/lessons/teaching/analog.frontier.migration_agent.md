# Process migration agents

Porting an analog block to a new process is the most automatable analog task: the topology and specs are known, and the work is re-sizing and re-verification. A migration agent maps each device to the new process through gm/Id lookup tables, re-checks the circuit-level consequences such as lost gain and headroom, and verifies the result in SPICE. Its core computations are table lookups.

## Holding gm/Id across processes

The gm/Id methodology gives a migration invariant: keep each device's role by holding its $g_m$ and $g_m/I_D$. Then the current $I_D = g_m/(g_m/I_D)$ is unchanged, and the new width is the current divided by the new process's current density at that $g_m/I_D$ and the chosen length: $W_{new} = I_D/(I_D/W)_{new}$. A device needing 1 mS at $g_m/I_D = 12$ carries 83.3 microamperes; if its density goes from 10 to 25 microamperes per micrometre, its width shrinks from 8.33 to 3.33 micrometres. The agent does this for every device from regenerated tables, a Tier 1 operation in effect: new values, same topology. Specs that depend on $g_m$, such as GBW, are preserved to first order.

## What does not carry over

Holding $g_m/I_D$ does not hold intrinsic gain, capacitances or noise. If a device's $g_m/g_{ds}$ at the chosen point falls from 40 to 20, a stage limited by its $r_o$ loses $20\log_{10}(20/40) = -6.02$ dB of gain. The agent must read these quantities from the same table rows, compute the circuit-level effect and decide on a structural fix: longer channels, which trade speed for gain, or a cascode, which is a topology change beyond parameter edits and needs review. Parasitic and device capacitances change too, so pole positions and phase margin must be re-checked. A migration agent that only re-maps widths and declares success has verified the easy part.

## Headroom at a lower supply

New processes usually lower the supply, and analog headroom shrinks with it. A branch stacking $k$ saturated devices, each needing $V_{DS,sat} \approx 2/(g_m/I_D)$, plus a required output swing, must satisfy $k \cdot 2/(g_m/I_D) + V_{swing} \le V_{DD}$. At 0.9 V supply with 3 stacked devices and a 0.3 V swing, $g_m/I_D \ge 2 \cdot 3/0.6 = 10$. Devices are pushed toward moderate inversion, cascodes may no longer fit, and topologies may need to change entirely. The agent can compute these checks for every branch before simulating, flag the branches that cannot meet them, and only then run SPICE across corners as the golden judge, with a reviewable report of what was re-sized, what was flagged and what the simulator certified.

## Going deeper

Jespers and Murmann, "Systematic Design of Analog CMOS Circuits Using Pre-Computed Lookup Tables" (2017); the agent_first_eda lesson on the harness loop.
