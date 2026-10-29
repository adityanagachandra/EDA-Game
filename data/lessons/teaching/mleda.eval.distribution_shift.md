# Distribution shift across designs and nodes

A learned EDA model meets designs, technology nodes and tool versions that differ from its training data. Its accuracy on a held-out slice of the training distribution says little about accuracy after such a shift. Measuring shift honestly, detecting it at run time and falling back to exact engines are what make a predictor safe to put in a tool.

## Kinds of shift

Shift comes from every input to the label. A new design style, such as a network-on-chip with long wires after training on CPU blocks, changes the relation between features and slack. A new node changes routing resources per gcell and the rules that turn demand into violations. A new power-grid style changes how currents become IR drop. A changed upstream tool version shifts the features themselves, a covariate shift, even when labels are unchanged. A tighter clock shifts slack labels for the same logic. Each requires measuring the model again; none is visible from validation error.

## Measuring generalization by design

A random split of endpoints or tiles across all designs puts the same designs in training and test, so it estimates error on seen designs. Leave-one-design-out trains on all designs but one and tests on the held-out one, repeated for each design. With five designs whose held-out MAEs are 23.6, 14.9, 12.0, 12.0 and 17.7 ps, the estimate for a new design is their mean, 16.0 ps, often several times the random-split figure. The spread across folds matters as much as the mean: a tool meets designs one at a time, and the worst fold shows how bad an unlucky one can be.

## Living with shift

A deployed model should detect when it is out of distribution: input features or embeddings far from the training data, or high ensemble disagreement, are available before the expensive label. When shift is detected, the tool falls back to the accurate engine for that design. Fine-tuning on a few runs of a new design is legitimate adaptation, provided accuracy is reported on held-out runs or regions of that design. Above all, if exact engines gate every accepted result, a shifted model costs runtime or quality, never a wrong signoff answer.

## Going deeper

Quinonero-Candela et al., eds., "Dataset Shift in Machine Learning" (2009); the CircuitNet dataset papers (Chai et al.), on cross-design evaluation.
