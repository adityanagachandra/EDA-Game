# Rule decks and derived layers

A rule deck is a program. It names the stream's layers, derives new layers by Boolean and sizing operations, and runs checks on them, each tagged with a rule name. Most of a deck's complexity sits in derived layers that express the physical objects a rule talks about: gates, source and drain, wide metal, via landings. Engine designers optimize deck execution as a dataflow graph of layer operations.

## Derived layers by Boolean operations

Derived layers come from AND, OR, NOT and XOR of drawn layers. GATE = POLY AND ACTIVE is the transistor channel; SD = ACTIVE NOT POLY the source and drain; FIELDPOLY = POLY NOT ACTIVE the poly used as wire. For POLY $[54, 71] \times [11, 134]$ and ACTIVE $[38, 150] \times [44, 99]$, GATE is $[54, 71] \times [44, 99]$ with area 935 square nm, and SD has area $112 \cdot 55 - 935 = 5225$. A rule written on GATE fires only on channels, not on poly wiring, which is the point of deriving it.

## Sizing inside a deck

Sizing grows or shrinks every shape by a distance, and decks combine it with Booleans. HALO = (M1 SIZED BY 15) AND M2 finds metal 2 within 15 nm of metal 1. For M1 $[23, 81] \times [2, 52]$ grown to $[8, 96] \times [-13, 67]$ and M2 $[3, 63] \times [45, 72]$, HALO is $[8, 63] \times [45, 67]$, area 1210, while the unsized overlap is only 280. Sizing down then up (an opening) removes features narrower than twice the distance, which is how decks find narrow necks and slivers.

## What a rule deck is

A deck contains layer definitions (stream layer and datatype to name), derived-layer operations, check statements with rule names and thresholds, and switches for process options. It is written in the language of one tool (Calibre's SVRF, KLayout's Ruby-based DSL, Magic's technology file), so a foundry maintains one deck per supported tool and qualifies each against test layouts. Rule names key everything downstream: reports, waivers, and automated fix strategies. The engine executes the deck as a graph: shared derived layers are computed once and freed after last use, and independent checks run in parallel.

## Going deeper

The KLayout DRC reference manual; the DRC decks distributed with the SkyWater SKY130 open PDK.
