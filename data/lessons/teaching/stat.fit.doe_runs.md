# Factorial designs and run counts

Every characterization study, every process-sensitivity experiment and every surrogate model is built from a set of simulation or silicon runs. Design of experiments decides which runs, so that the effects of interest can be estimated with the fewest. The first decisions are how many factors, how many levels, and how much of the full combination space to run, and each has a precise run count.

## Full factorial designs

A full factorial runs every combination of factor levels: $k$ factors at $L$ levels need $L^k$ runs. Two-level designs ($2^k$) are the workhorse because they estimate every main effect and every interaction for linear-in-each-factor responses. Five factors at two levels need 32 runs; at three levels, 243. The exponential growth is the whole problem: a library characterization over seven process and environment variables at three levels each would need 2187 runs per arc, before any statistical sampling. Full factorials are affordable only for a handful of factors.

## Fractional factorial designs

A $2^{k-p}$ fractional factorial runs a $1/2^p$ fraction of the full design, chosen so that main effects remain estimable. Seven factors in a $2^{7-4}$ design need only 8 runs. The price is **aliasing**: each estimated effect is the sum of several true effects that the fraction cannot separate. The design's **resolution** summarizes this. In a resolution III design main effects are aliased with two-factor interactions; resolution IV keeps main effects clear of two-factor interactions; resolution V keeps two-factor interactions clear of each other. Screening experiments use low resolution to find the few factors that matter, assuming most interactions are small, then follow up with a richer design on those factors.

## Runs for a quadratic model

Many responses have curvature, so the follow-up is often a second-order model in the important factors. A full quadratic in $k$ factors has an intercept, $k$ linear, $k$ pure quadratic and $\binom{k}{2}$ interaction terms, $(k+1)(k+2)/2$ coefficients in all: 10 for three factors, 21 for five. A design needs at least that many distinct runs. Two-level designs cannot estimate pure quadratic terms, so the central composite design adds $2k$ axial points and some centre runs to a $2^k$ core: $2^k + 2k + n_c$ runs. For $k = 3$ with 4 centre runs that is $8 + 6 + 4 = 18$. The Box-Behnken design is a common alternative that avoids extreme corner runs.

## Going deeper

Montgomery, "Design and Analysis of Experiments" (9th ed., 2017), chapters 6, 8 and 11; Box, Hunter and Hunter, "Statistics for Experimenters" (2nd ed., 2005).
