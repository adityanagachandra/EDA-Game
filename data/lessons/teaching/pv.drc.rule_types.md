# Design rule types

Design rules are the foundry's contract with the designer: if every rule passes, the layout is manufacturable with acceptable yield. A rule deck holds thousands of them, but they fall into a small set of types. A DRC engine designer needs that taxonomy because each type maps to a different geometric computation: single-shape measurements, relations between two layers, or checks that depend on context such as connectivity or a window.

## Single-layer rules

Single-layer rules measure one layer. Minimum width is the distance between a shape's own opposite edges; minimum spacing is the distance between facing edges of different shapes; a notch is spacing between facing edges of the same polygon; minimum area bounds the area of each polygon, so small islands (via landing pads) do not lift off. For example, with a 40 nm width rule a 38 nm wire fails, and with a 32 nm spacing rule two wires 30 nm apart fail. All of these operate on the merged layer, after overlapping shapes are combined, and they reduce to edge-pair distance tests or polygon area computations.

## Two-layer rules

Two-layer rules relate shapes on different layers. Enclosure requires an outer layer to surround an inner one by a margin, such as metal around a via cut. Extension requires one layer to continue past another's edge, such as poly past diffusion so the gate end does not shrink (the endcap). Inter-layer spacing keeps, for instance, a contact away from a neighbouring gate. Cut layers have their own spacing and size rules. Engines implement these with Boolean operations and sizing: an enclosure check is a test that the inner layer, grown by the margin, stays inside the outer one.

## Context-dependent rules

Some rules depend on more than local geometry. End-of-line spacing applies more spacing in front of narrow wire ends, because line ends print short and round. Width- and run-length-dependent spacing raises spacing for wide wires that run side by side. Density rules constrain the metal fraction in every window of a given size, which no single shape determines. Antenna rules depend on which metal is electrically connected to a gate at an intermediate manufacturing step, so the engine needs connectivity, not just geometry. These are the expensive checks, and they are where engines differ most.

## Going deeper

A foundry design rule manual for an open PDK, such as the SkyWater SKY130 or GF180MCU documentation.
