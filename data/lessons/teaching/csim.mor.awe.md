# Asymptotic waveform evaluation

Asymptotic waveform evaluation (AWE) turns a handful of moments into a reduced model with a few poles. It matches the moments with a Pade approximation, a rational function whose expansion agrees with the computed moments, and reads the dominant poles and residues from it. AWE made interconnect delay calculation far more accurate than Elmore at little cost, and its numerical weaknesses motivated the Krylov methods that replaced it.

## A one-pole fit from two moments

The simplest model is one pole: $H(s) \approx k/(1 - s/p)$, whose expansion is $k(1 + s/p + s^2/p^2 + \cdots)$. Matching the first two moments gives $m_0 = k$ and $m_1 = k/p$, so
$$p = \frac{m_0}{m_1}.$$
For an RC tree, $m_1 < 0$, so the pole is negative, as a stable pole must be. With $m_0 = 1$ and $m_1 = -2.5$ ns, $p = -0.4$ 1/ns and the time constant is $-1/p = 2.5$ ns, exactly the Elmore delay. The one-pole AWE model is therefore Elmore's model written as a waveform, with the response $1 - e^{pt}$ to a unit step.

## Delay from the dominant pole

A single-pole step response $1 - e^{-t/\tau}$ crosses 50% at $t_{50} = \tau \ln 2 \approx 0.693\tau$. With $\tau = -m_1/m_0$, the one-pole model predicts $t_{50} = 0.693 \times$ Elmore delay. For $m_1 = -0.3$ ns, $t_{50} = 0.208$ ns. The Elmore value itself upper-bounds the true 50% delay of an RC tree. Near the driver, where resistive shielding makes the response far from single-pole, both estimates can be very pessimistic; at the far end of a distributed line, $\ln 2$ times Elmore gives 0.35 RC against a true value near 0.38 RC. Higher-order AWE models were introduced to remove such errors.

## A two-pole Pade fit

A two-pole model $(a_0 + a_1 s)/(1 + b_1 s + b_2 s^2)$ matches four moments. Multiplying out and matching the coefficients of $s^2$ and $s^3$ gives the linear system
$$\begin{pmatrix} m_0 & m_1 \\ m_1 & m_2 \end{pmatrix} \begin{pmatrix} b_2 \\ b_1 \end{pmatrix} = -\begin{pmatrix} m_2 \\ m_3 \end{pmatrix}.$$
By Cramer's rule $b_1 = (m_1 m_2 - m_0 m_3)/(m_0 m_2 - m_1^2)$. For a response with time constants 2 and 0.5 ns weighted 0.8 and 0.2, the moments are $1, -1.7, 3.25, -6.425$, and $b_1 = 0.9/0.36 = 2.5$ ns, the sum of the time constants; the poles are the roots of the denominator. With higher orders the moment matrix becomes ill-conditioned and the Pade model can return unstable poles, which PRIMA's projection avoids.

## Going deeper

Pillage and Rohrer, "Asymptotic waveform evaluation for timing analysis" (IEEE TCAD, 1990); Feldmann and Freund, "Efficient linear circuit analysis by Pade approximation via the Lanczos process" (IEEE TCAD, 1995).
