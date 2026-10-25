> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed](https://discrete.openmathbooks.org/dmoi3/sec_intro-statements.html) by Oscar Levin, CC BY-SA 4.0

# Predicates and Quantifiers

##  Section 0.2 Mathematical Statements

######  Investigate!

While walking through a fictional forest, you encounter three trolls guarding a bridge. Each is either a _knight_ , who always tells the truth, or a _knave_ , who always lies. The trolls will not let you pass until you correctly identify each as either a knight or a knave. Each troll makes a single statement:

> Troll 1: If I am a knave, then there are exactly two knights here.
>
> Troll 2: Troll 1 is lying.
>
> Troll 3: Either we are all knaves or at least one of us is a knight.

Which troll is which?

In order to _do_ mathematics, we must be able to _talk_ and _write_ about mathematics. Perhaps your experience with mathematics so far has mostly involved finding answers to problems. As we embark towards more advanced and abstract mathematics, writing will play a more prominent role in the mathematical process.

Communication in mathematics requires more precision than many other subjects, and thus we should take a few pages here to consider the basic building blocks: _mathematical statements_.

###  Subsection Atomic and Molecular Statements

A statement is any declarative sentence which is either true or false. A statement is atomic if it cannot be divided into smaller statements, otherwise it is called molecular.

####  Example 0.2.1.

These are statements (in fact, _atomic_ statements):

  * Telephone numbers in the USA have 10 digits.

  * The moon is made of cheese.

  * 42 is a perfect square.

  * Every even number greater than 2 can be expressed as the sum of two primes.

  * \(\displaystyle 3+7 = 12\)

And these are not statements:

  * Would you like some cake?

  * The sum of two squares.

  * \(1+3+5+7+\cdots+2n+1\text{.}\)

  * Go to your room!

  * \(\displaystyle 3+x = 12\)

The reason the sentence “\(3 + x = 12\)” is not a statement is that it contains a variable. Depending on what \(x\) is, the sentence is either true or false, but right now it is neither. One way to make the _sentence_ into a _statement_ is to specify the value of the variable in some way. This could be done by specifying a specific substitution, for example, “\(3+x = 12\) where \(x = 9\text{,}\)” which is a true statement. Or you could _capture_ the free variable by _quantifying_ over it, as in, “for all values of \(x\text{,}\) \(3+x = 12\text{,}\)” which is false. We will discuss quantifiers in more detail at the end of this section.

You can build more complicated (molecular) statements out of simpler (atomic or molecular) ones using logical connectives. For example, this is a molecular statement:

> Telephone numbers in the USA have 10 digits and 42 is a perfect square.

Note that we can break this down into two smaller statements. The two shorter statements are _connected_ by an “and.” We will consider 5 connectives: “and” (Sam is a man and Chris is a woman), “or” (Sam is a man or Chris is a woman), “if…, then…” (if Sam is a man, then Chris is a woman), “if and only if” (Sam is a man if and only if Chris is a woman), and “not” (Sam is not a man). The first four are called binary connectives (because they connect two statements) while “not” is an example of a unary connective (since it applies to a single statement).

These molecular statements are of course still statements, so they must be either true or false. The absolutely key observation here is that which truth value the molecular statement achieves is completely determined by the type of connective and the truth values of the parts. We do not need to know what the parts actually say, only whether those parts are true or false. So to analyze logical connectives, it is enough to consider propositional variables (sometimes called _sentential_ variables), usually capital letters in the middle of the alphabet: \(P, Q, R, S, \ldots\text{.}\) We think of these as standing in for (usually atomic) statements, but there are only two _values_ the variables can achieve: true or false. 1

In computer programming, we should call such variables Boolean variables.

We also have symbols for the logical connectives: \(\wedge\text{,}\) \(\vee\text{,}\) \(\imp\text{,}\) \(\iff\text{,}\) \(\neg\text{.}\)

#### Logical Connectives.

  * \(P \wedge Q\) is read “\(P\) and \(Q\text{,}\)” and called a conjunction.

  * \(P \vee Q\) is read “\(P\) or \(Q\text{,}\)” and called a disjunction.

  * \(P \imp Q\) is read “if \(P\) then \(Q\text{,}\)” and called an implication or conditional.

  * \(P \iff Q\) is read “\(P\) if and only if \(Q\text{,}\)” and called a biconditional.

  * \(\neg P\) is read “not \(P\text{,}\)” and called a negation.

The truth value of a statement is determined by the truth value(s) of its part(s), depending on the connectives:

#### Truth Conditions for Connectives.

  * \(P \wedge Q\) is true when both \(P\) and \(Q\) are true.

  * \(P \vee Q\) is true when \(P\) or \(Q\) or both are true.

  * \(P \imp Q\) is true when \(P\) is false or \(Q\) is true or both.

  * \(P \iff Q\) is true when \(P\) and \(Q\) are both true, or both false.

  * \(\neg P\) is true when \(P\) is false.

Note that for us, _or_ is the inclusive or (and not the sometimes used _exclusive or_) meaning that \(P \vee Q\) is in fact true when both \(P\) and \(Q\) are true. As for the other connectives, “and” behaves as you would expect, as does negation. The biconditional (if and only if) might seem a little strange, but you should think of this as saying the two parts of the statements are _equivalent_ in that they have the same truth value. This leaves only the conditional \(P \imp Q\) which has a slightly different meaning in mathematics than it does in ordinary usage. However, implications are so common and useful in mathematics, that we must develop fluency with their use, and as such, they deserve their own subsection.

###  Subsection Implications

#### Implications.

An implication or conditional is a molecular statement of the form

$$
P \imp Q
$$

where \(P\) and \(Q\) are statements. We say that

  * \(P\) is the hypothesis (or antecedent).

  * \(Q\) is the conclusion (or consequent).

An implication is _true_ provided \(P\) is false or \(Q\) is true (or both), and _false_ otherwise. In particular, the only way for \(P \imp Q\) to be false is for \(P\) to be true _and_ \(Q\) to be false.

Easily the most common type of statement in mathematics is the implication. Even statements that do not at first look like they have this form conceal an implication at their heart. Consider the _Pythagorean Theorem_. Many a college freshman would quote this theorem as “\(a^2 + b^2 = c^2\text{.}\)” This is absolutely not correct. For one thing, that is not a statement since it has three variables in it. Perhaps they imply that this should be true for any values of the variables? So \(1^2 + 5^2 = 2^2\text{???}\) How can we fix this? Well, the equation is true as long as \(a\) and \(b\) are the legs of a right triangle and \(c\) is the hypotenuse. In other words:

> _If_ \(a\) and \(b\) are the legs of a right triangle with hypotenuse \(c\text{,}\) _then_ \(a^2 + b^2 = c^2\text{.}\)

This is a reasonable way to think about implications: our claim is that the conclusion (“then” part) is true, but on the assumption that the hypothesis (“if” part) is true. We make no claim about the conclusion in situations when the hypothesis is false. 2

However, note that in the case of the Pythagorean Theorem, it is also the case that _if_ \(a^2 + b^2 = c^2\text{,}\) _then_ \(a\) and \(b\) are the legs of a right triangle with hypotenuse \(c\text{.}\) So we could have also expressed this theorem as a biconditional: “\(a\) and \(b\) are the legs of a right triangle with hypotenuse \(c\) _if and only if_ \(a^2 + b^2 = c^2\text{.}\)”

Still, it is important to remember that an implication is a statement, and therefore is either true or false. The truth value of the implication is determined by the truth values of its two parts. To agree with the usage above, we say that an implication is true either when the hypothesis is false, or when the conclusion is true. This leaves only one way for an implication to be false: when the hypothesis is true and the conclusion is false.

####  Example 0.2.2.

Consider the statement:

> If Bob gets a 90 on the final, then Bob will pass the class.

This is definitely an implication: \(P\) is the statement “Bob gets a 90 on the final,” and \(Q\) is the statement “Bob will pass the class.”

Suppose I made that statement to Bob. In what circumstances would it be fair to call me a liar? What if Bob really did get a 90 on the final, and he did pass the class? Then I have not lied; my statement is true. However, if Bob did get a 90 on the final and did not pass the class, then I lied, making the statement false. The tricky case is this: what if Bob did not get a 90 on the final? Maybe he passes the class, maybe he doesn’t. Did I lie in either case? I think not. In these last two cases, \(P\) was false, and the statement \(P \imp Q\) was true. In the first case, \(Q\) was true, and so was \(P \imp Q\text{.}\) So \(P \imp Q\) is true when either \(P\) is false or \(Q\) is true.

Just to be clear, although we sometimes read \(P \imp Q\) as “\(P\) _implies_ \(Q\)”, we are not insisting that there is some causal relationship between the statements \(P\) and \(Q\text{.}\) In particular, if you claim that \(P \imp Q\) is _false_ , you are not saying that \(P\) does not imply \(Q\text{,}\) but rather that \(P\) is true and \(Q\) is false.

####  Example 0.2.3.

Decide which of the following statements are true and which are false. Briefly explain.

  1. If \(1=1\text{,}\) then most horses have 4 legs.

  2. If \(0=1\text{,}\) then \(1=1\text{.}\)

  3. If 8 is a prime number, then the 7624th digit of \(\pi\) is an 8.

  4. If the 7624th digit of \(\pi\) is an 8, then \(2+2 = 4\text{.}\)

It is important to understand the conditions under which an implication is true not only to decide whether a mathematical statement is true, but in order to _prove_ that it is. Proofs might seem scary (especially if you have had a bad high school geometry experience) but all we are really doing is explaining (very carefully) why a statement is true. If you understand the truth conditions for an implication, you already have the outline for a proof.

#### Direct Proofs of Implications.

To prove an implication \(P \imp Q\text{,}\) it is enough to assume \(P\text{,}\) and from it, deduce \(Q\text{.}\)

Perhaps a better way to say this is that to prove a statement of the form \(P \imp Q\) directly, you must explain why \(Q\) is true, but you _get to_ assume \(P\) is true first. After all, you only care about whether \(Q\) is true in the case that \(P\) is as well.

There are other techniques to prove statements (implications and others) that we will encounter throughout our studies, and new proof techniques are discovered all the time. Direct proof is the easiest and most elegant style of proof and has the advantage that such a proof often does a great job of explaining _why_ the statement is true.

####  Example 0.2.4.

Prove: If two numbers \(a\) and \(b\) are even, then their sum \(a+b\) is even.

This sort of argument shows up outside of math as well. If you ever found yourself starting an argument with “hypothetically, let’s assume …,” then you have attempted a direct proof of your desired conclusion.

An implication is a way of expressing a relationship between two statements. It is often interesting to ask whether there are other relationships between the statements. Here we introduce some common language to address this question.

#### Converse and Contrapositive.

  * The converse of an implication \(P \imp Q\) is the implication \(Q \imp P\text{.}\) The converse is NOT logically equivalent to the original implication. That is, whether the converse of an implication is true is independent of the truth of the implication.

  * The contrapositive of an implication \(P \imp Q\) is the statement \(\neg Q \imp \neg P\text{.}\) An implication and its contrapositive are logically equivalent (they are either both true or both false).

Mathematics is overflowing with examples of true implications which have a false converse. If a number greater than 2 is prime, then that number is odd. However, just because a number is odd does not mean it is prime. If a shape is a square, then it is a rectangle. But it is false that if a shape is a rectangle, then it is a square.

However, sometimes the converse of a true statement is also true. For example, the Pythagorean theorem has a true converse: if \(a^2 + b^2 = c^2\text{,}\) then the triangle with sides \(a\text{,}\) \(b\text{,}\) and \(c\) is a _right_ triangle. Whenever you encounter an implication in mathematics, it is always reasonable to ask whether the converse is true.

The contrapositive, on the other hand, always has the same truth value as its original implication. This can be very helpful in deciding whether an implication is true: often it is easier to analyze the contrapositive.

####  Example 0.2.5.

True or false: If you draw any nine playing cards from a regular deck, then you will have at least three cards all of the same suit. Is the converse true?

Understanding converses and contrapositives can help understand implications and their truth values:

####  Example 0.2.6.

Suppose I tell Sue that if she gets a 93% on her final, then she will get an A in the class. Assuming that what I said is true, what can you conclude in the following cases:

  1. Sue gets a 93% on her final.

  2. Sue gets an A in the class.

  3. Sue does not get a 93% on her final.

  4. Sue does not get an A in the class.

As we said above, an implication is not logically equivalent to its converse, but it is possible that both the implication and its converse are true. In this case, when both \(P \imp Q\) and \(Q \imp P\) are true, we say that \(P\) and \(Q\) are equivalent and write \(P \iff Q\text{.}\) This is the biconditional we mentioned earlier.

#### If and only if.

> \(P \iff Q\) is logically equivalent to \((P \imp Q) \wedge (Q \imp P)\text{.}\)

Example: Given an integer \(n\text{,}\) it is true that \(n\) is even if and only if \(n^2\) is even. That is, if \(n\) is even, then \(n^2\) is even, as well as the converse: if \(n^2\) is even, then \(n\) is even.

You can think of “if and only if” statements as having two parts: an implication and its converse. We might say one is the “if” part, and the other is the “only if” part. We also sometimes say that “if and only if” statements have two directions: a forward direction \((P \imp Q)\) and a backwards direction (\(P \leftarrow Q\text{,}\) which is really just sloppy notation for \(Q \imp P\)).

Let’s think a little about which part is which. Is \(P \imp Q\) the “if” part or the “only if” part? Consider an example.

####  Example 0.2.7.

Suppose it is true that I sing if and only if I’m in the shower. We know this means both that if I sing, then I’m in the shower, and also the converse, that if I’m in the shower, then I sing. Let \(P\) be the statement, “I sing,” and \(Q\) be, “I’m in the shower.” So \(P \imp Q\) is the statement “if I sing, then I’m in the shower.” Which part of the if and only if statement is this?

What we are really asking for is the meaning of “I sing _if_ I’m in the shower” and “I sing _only if_ I’m in the shower.” When is the first one (the “if” part) _false_? When I am in the shower but not singing. That is the same condition on being false as the statement “if I’m in the shower, then I sing.” So the “if” part is \(Q \imp P\text{.}\) On the other hand, to say, “I sing only if I’m in the shower” is equivalent to saying “if I sing, then I’m in the shower,” so the “only if” part is \(P \imp Q\text{.}\)

It is not terribly important to know which part is the “if” or “only if” part, but this does illustrate something very, very important: _there are many ways to state an implication!_

####  Example 0.2.8.

Rephrase the implication, “if I dream, then I am asleep” in as many different ways as possible. Then do the same for the converse.

Hopefully you agree with the above example. We include the “necessary and sufficient” versions because those are common when discussing mathematics. In fact, let’s agree once and for all what they mean.

#### Necessary and Sufficient.

  * “\(P\) is necessary for \(Q\)” means \(Q \imp P\text{.}\)

  * “\(P\) is sufficient for \(Q\)” means \(P \imp Q\text{.}\)

  * If \(P\) is necessary and sufficient for \(Q\text{,}\) then \(P \iff Q\text{.}\)

To be honest, I have trouble with these if I’m not very careful. I find it helps to keep a standard example for reference.

####  Example 0.2.9.

Recall from calculus, if a function is differentiable at a point \(c\text{,}\) then it is continuous at \(c\text{,}\) but that the converse of this statement is not true (for example, \(f(x) = |x|\) at the point 0). Restate this fact using “necessary and sufficient” language.

Thinking about the necessity and sufficiency of conditions can also help when writing proofs and justifying conclusions. If you want to establish some mathematical fact, it is helpful to think what other facts would _be enough_ (be sufficient) to prove your fact. If you have an assumption, think about what must also be necessary if that hypothesis is true.

###  Subsection Predicates and Quantifiers

######  Investigate!

Consider the statements below. Decide whether any are equivalent to each other, or whether any imply any others.

  1. You can fool some people all of the time.

  2. You can fool everyone some of the time.

  3. You can always fool some people.

  4. Sometimes you can fool everyone.

It would be nice to use variables in our mathematical sentences. For example, suppose we wanted to claim that if \(n\) is prime, then \(n+7\) is not prime. This looks like an implication. I would like to write something like

$$
P(n) \imp \neg P(n+7)
$$

where \(P(n)\) means “\(n\) is prime.” But this is not quite right. For one thing, because this sentence has a free variable (that is, a variable that we have not specified anything about), it is not a statement. A sentence that contains variables is called a predicate.

Now, if we plug in a specific value for \(n\text{,}\) we do get a statement. In fact, it turns out that no matter what value we plug in for \(n\text{,}\) we get a true implication in this case. What we really want to say is that _for all_ values of \(n\text{,}\) if \(n\) is prime, then \(n+7\) is not. We need to _quantify_ the variable.

Although there are many types of _quantifiers_ in English (e.g., many, few, most, etc.) in mathematics we, for the most part, stick to two: existential and universal.

#### Universal and Existential Quantifiers.

The existential quantifier is \(\exists\) and is read “there exists” or “there is.” For example,

$$
\exists x (x \lt 0)
$$

asserts that there is a number less than 0.

The universal quantifier is \(\forall\) and is read “for all” or “every.” For example,

$$
\forall x (x \ge 0)
$$

asserts that every number is greater than or equal to 0.

As with all mathematical statements, we would like to decide whether quantified statements are true or false. Consider the statement

$$
\forall x \exists y (y \lt x)\text{.}
$$

You would read this, “for every \(x\) there is some \(y\) such that \(y\) is less than \(x\text{.}\)” Is this true? The answer depends on what our _domain of discourse_ is: when we say “for all” \(x\text{,}\) do we mean all positive integers or all real numbers or all elements of some other set? Usually this information is implied. In discrete mathematics, we almost always quantify over the _natural numbers_ , 0, 1, 2, …, so let’s take that for our domain of discourse here.

For the statement to be true, we need it to be the case that no matter what natural number we select, there is always some natural number that is strictly smaller. Perhaps we could let \(y\) be \(x-1\text{?}\) But here is the problem: what if \(x = 0\text{?}\) Then \(y = -1\) and that is _not a number!_ (in our domain of discourse). Thus we see that the statement is false because there is a number which is less than or equal to all other numbers. In symbols,

$$
\exists x \forall y (y \ge x)\text{.}
$$

To show that the original statement is false, we proved that the _negation_ was true. Notice how the negation and original statement compare. This is typical.

#### Quantifiers and Negation.

> \(\neg \forall x P(x)\) is equivalent to \(\exists x \neg P(x)\text{.}\)
>
> \(\neg \exists x P(x)\) is equivalent to \(\forall x \neg P(x) \text{.}\)

Essentially, we can pass the negation symbol over a quantifier, but that causes the quantifier to switch type. This should not be surprising: if not everything has a property, then something doesn’t have that property. And if there is not something with a property, then everything doesn’t have that property.

#### Implicit Quantifiers.

It is always a good idea to be precise in mathematics. Sometimes though, we can relax a little bit, as long as we all agree on a convention. An example of such a convention is to assume that sentences containing predicates with free variables are intended as statements, where the variables are universally quantified.

For example, do you believe that if a shape is a square, then it is a rectangle? But how can that be true if it is not a statement? To be a little more precise, we have two predicates: \(S(x)\) standing for “\(x\) is a square” and \(R(x)\) standing for “\(x\) is a rectangle”. The _sentence_ we are looking at is,

$$
S(x) \imp R(x)\text{.}
$$

This is neither true nor false, as it is not a statement. But come on! We all know that we meant to consider the statement,

$$
\forall x (S(x) \imp R(x))\text{,}
$$

and this is what our convention tells us to consider.

Similarly, we will often be a bit sloppy about the distinction between a predicate and a statement. For example, we might write, _let \(P(n)\) be the _statement_ , “\(n\) is prime,”_ which is technically incorrect. It is implicit that we mean that we are defining \(P(n)\) to be a predicate, which for each \(n\) becomes the statement, \(n\) is prime.
