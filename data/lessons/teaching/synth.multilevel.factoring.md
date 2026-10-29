# Factored forms and literal counts

A factored form is a parenthesized expression of sums and products, such as $(a + b)(c + d) + e$. It is the multi-level analogue of an SOP and the best simple model of a static CMOS complex gate, since each literal becomes a transistor pair. Converting between factored forms and SOPs, and counting literals of each, is the core arithmetic of multi-level area estimation.

## Expanding a factored form

Multiplying a factored form out gives an SOP. Cube counts follow simple rules: a sum adds the cube counts of its terms, and a product multiplies them. For $F = (c(e + bg) + h)d$: the inner $e + bg$ has 2 cubes, times $c$ still 2, plus $h$ makes 3, times $d$ stays 3. So the SOP $cde + bcdg + dh$ has 3 cubes. When every variable appears once (a read-once form), there is no simplification during expansion, so these rules are exact. With repeated variables, expansion can produce duplicate or absorbed cubes and the counts become upper bounds.

## SOP literals from a factored form

Literals multiply too. For a product $AB$ the SOP literal count is $l_A c_B + l_B c_A$, where $c$ is a cube count and $l$ a literal count, because each literal of $A$ appears once with every cube of $B$. For a sum, literals add. For $F = ga(h + b + f + ed)$: $ga$ has 1 cube and 2 literals, $h + b + f + ed$ has 4 cubes and 5 literals, so the SOP has $2 \cdot 4 + 5 \cdot 1 = 13$ literals: $agh + abg + afg + adeg$. The factored form has 6. The gap grows quickly with nesting, which is why SOP literal count overstates the area of multi-level logic.

## Finding a good factored form

The literal count of a factored form is the number of leaves. Finding the minimum-literal factored form is hard, so tools use heuristics. **Quick factor** divides $F$ by a level-0 kernel $K$, giving $F = KQ + R$, and recursively factors $K$, $Q$ and $R$. **Good factor** chooses the kernel that saves the most literals. For $bceg + bceh + ad$ (10 literals), dividing by the kernel $g + h$ gives $bce(g + h) + ad$ with 7 literals, a saving of 3. Read-once functions are the easy case, with one leaf per variable, which is optimal. Factored-form literal count is SIS's default area cost and the basis of refactoring in ABC.

## Going deeper

Brayton, "Factoring logic functions", IBM Journal of Research and Development (1987); Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms" (1996), on multi-level factoring.
