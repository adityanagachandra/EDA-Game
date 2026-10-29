# Memory fault models

Embedded memories occupy a large share of a modern die and are built at the tightest design rules, so they collect a large share of defects. Their dense, analog-like arrays fail in characteristic ways that gate-level stuck-at faults do not describe. Memory test therefore uses functional fault models, defined on cells, decoders and read/write logic, and March tests are designed against exactly these models.

## Cell faults

A **stuck-at fault** (SAF) leaves a cell permanently at 0 or 1. A **transition fault** (TF) lets a cell hold either value but blocks one transition: an up transition fault cannot go from 0 to 1, a down fault cannot go from 1 to 0. A **data retention fault** (DRF) loses a stored value after a pause without access, from leakage; detecting it needs a delay before the read. A **read destructive fault** returns the correct value but flips the cell, so a second read is needed to see it; a **deceptive read** variant returns the right value but leaves the cell damaged for later reads.

## Coupling faults

Coupling faults involve an **aggressor** cell $a$ and a **victim** $v$. In an **inversion** coupling fault (CFin), a transition written into $a$ (rising or falling, per fault) inverts $v$. In an **idempotent** coupling fault (CFid), a transition in $a$ forces $v$ to a fixed value. In a **state** coupling fault (CFst), $v$ is forced whenever $a$ holds a particular value, with no write event needed. Coupling grows as cells shrink, because neighbours share closer wires and wells. Because a victim can sit at a higher or lower address than its aggressor, tests must sweep addresses in both directions.

## Decoder and other faults

An **address decoder fault** (AF) makes an address select no cell, several cells, or a cell also selected by another address; writing address 5 might also write address 13. **Neighbourhood pattern-sensitive faults** make a cell fail only when its physical neighbours hold particular values; testing them needs the physical address scrambling and many more operations than a linear March test. DRAM adds capacitor leakage, refresh and sense-margin faults that a latch-based SRAM cell does not have. Each model is a target that fixes what a test must do: write, read, wait, and in which address order.

## Going deeper

van de Goor, "Testing Semiconductor Memories: Theory and Practice" (1991); Bushnell and Agrawal (2000), chapter 9; Hamdioui, "Testing Static Random Access Memories" (2004).
