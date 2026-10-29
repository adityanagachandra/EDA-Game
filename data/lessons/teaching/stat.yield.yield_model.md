# Yield as a probability

Yield is the fraction of manufactured dies that work and meet spec, and it is the number that turns variation statistics into cost. Formally it is a probability over the joint distribution of defects and process parameters, and building it up from independent pieces, catastrophic defects and parametric specs, is how a yield model is assembled and how an engine decides where a yield loss comes from.

## Parametric yield

A die can be free of defects and still fail a performance spec. If a performance $P$ is Gaussian with mean $\mu$ and standard deviation $\sigma$ and the spec accepts $L \le P \le U$, the parametric yield is
$$Y = \Phi\!\left(\frac{U - \mu}{\sigma}\right) - \Phi\!\left(\frac{L - \mu}{\sigma}\right).$$
A ring oscillator at $100 \pm 4$ MHz with a spec from 94 to 108 MHz has $z$ values -1.5 and 2, so $Y = 0.9772 - 0.0668 = 91.0$ percent. The two tails are budgeted separately, and the asymmetric spec makes the low-frequency tail dominate. For a non-Gaussian or multi-dimensional performance, the same definition is an integral over the acceptance region, estimated by Monte Carlo.

## Combining independent requirements

A die is good only if every requirement holds. If they are independent, probabilities multiply: a defect-limited yield of 0.85 and three parametric yields of 0.98, 0.97 and 0.99 combine to $0.85 \times 0.98 \times 0.97 \times 0.99 = 0.80$. Independence is an assumption that must be checked. Parametric specs driven by the same global process shift are positively correlated, and the joint yield is then higher than the product; the product is a pessimistic estimate in that case. The decomposition is still useful, because it shows where the loss is: here the defect-limited term costs 15 points and all parametric specs together cost about 6.

## Good dies per wafer

The economic output is good dies per wafer. Gross dies per wafer on a wafer of diameter $d$ with die area $A$ is commonly approximated as
$$DPW = \frac{\pi d^2}{4A} - \frac{\pi d}{\sqrt{2A}},$$
the wafer area divided by die area minus a correction for partial dies at the edge. A 300 mm wafer with 100 mm$^2$ dies gives $706.9 - 66.6 = 640$ gross dies; at 80 percent yield, 512 good ones. Larger dies lose twice: fewer fit, and each is more likely to catch a defect. Yield models therefore enter architecture decisions such as chiplet partitioning, where splitting one large die into several small ones raises yield per die at the cost of packaging and interconnect.

## Going deeper

Koren and Koren, "Defect tolerance in VLSI circuits: techniques and yield analysis," Proceedings of the IEEE, 1998; Orshansky, Nassif and Boning (2008), the chapter on parametric yield.
