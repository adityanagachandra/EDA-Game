# Grounding, shielding and the proximity effect

Electromagnetic compatibility is about keeping a design from interfering with itself and with its surroundings. Three ideas recur: grounding topology decides where return currents flow, shielding decides how much field escapes or enters an enclosure, and the proximity effect decides how current crowds in nearby conductors. Tool builders encode them as layout rules, extraction models and design checks.

## Ground loops and grounding topologies

A ground loop forms when two circuits share ground through more than one path. A changing magnetic field through the loop, or a voltage difference between the two ground points, drives a circulating current that flows through the shared conductors and appears as noise on signals referenced to them. Hum at the power-line frequency on a bench setup is the classic symptom.

Single-point (star) grounding avoids loops by connecting each circuit's ground to one point, which works at low frequency where resistance dominates. At high frequency, long single-point leads are inductive and resonant, so high-speed design uses multipoint grounding to a solid plane, keeping every return path short and close to its signal. Mixed-signal designs often combine the two, with a carefully chosen connection between analog and digital grounds. A slot in a reference plane breaks the high-frequency return path and forces a detour.

## How a shield works

A conducting shield attenuates an incident field by two mechanisms. Reflection loss comes from the mismatch between the wave impedance of the incident field and the shield's very low intrinsic impedance; for a plane wave in air (377 ohm) hitting a good conductor, most of the energy reflects at the surface. Absorption loss comes from the field decaying inside the metal as $e^{-t/\delta}$, about 8.7 dB per skin depth of thickness $t$.

At low frequency the skin depth is large (about 2 mm in copper at 1 kHz), so thin shields absorb little, and low-impedance magnetic fields reflect poorly; thick or high-permeability materials are needed. At high frequency a thin shield is ample, and real enclosures are instead limited by apertures and seams: a slot whose length approaches half a wavelength radiates efficiently. Shield effectiveness is therefore usually set by the openings, not by the metal.

## The proximity effect

When conductors carrying alternating current are close together, each sits in the magnetic field of the others, and eddy currents redistribute the current density. For two conductors carrying opposite currents, current crowds onto the facing surfaces; for currents in the same direction, onto the far surfaces. The crowding minimizes the loop's magnetic energy, and it shrinks the effective conducting area beyond what the skin effect alone would.

Consequences for tools: the AC resistance of a closely spaced go-and-return pair is higher than the isolated-wire skin-effect value; multi-layer windings and tightly packed bus lines show resistance far above their strand count; and the return current in a plane concentrates directly under its trace. Extractors capture this by subdividing conductors into filaments and solving for the current distribution, which a uniform-current model cannot represent.

## Going deeper

Ott, "Electromagnetic Compatibility Engineering", chapters 3 and 6; Paul, "Introduction to Electromagnetic Compatibility" (2nd ed.), chapter 10 on shielding.
