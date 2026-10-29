# Labels from tool runs

In EDA, a training label is almost never a human annotation. It is the output of a tool: the overflow a global router reports, the slack signoff STA computes, the violation count a DRC engine returns. Labels are therefore expensive, tied to one tool and setting, and correlated within a design. A tool designer must cost them before deciding a learned model is worth building.

## The cost of a label

Each label set costs a flow run. A congestion dataset of 20 designs with 10 parameter variants each needs 200 placement-and-route runs. At 6 machine-hours per run that is 1200 machine-hours; on 25 machines running in parallel the wall-clock time is $1200/25 = 48$ hours, or 2 days. Timing labels need routing, extraction and STA; DRC labels need detailed routing and a full rule check, which are slower still. Labels also inherit the tool's version and settings. If the router changes, labels made with the old router describe a different tool, and a model trained on them predicts the wrong outcome.

## Samples per run

One run yields many labels. A 400 by 400 um die cut into 4 um gcells has $100 \times 100 = 10{,}000$ gcells, so 200 runs give 2 million labeled gcells. The number is misleading. Gcells of one run share a netlist, a floorplan and a placement; gcells of variants of one design share most of that too. For generalization to a new design, the effective sample size is closer to the number of distinct designs than to the number of gcells. This is why random tile-level splits inflate test scores, and why datasets are judged by design diversity, not tile count.

## Break-even reuse

A predictor has a fixed cost, labels plus training, and a per-query saving: the tool run it replaces minus its own inference time. Suppose labels cost 1200 machine-hours and training 50 hours, so 1250 hours in all, and each query takes 0.05 hours instead of a 6-hour run. Each query saves 5.95 hours, and the smallest number of queries with net savings is $\lfloor 1250/5.95 \rfloor + 1 = 211$. If the exploration loop issues thousands of queries per project, the model pays. If it is used a few dozen times, or if every query still needs a confirming tool run, it never does.

## Going deeper

The CircuitNet dataset papers (Chai et al.) describe how flow runs become features and labels; the OpenROAD-flow-scripts documentation shows the runs such labels come from.
