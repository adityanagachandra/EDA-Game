# The DTCO evaluation loop

Design-technology co-optimization is an engineering loop, not a single calculation. Candidate process options and cell architectures are turned into rules and models, then into libraries, then into placed and routed benchmark blocks, whose power, performance, area and cost decide which option survives. Every engine in the EDA flow runs inside the loop, which is why tool designers are part of technology definition.

## What DTCO optimizes

The variables are coupled design and technology knobs: CPP, metal pitches, track height, fins or sheets per device, diffusion-break style, rail scheme, contact architecture, and the rules each implies. The objective is block-level power, performance and area at a target, plus wafer cost and yield. Device metrics alone mislead: a faster transistor that needs a taller cell or adds parasitics can lose at the block level. Early in a node's definition there is no signoff PDK, so the loop uses a pathfinding or predictive PDK with projected rules and compact models, and it updates them as process data arrives.

## The loop step by step

For each candidate: write the rules and technology files; generate or draw cell layouts and check them; characterize the cells from the models into timing and power libraries; synthesize, place and route several benchmark blocks; extract and time them; record area, frequency, power and routability. Then compare candidates at matched targets and feed the mechanisms behind wins and losses back into the next candidates. Several benchmarks are needed because a CPU core, a GPU tile and a memory controller stress pin density and wiring differently. A cell that is 12% smaller can give a block only 3% smaller if fewer tracks forced lower utilization.

## From DTCO to STCO

Interpreting results takes care: cell-level metrics such as NAND2 area are proxies that can be optimized at the expense of routability. When a candidate loses frequency on one block, DTCO looks for compensating knobs and reruns rather than accepting or discarding it outright. System-technology co-optimization (STCO) widens the scope to how a system is partitioned across dies, packages and 3D stacks, which technology each part uses, and how memory is attached. The arithmetic of cost per good die, yield versus area, and interface overhead then joins the PPA results, as the next concept on chiplets shows.

## Going deeper

Clark et al., "ASAP7: A 7-nm finFET predictive process design kit" (2016); the OpenROAD project documentation, as an open flow used for DTCO-style experiments.
