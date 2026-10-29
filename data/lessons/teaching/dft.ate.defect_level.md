# Defect level and the Williams-Brown model

Customers buy quality, not coverage. The defect level, the fraction of shipped parts that are defective, is usually quoted in defective parts per million (DPPM), and it depends on two quantities a test engineer knows: process yield and fault coverage. The Williams-Brown model links them and turns a DPPM target into a coverage target, which is how coverage sign-off thresholds are justified.

## The model

Let $Y$ be the yield (fraction of good dies) and $T$ the fault coverage of the test. Williams and Brown modeled faults as independent and equally likely, so a die is good with probability $Y$ and a die passes the test if none of its faults is in the detected fraction. The probability that a passing part is defective is
$$DL = 1 - Y^{1 - T}.$$
At $T = 1$, $DL = 0$; at $T = 0$, $DL = 1 - Y$, the fraction of all dies that are bad. For $Y = 0.4$ and $T = 0.9324$, $DL = 1 - 0.4^{0.0676} = 0.0601$, about 60,000 DPPM: even 93 percent coverage ships many bad parts when yield is low.

## Coverage for a target

Inverting the model gives the coverage needed for a quality goal:
$$T = 1 - \frac{\ln(1 - DL)}{\ln Y}.$$
For $Y = 0.95$ and 500 DPPM, $T = 1 - \ln(0.9995)/\ln(0.95) = 99.025$ percent. Lower yield demands higher coverage for the same DPPM, which is why new processes, with low yield, are the hardest to ship at quality. Automotive targets near single-digit DPPM push coverage requirements toward the limit of the stuck-at model and motivate additional fault models.

## Connecting to defect density

With Poisson random defects, $Y = e^{-AD_0}$ for die area $A$ and defect density $D_0$, and then $Y^{1-T} = e^{-AD_0(1-T)}$: only the untested fraction of the defects reaches the customer. A 56 mm$^2$ die at 0.37 defects per cm$^2$ has $Y = 0.813$; at 97.94 percent coverage, $DL \approx 4260$ DPPM. The model assumes fault coverage measures defect coverage, which overstates quality when real defects behave unlike modeled faults; refinements weight faults by their likelihood or use multiple-detect coverage.

## Going deeper

Williams and Brown, "Defect level as a function of fault coverage," IEEE Transactions on Computers, 1981; Bushnell and Agrawal (2000), chapter 3.
