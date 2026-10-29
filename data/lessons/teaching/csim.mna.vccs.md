# Voltage-controlled current source stamps

Controlled sources are how small-signal device models enter MNA: a transistor's transconductance is a voltage-controlled current source. The VCCS stamp is a four-entry pattern that couples the output nodes' rows to the controlling nodes' columns. It needs no new unknown, but it breaks the symmetry that resistors alone produce, which is one reason SPICE uses a general LU factorization.

## Four entries off the controlling nodes

The card G1 n+ n- nc+ nc- gm drives a current $g_m(v_{nc+} - v_{nc-})$ from n+ through the source to n-. That current leaves node n+, so in the KCL row of n+ it appears with a plus sign: $+g_m$ at column nc+ and $-g_m$ at column nc-. It enters node n-, so row n- gets $-g_m$ at nc+ and $+g_m$ at nc-. Any row or column belonging to ground is dropped. For G1 1 3 2 0 3m, the stamp is $+3$ mS at (1,2) and -3 mS at (3,2); the nc- column is ground and vanishes. A VCCS adds no unknown because its current is already expressed in node voltages.

## Adding a VCCS to a resistive matrix

Stamps accumulate, so an entry of the full matrix is the resistor contribution plus any controlled-source contribution at that position. If a 1 k$\Omega$ resistor joins nodes 2 and 3 and a VCCS G1 3 2 1 2 4.5m is present, entry (3,2) gets -1 from the resistor and -4.5 from the VCCS (row n+ = 3, column nc- = 2), so -5.5 mS, while (3,1) gets only $+4.5$ mS. When a VCCS output node coincides with one of its controlling nodes, the stamp lands on a diagonal and behaves like a conductance, positive or negative. A negative effective conductance is the signature of gain and positive feedback.

## Controlled sources break symmetry

Resistor stamps are symmetric, so a passive resistive network has $G = G^T$, and its matrix is positive definite once grounded. A VCCS puts entries in output rows and control columns without the mirror image, so $G_{ab} - G_{ba}$ equals the controlled-source contribution alone. For G1 4 3 2 1 3m, $G_{42} = +3$ mS but $G_{24}$ has no VCCS entry, so $G_{24} - G_{42} = -3$ mS. The loss of symmetry rules out Cholesky factorization and means the adjoint system uses $G^T$, not $G$. SPICE-class simulators therefore use sparse LU with pivoting chosen for both stability and sparsity.

## Going deeper

Vlach and Singhal, "Computer Methods for Circuit Analysis and Design" (2nd ed., 1994), chapter on MNA; Najm, "Circuit Simulation" (Wiley, 2010).
