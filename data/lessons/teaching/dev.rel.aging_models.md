# Aging-aware device models

Aging-aware design turns device degradation physics into numbers that timing and power tools can use. A reliability simulator applies BTI and HCI shifts to compact-model parameters according to each device's stress, and characterization produces aged libraries for end-of-life signoff. The modeling chain, from mission profile to aged Liberty file, is where device physics meets EDA most directly in reliability.

## From stress to aged models

An aging flow has three parts.

A mission profile describes how the product is used: supply voltages, temperatures, operating modes and their fractions of life, often over ten years.

Device stress comes from simulation or signal statistics: for BTI, the fraction of time each transistor's gate is under stress (static probability) and the recovery between stresses; for HCI, the switching activity and the drain voltage during transitions.

Aging models then compute parameter shifts per device, mostly in threshold and mobility, using fitted equations in stress time, voltage, temperature and duty cycle. Reliability-aware simulators apply these shifts to the base compact model, leaving its I-V equations unchanged, and re-simulate to see the degraded behavior. Compact-model aging extensions standardize this interface.

Characterizing standard cells with the shifted models produces aged Liberty libraries, which static timing analysis uses exactly like fresh ones.

## Timing consequences

Aging mostly slows transistors, so setup checks are usually worst with aged libraries. But aging is not uniform.

Different devices see different stress. A clock-tree branch whose buffers idle at one logic level for long periods ages differently from a branch that toggles, so clock skew drifts over life. Skew is a difference, so asymmetric aging can create new setup or hold violations that time-zero analysis never saw.

Hold checks depend on relative delays and can be worst with fresh, fast devices or when aging slows the capture clock more than the data path. Signoff therefore uses both fresh and aged views.

A single worst-case aged library assumes every cell was maximally stressed, which is safe but pessimistic. Instance-specific or activity-aware aging analysis recovers margin by aging each cell according to its own stress. The cost is a much larger characterization and analysis problem, a natural target for fast surrogate models.

## Model risks

Aging predictions rest on extrapolation in three directions: from short stress times to a decade, from high stress voltage to use voltage, and from high temperature to use temperature. A power-law exponent fitted over a few decades of time and extrapolated two more magnifies fitting error; mechanisms may also change outside the stressed window.

Recovery complicates the profile: the same total stress time gives different shifts depending on how stress and relaxation alternate. Static probability captures much of this for logic, but high-frequency switching and long idle periods behave differently.

Temperature must be the local channel temperature. FinFET and nanosheet channels self-heat because heat escapes poorly through narrow fins and sheets, so devices age as if hotter than the ambient.

Static circuits have their own failure modes: an SRAM cell storing one value for years stresses one side of its cross-coupled pair, and the growing mismatch erodes its stability margin. Adaptive voltage scaling that raises supply to compensate aging also accelerates aging, so its schedule must be checked against the models.

## Going deeper

Reddy and others, "Impact of negative bias temperature instability on digital circuit reliability", IRPS (2002); the Si2 Compact Model Coalition documentation on the open model interface for reliability (OMI) and aging.
