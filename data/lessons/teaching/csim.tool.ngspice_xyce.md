# Running decks in ngspice and Xyce

ngspice and Xyce are the two open-source simulators a tool designer is most likely to read, run and extend. ngspice descends from Berkeley SPICE3 with many extensions; Xyce was written at Sandia National Laboratories for large parallel simulations. Both read SPICE decks, both expose their numerical machinery through analysis cards and options, and both write raw files that scripts can parse. Knowing which knob reaches which part of the engine is how you debug a run instead of guessing.

## Analysis cards and what they compute

Each dot card requests an analysis. .op solves the DC operating point by Newton. .dc sweeps a source, re-solving the operating point at each value from the previous solution. .ac linearizes all devices at the operating point and solves the complex system $(G + j\omega C)x = b$ over a frequency sweep. .noise is a small-signal sweep that sums device noise densities to an output and refers them to an input. .tran tstep tstop integrates in time; tstep is the print and initial suggested step, not a fixed step, because the internal step is chosen by truncation-error control. A card .tran 1p 2n therefore asks for output every picosecond up to 2 ns, while the engine may take far fewer or more internal steps.

## Options that steer the numerical engine

The .options card reaches the solver directly. reltol (default $10^{-3}$), abstol (default 1 pA) and vntol (default 1 uV) set the Newton convergence tests and, through reltol, the truncation-error tolerance used in step control; trtol scales how strictly truncation error is enforced. gmin is the minimum conductance added across junctions and used by gmin stepping. itl1 and itl4 limit Newton iterations for the operating point and for each transient step. method=gear switches integration from trapezoidal to backward differentiation formulas, with maxord setting the order. Loosening reltol to $10^{-2}$ trades accuracy for speed everywhere; setting abstol to 1 uA on a nanoampere circuit silently accepts errors larger than the signal.

## Reading output and choosing a simulator

In batch mode, ngspice -b deck.cir -r out.raw runs without a user interface and writes every vector (time, node voltages, branch currents) to a raw file that post-processing reads. .measure statements compute scalar results inside the run, such as the delay between 50% crossings of two waveforms, so a characterization script needs only the log. When two simulators disagree, compare their default options and convergence strategies first: different gmin, tolerances and homotopy choices change results more often than bugs do. Xyce was designed for parallel execution with distributed-memory linear solvers, which matters for very large post-layout decks, and it offers homotopy methods for hard operating points.

## Going deeper

The ngspice manual; the Xyce Users' Guide and Reference Guide (Sandia National Laboratories).
