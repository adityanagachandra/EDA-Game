# Statistical process control

A fab does not wait for yield loss to notice that its process drifted. It measures critical dimensions, film thicknesses and electrical test structures on every lot and plots them on control charts. Statistical process control is where variation statistics are used in real time, and its numbers, such as the capability index, are what process and design teams negotiate over when a spec is set.

## Control limits for a mean

An $\bar{X}$ chart plots the average of $n$ measurements per sample. When the process is in control with mean $\mu$ and standard deviation $\sigma$, the sample mean has standard error $\sigma/\sqrt{n}$, so the chart draws control limits at
$$\mu \pm 3\frac{\sigma}{\sqrt{n}}.$$
With $\mu = 30$ nm, $\sigma = 1.2$ nm and $n = 9$ sites, the limits are $30 \pm 1.2$ nm. A point outside them is a signal to stop and investigate. Control limits describe what the process does; they are not spec limits, which describe what the product needs. Further rules (runs of points on one side of the centre line, trends) catch smaller drifts at the cost of more false alarms.

## Capability indices

Capability compares the spread to the spec window from $LSL$ to $USL$. $C_p = (USL - LSL)/(6\sigma)$ ignores centring; $C_{pk} = \min(USL - \mu,\ \mu - LSL)/(3\sigma)$ penalizes an off-centre mean. A spec of 27 to 33 nm with $\sigma = 0.8$ nm and $\mu = 30.6$ nm gives $C_p = 6/4.8 = 1.25$ and $C_{pk} = 2.4/2.4 = 1.0$. A $C_{pk}$ of 1 puts the nearest spec limit at $3\sigma$, which for a Gaussian is about 1350 ppm beyond it; $C_{pk} = 1.33$ is $4\sigma$, about 32 ppm; $C_{pk} = 2$ is $6\sigma$. Capability assumes normality and stability, and a fat-tailed process can have a good $C_{pk}$ and a bad defect rate.

## How fast a chart detects a shift

The average run length (ARL) is the expected number of plotted points before a signal. With no shift, the chance a point falls outside 3-sigma limits is 0.0027, so $ARL_0 = 1/0.0027 \approx 370$: one false alarm every 370 samples. If the mean shifts by $\delta\sigma$, the plotted mean shifts by $\delta\sqrt{n}$ of its own standard error, the signal probability becomes $P = Q(3 - \delta\sqrt{n}) + \Phi(-3 - \delta\sqrt{n})$, and $ARL = 1/P$. A one-sigma shift with $n = 4$ gives $Q(1) = 0.159$ and an ARL of about 6.3. Larger samples detect faster but cost metrology time, the same budget trade-off any statistical monitor faces.

## Going deeper

Montgomery, "Introduction to Statistical Quality Control" (8th ed.), the chapters on Shewhart charts and process capability; May and Spanos, "Fundamentals of Semiconductor Manufacturing and Process Control" (2006).
