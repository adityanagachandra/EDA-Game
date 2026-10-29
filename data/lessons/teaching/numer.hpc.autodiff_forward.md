# Forward-mode automatic differentiation

Newton iteration needs Jacobians, sensitivity analysis needs gradients, and analytical placers optimize by gradient. Automatic differentiation computes exact derivatives of a program by applying the chain rule to each operation. Forward mode, the simplest kind, carries a derivative alongside every value and is how many device-model compilers generate the derivatives SPICE needs.

## Dual numbers carry a value and a derivative

A dual number $a + b\varepsilon$ with $\varepsilon^2 = 0$ carries a value $a$ and a derivative $b$. Arithmetic follows: $(a + b\varepsilon) + (c + d\varepsilon) = (a + c) + (b + d)\varepsilon$ and
$$(a + b\varepsilon)(c + d\varepsilon) = ac + (ad + bc)\varepsilon,$$
the product rule. For any smooth $f$, $f(a + b\varepsilon) = f(a) + f'(a)b\,\varepsilon$.

Evaluate $f(x) = (x^2 + 3x)(x - 1)$ at $x = 2 + 1\varepsilon$: $x^2 + 3x = 10 + 7\varepsilon$, $x - 1 = 1 + \varepsilon$, product $10 + (10 + 7)\varepsilon = 10 + 17\varepsilon$. So $f(2) = 10$ and $f'(2) = 17$, exactly. Finite differences would trade truncation error against cancellation; AD has neither, only ordinary roundoff.

## Propagating tangents through a computation

For several inputs, forward mode seeds a tangent direction $\dot x$ and propagates $\dot v = \sum_j(\partial v/\partial u_j)\dot u_j$ through every intermediate $v$. The output tangent is the directional derivative $J\dot x$.

For $f(x, y) = x^2y + y\sin x$ at $(1, 2)$ with seed $(1, 0)$: $\dot{(x^2)} = 2x = 2$, $\dot{(x^2y)} = 2 \cdot 2 = 4$, $\dot{(\sin x)} = \cos 1 = 0.5403$, $\dot{(y\sin x)} = 2 \cdot 0.5403 = 1.0806$. Total $\partial f/\partial x = 5.0806$. Seeding $(0, 1)$ instead gives $\partial f/\partial y = x^2 + \sin x = 1.8415$.

Compact-model compilers such as OpenVAF apply this idea to Verilog-A source, generating the current derivatives that fill the MNA Jacobian, so model writers never hand-code derivatives, a major source of Newton failures.

## Cost: one pass per input direction

One forward pass costs a small constant times one function evaluation and gives one column, or one directional derivative, of the Jacobian. A full Jacobian of $f: \mathbb{R}^n \to \mathbb{R}^m$ needs $n$ passes; vector mode carries $w$ tangents per pass, so $\lceil n/w\rceil$ passes.

That is cheap when inputs are few: a MOSFET model has a handful of terminal voltages, so forward mode gives its full Jacobian in about four passes regardless of how many outputs it has. It is expensive when inputs are many and outputs few, as in placement, where a scalar wirelength-plus-density objective depends on millions of coordinates; there reverse mode wins. Sparsity also helps: if Jacobian columns that never share a row are grouped (graph colouring), one pass computes several columns at once.

## Going deeper

Griewank and Walther, "Evaluating Derivatives: Principles and Techniques of Algorithmic Differentiation" (2nd ed., SIAM, 2008), chapters 3 and 8; Baydin et al., "Automatic differentiation in machine learning: a survey", JMLR (2018).
