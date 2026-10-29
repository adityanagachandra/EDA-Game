# Signature aliasing

A MISR compresses millions of response bits into a few dozen, so information is lost. A faulty circuit can produce a signature equal to the good one, a loss called aliasing, and the defective die then passes. The probability is small and easy to estimate, and it sets the minimum signature length a BIST design needs for a quality target.

## The basic estimate

Model the faulty response as any of the possible erroneous streams with equal likelihood. Because the MISR is linear, the faulty signature equals the good one exactly when the error stream's signature is zero. For a long stream, the error signatures spread almost uniformly over the $2^n$ values, so the aliasing probability approaches $2^{-n}$. A 16-bit MISR aliases with probability about $1.5 \times 10^{-5}$, or 15.26 per million; 32 bits bring it to about $2.3 \times 10^{-10}$. For a single-input signature register fed $L$ bits, the exact value under this model is $(2^{L-n} - 1)/(2^L - 1)$, which tends to $2^{-n}$.

## Sizing the signature

To bring aliasing below a target $\epsilon$, choose $n \ge \log_2(1/\epsilon)$ and round up. For $\epsilon = 2 \times 10^{-4}$, $\log_2(5000) = 12.29$, so 13 bits. In practice MISRs are 16 to 64 bits wide, set by the number of chains they compact rather than by aliasing, so aliasing is rarely the limiting factor; it is the X values in responses, which make the signature unpredictable, that dominate BIST design effort.

## Escapes at volume

Expected aliased escapes are the number of defective dies reaching the MISR times $2^{-n}$. With 1.88 million such dies and a 15-bit MISR, about 57 alias, which is significant; at 32 bits the expectation is negligible. The equal-likelihood model is an assumption: real error streams are structured, and some structures alias more often, for example errors that repeat with the polynomial's period. Since the signature is checked only at the end, a MISR also hides which pattern failed, so diagnosis needs intermediate signatures or a fall-back to scan unload.

## Going deeper

Bardell, McAnney and Savir (1987), chapter on signature analysis; Bushnell and Agrawal (2000), chapter 15; Williams, Daehn, Gruetzner and Starke, "Bounds and analysis of aliasing errors in linear feedback shift registers," IEEE Transactions on CAD, 1988.
