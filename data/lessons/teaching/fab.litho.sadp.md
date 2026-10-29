# Self-aligned double and quadruple patterning

Self-aligned patterning halves the pitch without a second critical exposure. Mandrel lines are printed at a relaxed pitch, a conformal spacer film is deposited and etched back, and the mandrels are removed, leaving two spacer lines per mandrel. Repeat once more and the pitch quarters. The line width is set by film thickness rather than lithography, which gives excellent uniformity but imposes a gridded, unidirectional layout with separate cut masks.

## From mandrel to spacer lines

Print mandrels at pitch $P$. Deposit a conformal film of thickness $s$ by ALD or CVD and etch it anisotropically so only sidewall spacers remain. Remove the mandrels. Each mandrel leaves two spacers, so the spacer pitch is $P/2$: an 80 nm mandrel pitch gives 40 nm. SAQP repeats the process on the first spacers, giving $P/4$, so a 128 nm mandrel pitch becomes 32 nm. In the spacer-is-line flavour the spacers become the final lines, with width $s$. Because a deposited film is uniform to a fraction of a nanometre, line width varies far less than in litho-etch patterning.

## Alternating spaces

After mandrel removal, two kinds of gap alternate. Where a mandrel stood, the gap equals the mandrel width $m$. Between spacers of neighbouring mandrels, the gap is $P - m - 2s$. With $P = 96$ nm, $m = 30$ nm and $s = 16$ nm: the first gap is 30 nm and the second is $96 - 30 - 32 = 34$ nm. The two gaps come from different physical causes, mandrel lithography and etch for one, mandrel pitch and spacer thickness for the other, so they vary independently. Extraction and timing must treat the two spaces as two populations.

## Pitch walk and cut masks

Equal spaces require $m = P - m - 2s$, so the ideal mandrel width is $m = P/2 - s$: 32 nm for the example above. Any error in printed mandrel width makes line-to-line pitches alternate, $m + s$ across a mandrel and $P - m - s$ between mandrels, and their difference $|2m + 2s - P|$ is the pitch walk. A 4 nm mandrel error gives 8 nm of pitch walk. SADP also produces only continuous parallel lines, so line ends, gaps and jogs need separate cut or block masks, and those cuts carry their own rules on size, spacing and alignment. That is why SADP metal layers are fully gridded and unidirectional, and why routers for them reason about cut placement.

## Going deeper

Bencher, Chen, Dai and Yu, "22nm half-pitch patterning by CVD spacer self alignment double patterning (SADP)", Proc. SPIE (2008).
