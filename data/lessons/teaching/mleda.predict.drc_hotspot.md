# DRC hotspot prediction

A DRC hotspot predictor flags gcells likely to contain design-rule violations after detailed routing, so the placer or router can act early. It is a classification problem with extreme class imbalance: violations sit in a tiny fraction of gcells. That makes accuracy meaningless and makes precision depend on how rare hotspots are in the design at hand.

## Precision and recall

Compare the flagged gcells with the DRC engine's report. True positives (TP) are flagged gcells with violations, false positives (FP) are flagged clean gcells, and false negatives (FN) are violating gcells not flagged. Precision is $TP/(TP + FP)$, the fraction of flags that are real; recall is $TP/(TP + FN)$, the fraction of real hotspots caught. With 40 true flags, 60 false alarms and 20 misses, precision is 40 percent and recall 66.7 percent. Accuracy would count the many true negatives too: on 250,000 gcells, flagging nothing scores over 99.9 percent accuracy while catching no hotspot at all.

## The F1 score

F1 is the harmonic mean of precision and recall, $2PR/(P + R)$, which equals $2TP/(2TP + FP + FN)$. In the example, $F_1 = 80/(80 + 60 + 20) = 0.5$. True negatives never appear, so F1 is not inflated by the sea of clean gcells. The harmonic mean is pulled toward the smaller of the two, so a model cannot score well by maximizing one and ignoring the other. F1 still hides the trade-off a tool designer cares about: a router may tolerate many false alarms (extra spreading costs a little wirelength) but not misses, in which case recall at a fixed false-alarm budget is the better target.

## Precision depends on prevalence

Recall and the false-positive rate are properties of the classifier; precision also depends on how common hotspots are. Per gcell, the expected true flags are recall times prevalence and the false flags are the false-positive rate times one minus prevalence. A classifier with recall 0.9 and false-positive rate 0.01 on a design where 0.5 percent of gcells are hotspots gives true flags 0.0045 and false flags 0.00995, so precision is $0.0045/0.01445 = 31$ percent. The same classifier on a design with 5 percent hotspots reaches 83 percent. A precision reported on one design does not transfer to a cleaner one.

## Going deeper

Huang et al., "Machine Learning for Electronic Design Automation: A Survey," ACM TODAES (2021), section on DRC prediction; the CircuitNet dataset papers (Chai et al.) for public DRC-hotspot labels.
