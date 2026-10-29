# Hold fixing with delay buffers

Hold violations are fixed by adding delay to the short paths, usually by inserting delay buffers. The difficulty is that the same buffers also slow the path at the slow corner, where it may be close to failing setup. A hold fix is therefore sized at one corner and budgeted at another, and a path with too little setup margin cannot be fixed by buffering alone.

## Buffers needed at the fast corner

Hold usually fails at the fast corner, where the short path is shortest. Each buffer adds its fast-corner delay $b_{min}$. To clear a violation of $h$ ps, insert $n = \lceil h / b_{min} \rceil$ buffers. For a 52 ps violation and buffers of 19 ps at the fast corner, $n = \lceil 2.74 \rceil = 3$. Hold fixers prefer to insert buffers close to the capture flop's data pin when the short path is shared with long paths elsewhere, so that the added delay hits only the short path. Hold is fixed after clock tree synthesis, when real skew is known, so that buffers are not spent against estimated skew that never appears.

## Setup cost at the slow corner

The same buffers add their slow-corner delay $b_{max}$, which is larger than $b_{min}$. After inserting $n$ buffers the setup slack at the slow corner becomes $S - n\,b_{max}$. With 191 ps of setup slack, 2 buffers and $b_{max} = 38$ ps, setup falls to 115 ps. The ratio $b_{max}/b_{min}$ is often well above 1, so each picosecond of hold fixed costs more than a picosecond of setup. Paths that fail hold at the fast corner and are close to failing setup at the slow corner are the hard cases; they appear wherever large clock skew or a short logic path meets a tight cycle.

## The feasibility bound

The setup slack limits how many buffers a path can take: $n_{max} = \lfloor S / b_{max} \rfloor$. With 145 ps of setup slack and 36 ps buffers, at most 4 buffers fit. A hold fix is feasible only if $n \le n_{max}$. When it is not, buffering alone cannot work, and the fix moves elsewhere: reduce the clock skew that causes the hold problem, use a buffer with a better min-to-max ratio, insert the delay on a branch of the path that the setup-critical path does not share, or change the launch or capture clock arrival. This interaction between corners is why hold fixing is run with all scenarios active.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter on timing fixes; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8.
