# The RL macro-placement reproducibility debate

The RL macro placer published in Nature in 2021 became a case study in how ML results in EDA should be evaluated. Independent groups ran the released code on open designs and reported that classical methods matched or beat it; the original authors disputed how the reproductions were run. A tool designer does not need to settle the dispute to learn from what each side measured.

## The original claim

The paper described an RL agent with a graph-based encoder that places macros sequentially, pretrained across many blocks and fine-tuned on new ones, and reported quality comparable to or better than human experts in much less time, evaluated on Google's production blocks. Google released an open-source implementation, Circuit Training. Because the evaluation blocks and parts of the flow were proprietary, outsiders could not rerun the original experiments; they could only run the released method on other designs, which is a different experiment. Results were reported on both proxy cost and post-route metrics, and which of the two a comparison uses changes its conclusion.

## What the critiques measured

An academic assessment ran Circuit Training on open designs and enablements and compared it with simulated annealing, analytical placers, a commercial macro placer and human placements. Its questions were concrete: with comparable inputs and evaluation, does the learned policy beat strong non-learned optimizers; how well does the proxy cost track routed timing, power and wirelength; does placement information produced by a commercial tool and fed to the method account for part of its quality. Proponents replied that the reproductions differed from the published method, for example in pretraining and compute. Both arguments are about accounting: what was run, with what inputs and budget, on which metric.

## Lessons for evaluation

The debate lasted because the evidence was not on shared ground. The fix is a standard for ML-in-EDA claims: open designs and flows that anyone can rerun; strong baselines on the same objective with matched compute; several seeds and several designs, since one run on one design says nothing about variance; downstream metrics, not only the proxy that was optimized; and ablations of every input that could carry information from another optimizer. A team considering a learned macro placer should measure post-route PPA and runtime on its own blocks against its current method, which is the only result that answers its question.

## Going deeper

Mirhoseini et al., "A graph placement methodology for fast chip design," Nature (2021); Cheng, Kahng, Kundu, Wang and Wang, "Assessment of Reinforcement Learning for Macro Placement" (ISPD 2023); Markov, "Reevaluating Google's Reinforcement Learning for IC Macro Placement," Communications of the ACM (2024).
