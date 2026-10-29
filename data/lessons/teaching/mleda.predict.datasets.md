# Datasets for ML in EDA

Progress on learned EDA models depends on shared data, and EDA data is scarce: designs are proprietary, labels need expensive tool runs, and results depend on the node and the flow. Open datasets such as CircuitNet package flow runs as features and labels. Knowing how they are built tells a tool designer what a reported result on them does and does not show.

## What a dataset contains

An ML-for-EDA dataset pairs inputs a tool sees early with labels a tool computes later. CircuitNet, for example, provides layout-derived feature maps (cell density, pin density, RUDY, macro regions) with labels such as routing congestion, DRC violations and IR drop, produced by running designs through physical-design flows under many settings. Synthesis datasets such as OpenABC-D pair and-inverter graphs with the QoR that ABC reports after many recipes. The input representation decides the task's form: grids give image-to-image problems, netlists give graph problems. Each label is only as trustworthy as the run that produced it, so the dataset's flow, node and tool versions belong with the data.

## How splits mislead

Because labels are costly and open designs few, datasets grow by running variants of a small number of designs. Samples from one design share netlist structure, floorplan and flow style. A random split of tiles or runs puts the same design in training and test, so the test score measures recognition of seen designs. The honest split holds out whole designs, and the effective sample size is the number of distinct designs. Labels also carry their tool: a model trained on labels from one router, rule deck and node predicts that combination, and must be re-evaluated before use with another.

## Comparable results

A benchmark is the data, the split and the metric together. Two papers that use the same dataset with different splits, metrics or preprocessing report incomparable numbers. Reproducibility also needs the flow scripts and parameters that produced each label, so that others can regenerate and audit them. A tool designer reading a result on an open dataset should ask which node and PDK the labels came from, whether the split held out designs, and whether the gain was measured in a flow, such as improved routed congestion, or only as prediction error.

## Going deeper

The CircuitNet dataset papers (Chai et al.); Chowdhury et al., "OpenABC-D: A Large-Scale Dataset for Machine Learning Guided Integrated Circuit Synthesis" (2021).
