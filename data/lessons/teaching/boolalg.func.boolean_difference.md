# Boolean difference and sensitivity

The Boolean difference of $f$ with respect to $x$ is the function that is 1 exactly where changing $x$ changes $f$. It is the Boolean analogue of a partial derivative. Test generation uses it to decide when a fault on $x$ is observable, synthesis uses it to compute observability don't cares, and timing and power analysis use it for sensitization and switching activity. It is computed from the two cofactors with one XOR.

## The XOR of the two cofactors

Define $\dfrac{\partial f}{\partial x} = f_x \oplus f_{\overline{x}}$. It is a function of the variables other than $x$, and it is 1 at an assignment of those variables exactly when the two cofactors disagree there, so that flipping $x$ flips $f$.

For $f(a, b, c) = ab + \overline{c}$: $f_a = b + \overline{c}$ and $f_{\overline{a}} = \overline{c}$. They differ only when $b = 1$ and $c = 1$ (1 against 0), so $\partial f / \partial a = bc$, one ON row of four. With respect to $c$: $f_c = ab$ and $f_{\overline{c}} = 1$, so $\partial f / \partial c = \overline{ab}$, three rows. A variable is in the support of $f$ exactly when its Boolean difference is not identically 0.

## Inputs where a variable matters

Each ON row of $\partial f / \partial x$ is an assignment of the other variables, and it pairs with both values of $x$ to give two full input vectors at which flipping $x$ changes $f$. So the number of sensitive input vectors is $2 \cdot |\partial f / \partial x|$. For $f = ab + \overline{c}$ and $x = a$, the sensitive vectors are 011 and 111.

Divided by $2^n$, that count is the influence of $x$ on $f$, the probability that $x$ is pivotal under uniform inputs. It feeds switching-activity estimates in power analysis and is a quick filter when choosing which variables to split on first.

## Sensitization in test generation

Model input $x$ of a block as stuck at 0. A test vector must set $x = 1$, so the good and faulty circuits differ at the fault site (activation), and must make $\partial f / \partial x = 1$, so the difference reaches the output (propagation). The tests are the ON rows of $x \cdot \partial f / \partial x$; for stuck-at 1 they are those of $\overline{x} \cdot \partial f / \partial x$.

For $f = ab + \overline{c}$ and $a$ stuck at 0: $a = 1$ and $bc = 1$, so the only test is 111. The good circuit outputs 1; the faulty one computes $f$ at 011, which is 0. Further constraints, such as a test-mode pin fixed to a value, simply intersect the test set with that literal.

## Going deeper

Sellers, Hsiao and Bearnson, "Analyzing errors with the Boolean difference", IEEE Transactions on Computers (1968). Bushnell and Agrawal, "Essentials of Electronic Testing", chapter 7.
