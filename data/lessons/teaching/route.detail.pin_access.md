# Pin access points

Standard-cell pins at advanced nodes are small, crowded shapes on the lowest layers, and reaching them legally is one of the hardest parts of detailed routing. Pin-access analysis enumerates, for each pin, the points where a via or wire can land without violating rules, and chooses combinations that do not block each other. A router that gets this wrong fails locally no matter how good its global plan is.

## On-track access points

An access point is a location on a pin where a via up to the next routing layer can land. With vertical M2 tracks at $x = o + kp$, a via on track $x$ fits inside a pin spanning $[x_1, x_2]$ only if its metal landing fits with the required enclosure $e$ on each side: $x_1 + e \le x \le x_2 - e$.

Example: tracks at $40 + 80k$ nm, $e = 10$ nm and a pin from 320 to 630 nm give the window [330, 620], which holds tracks 360, 440, 520 and 600: four access points. A pin narrower than one pitch plus twice the enclosure may have none at all on track, which is a real failure mode in dense libraries.

## Compatible choices for neighbours

Access points of neighbouring pins interact. Two vias too close together violate via-to-via spacing, and a via on one pin can block the only access of the next. So pin-access engines choose a combination of access points per cell instance: one point for each pin such that every pair is legal.

Example: pin A with an access point at 196 nm and pin B with points at 308 and 364 nm, with vias required at least 112 nm apart. Both pairs (196, 308) and (196, 364) qualify, so 2 combinations remain. Since cells repeat, engines analyse each unique cell and placement context once and reuse the result across all its instances, which is what makes the enumeration affordable.

## Off-track access

When on-track points are too few, an engine also considers off-track positions, for example halfway between tracks. A via there is joined to a nearby track by a short jog, which must itself be DRC-clean. Example: allowing half-track positions in a pin window [305, 425] with tracks every 56 nm turns 3 on-track points into 5 candidates.

Off-track access trades simplicity for success. It raises the number of choices, which helps crowded cells, but every jog is extra metal near other pins, so the engine must check spacing, end-of-line and minimum-area rules on it. TritonRoute's pin-access step generates on-track and off-track candidates in this spirit, validates each against the rules, and prefers clean, on-track ones.

## Going deeper

A. B. Kahng, L. Wang and B. Xu, "The Tao of PAO: Anatomy of a Pin Access Oracle for Detailed Routing" (DAC 2020); the TritonRoute paper (IEEE TCAD, 2021).
