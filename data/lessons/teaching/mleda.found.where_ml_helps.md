# Where machine learning helps in EDA

Machine learning is one branch of EDA, not a replacement for it. Placers, routers, timers and verifiers are built on algorithms with guarantees and on physics with known error. A learned model is useful when it gives a fast, approximate answer that steers one of those engines, and harmful when it is trusted where a guarantee is needed. A tool designer's first job is to tell the two cases apart.

## Where learning pays

Learning pays when three conditions hold. The exact evaluation is expensive or not yet available: post-route wire delay is unknown at placement time, and a full global route per placement move is too slow. Many similar instances recur, so experience transfers: every block in a project runs the same flow with similar netlists. And an approximate answer is good enough because something exact runs later. Under those conditions a model can rank 200 flow configurations from early features so that only ten get full runs, or give the placer a congestion map between expensive routes. Where the exact answer is cheap (the HPWL of a net, an overlap check, a topological order), a model only adds error and runtime.

## What learning cannot replace

Some questions have answers that must hold for every input, and a model trained on samples cannot give them. Legality is one: a placement either has overlapping cells or it does not, and a legalizer enforces that exactly. Manufacturability is another: only a DRC engine running the real rule deck on the real geometry decides it. Functional equivalence of two netlists is a formal property proven by an equivalence checker. Signoff timing is a conservative bound computed by propagating worst-case delays over every path. A predictor's error is measured on a test set and has no guaranteed bound on an unseen design. So a prediction of zero DRC violations, or of positive slack, licenses nothing at signoff.

## Propose with a model, decide with an engine

The pattern that makes learning safe is simple: the model proposes and an exact engine decides. A model suggests 50 buffer insertions; incremental STA checks each one. A learned net ordering guides a router; the router still enforces design rules. In the four-layer engine model the learned part lives in the strategy layer, feeding guidance to a core that keeps invariants, so the worst a bad model can do is a legal but poor result. The other half of the decision is cost. A method that matches a heuristic's quality but needs three GPU-days of training per design is useful only if that cost is amortized across many designs.

## Going deeper

Huang et al., "Machine Learning for Electronic Design Automation: A Survey," ACM TODAES (2021); Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), for the engines a model would assist.
