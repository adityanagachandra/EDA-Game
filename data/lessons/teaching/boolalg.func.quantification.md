# Existential and universal quantification

Quantifying a Boolean variable removes it from a function. The existential quantification $\exists x\, f$ asks whether some value of $x$ makes $f$ true; the universal quantification $\forall x\, f$ asks whether every value does. These are predicate logic restricted to bits, and they are everywhere in hardware reasoning: image computation in symbolic model checking quantifies out current-state and input variables, satisfiability is an existential question over all inputs, validity is a universal one, and synthesis problems with a free parameter take the form "exists a setting such that for all inputs".

## Existential quantification as an OR of cofactors

Since $x$ takes only two values, $\exists x\, f = f_x + f_{\overline{x}}$. The result does not depend on $x$, and it is 1 at an assignment of the other variables exactly when at least one cofactor is 1 there. Existential quantification of every variable in turn gives the constant 1 exactly when $f$ is satisfiable.

For $f(a, b, c) = ab + \overline{c}$: $\exists a\, f = (b + \overline{c}) + \overline{c} = b + \overline{c}$, which is 1 on 3 of the 4 rows over $b, c$. Geometrically, $\exists x$ projects the ON-set onto the other variables, the Boolean version of a shadow.

## Universal quantification as an AND of cofactors

Dually, $\forall x\, f = f_x \cdot f_{\overline{x}}$, which is 1 where both cofactors are 1. By De Morgan, $\forall x\, f = \overline{\exists x\, \overline{f}}$. For $f = ab + \overline{c}$: $\forall a\, f = (b + \overline{c})\,\overline{c} = \overline{c}$, 2 of 4 rows.

Quantifying every variable universally gives 1 exactly when $f$ is a tautology. On BDDs, quantification costs one apply per variable in the worst case and can grow the diagram; the relational product, which conjoins and quantifies in one pass, exists because image computation must avoid building the large intermediate conjunction.

## Nested quantifiers in hardware properties

Quantifiers over different variables nest, innermost first, and their order matters when they are mixed. $\exists u\, \forall v\, f$ asks for one value of $u$ that works for every $v$; $\forall v\, \exists u\, f$ lets $u$ depend on $v$, a weaker demand.

With $f = ab + \overline{c}$: $\exists a\, f = b + \overline{c}$, and then $\forall c$ gives $(b)(b + 1) = b$, so $\forall c\, \exists a\, f = b$, true on 1 of the 2 rows over $b$. In hardware terms: a controller with a free configuration input $u$ meets its specification for every environment input $v$ exactly when $\exists u\, \forall v\, (\text{spec})$ holds. Two alternations make this harder than SAT in general (it is a quantified Boolean formula), which is why such questions are solved by counterexample-guided loops of SAT calls rather than by one call.

## Going deeper

Clarke, Grumberg, Kroening, Peled and Veith, "Model Checking" (2nd ed., 2018), chapter on symbolic model checking. Biere et al., "Handbook of Satisfiability" (2nd ed.), chapter on QBF.
