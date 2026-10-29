# Mesh analysis

Mesh analysis is the dual of nodal analysis: it takes a current circulating in each window of a planar circuit as the unknown and writes KVL around each window. Simulators rarely use it, but it is the natural way to reason about loops, which is how inductance, return paths and current crowding are understood. It also shows how a shared element couples two unknowns.

## Writing a mesh equation

Assign a clockwise mesh current to each window of a planar circuit. Walking around a mesh in the direction of its current, each resistor contributes a drop of $R$ times the net current through it, and each source a rise or a drop depending on its polarity. For a single loop there is one unknown.

Example: a loop holds a 5 V source aiding the clockwise direction, a 2 V source opposing it, and resistors of 1 kohm and 2 kohm. KVL: $5 - 2 = (1 + 2) I$, so $I = 1$ mA. A source aiding the mesh current appears with a plus sign on the source side of the equation; an opposing one with a minus sign. Volts over kilohms give milliamperes, which keeps on-chip numbers readable.

## Solving two meshes

When two meshes share a resistor $R_3$, the current through it is the difference of the mesh currents. Writing KVL around each mesh gives a symmetric 2 by 2 system: the diagonal holds each mesh's total resistance and the off-diagonal holds minus the shared resistance:
$$(R_1 + R_3)I_1 - R_3 I_2 = V_1, \qquad -R_3 I_1 + (R_2 + R_3) I_2 = V_2,$$
with each $V$ the net source rise in its mesh's direction.

Example: $R_1 = 1$, $R_2 = 2$, $R_3 = 1$ kohm, a 4 V source aiding $I_1$, and a 1 V source aiding $I_2$. Then $2I_1 - I_2 = 4$ and $-I_1 + 3I_2 = 1$. Eliminating, $I_1 = 13/5 = 2.6$ mA and $I_2 = 1.2$ mA. The structure mirrors the nodal conductance matrix with resistances in place of conductances, and it stays symmetric because the shared resistor appears with the same sign in both equations.

## Branch currents from mesh currents

Mesh currents are fictitious loop variables; real branch currents are recovered from them. A branch on the outside of a single mesh carries that mesh's current. A shared branch carries the difference, $I_1 - I_2$ in the direction of $I_1$ when both meshes are clockwise.

Continuing the example with $I_1 = 2.6$ mA and $I_2 = 1.2$ mA, the shared 1 kohm resistor carries 1.4 mA in the direction of $I_1$ and drops 1.4 V. A sign error here is the most common mistake: the mesh currents cross a shared branch in opposite directions. This difference-of-loops picture is exactly how loop and partial inductances combine: the current in a shared return conductor is the sum or difference of the loop currents that use it, and its inductance couples those loops.

## Going deeper

Hayt, Kemmerly and Durbin, "Engineering Circuit Analysis", chapter 4; Desoer and Kuh, "Basic Circuit Theory", chapter 10.
