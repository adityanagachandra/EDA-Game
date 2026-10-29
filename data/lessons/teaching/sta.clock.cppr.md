# Clock reconvergence pessimism removal

Under on-chip variation, a timer takes the launch clock path late and the capture clock path early for setup. But the two paths usually share a trunk from the clock root to a branch point, and a physical buffer cannot be slow and fast at the same moment. Clock reconvergence pessimism removal, CPPR or CRPR, adds back the difference on that shared segment.

## The common point and its arrivals

Trace the launch and capture clock paths from the root. They share buffers until the branch point where they split; this last shared pin is the common point. With min and max delays per buffer, the common segment has an early arrival, the sum of its early delays, and a late arrival, the sum of its late delays. For buffers of 38/46, 37/49 and 34/41 ps (early/late), the common point arrives at 109 ps early and 136 ps late. The setup check uses 136 ps on the launch side and 109 ps on the capture side, although the real edge reaches the common point only once.

## The pessimism credit

The double counting equals the late arrival minus the early arrival at the common point: credit $= \sum_{\text{common}} (d_{late} - d_{early})$. For the segment above it is $136 - 109 = 27$ ps. With a single shared buffer of 49/60 ps, the credit is 11 ps. Only the common segment counts; the branches after the split can genuinely differ, because they are different buffers that may vary independently. A timer computes the credit per path by finding the common point of its launch and capture clocks. For hold, launch is taken early and capture late; the shared segment again cannot be both, so the same credit is added to the hold slack.

## Slack after CPPR

The corrected slack is the raw slack plus the credit. For a 1000 ps period, a launch clock of 173 ps late, a capture clock of 168 ps early, clock-to-Q 43 ps, data 470 ps and setup 30 ps, the raw slack is $1000 + 168 - 30 - 686 = 452$ ps, and adding the 27 ps credit gives 479 ps. Reports show the credit as a clock reconvergence pessimism line in the required-time section. CPPR matters most for paths whose flops share a long tree, which is most local paths, and for hold checks between nearby flops. It also matters in multi-cycle and generated-clock paths, where the common path can include the generating logic. Without it, OCV derating would be prohibitively pessimistic.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 10 on on-chip variation and CRPR; the OpenSTA documentation on CRPR.
