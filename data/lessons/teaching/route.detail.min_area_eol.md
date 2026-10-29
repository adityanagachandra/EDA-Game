# Minimum-area and end-of-line rules

Beyond spacing, modern rule decks constrain the shape of every wire end. Minimum-area rules forbid tiny isolated metal, end-of-line rules require extra room in front of a line end, and via enclosure rules require metal around each cut. A detailed router has to generate shapes that satisfy all of them at once, which is why short stubs and line ends are where violations cluster.

## Minimum area

Lithography cannot reliably print very small isolated shapes, so each layer has a minimum area $A_{min}$. A stub of width $w$ needs length at least $A_{min}/w$, rounded up to the manufacturing grid. Example: $A_{min} = 9000$ nm$^2$ and $w = 28$ nm need $9000/28 = 321.4$, so 322 nm; a 209 nm stub must grow by 113 nm.

The router fixes this by extending the wire end in its direction, which may itself collide with neighbours, so the extension is part of the search or of a post-processing patch with its own checks. Via landing pads on a layer that carries nothing else are the classic case: the pad is tiny unless deliberately grown.

## End-of-line spacing

Two line ends facing each other on the same track need a larger gap than the normal side-to-side spacing, because line ends print shorter and rounder than drawn. The end-of-line rule requires facing ends to be at least $s_{eol}$ apart. Checking a track is a sweep: sort its segments, subtract each end from the next start, and flag every gap below $s_{eol}$.

Example: segments [0, 83], [153, 256], [296, 566] with $s_{eol} = 50$ nm have gaps 70 and 40, so one violation. Routers price line ends that stop too close to another, and prefer to extend a wire through a short gap (merging the ends) or to shift where it ends. Full end-of-line rules also look sideways within a window, which makes them context dependent and expensive to evaluate.

## Rules that interact

A via lands on a metal stub that must enclose the cut by $e$ beyond each end in the wire direction, so the stub needs length at least $v + 2e$ for cut length $v$. The same stub must also meet minimum area. So the required length is $\max(v + 2e, \lceil A_{min}/w \rceil)$.

Example: a 32 nm cut with 30 nm enclosure needs 92 nm, while $A_{min} = 7000$ nm$^2$ at width 32 nm needs 219 nm, so the stub must be 219 nm. Rules like these stack, and some are conditional on neighbours (end-of-line within a parallel run, spacing by width and run length). That is why routers evaluate shapes against the full rule deck with a geometric engine rather than one formula, and why rule changes between PDK versions break hard-coded router logic.

## Going deeper

The LEF/DEF reference on MINAREA, SPACING ENDOFLINE and via rules; the TritonRoute paper (IEEE TCAD, 2021).
