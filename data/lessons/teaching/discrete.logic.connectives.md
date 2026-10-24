> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed](https://discrete.openmathbooks.org/dmoi3/sec_propositional.html) by Oscar Levin, CC BY-SA 4.0

# Propositional Logic

##  Section 3.1 Propositional Logic

######  Investigate!

You stumble upon two trolls playing Stratego®. They tell you:

> Troll 1: If we are cousins, then we are both knaves.
>
> Troll 2: We are cousins or we are both knaves.

Could both trolls be knights? Recall that all trolls are either always-truth-telling knights or always-lying knaves.

A proposition is simply a statement. Propositional logic studies the ways statements can interact with each other. It is important to remember that propositional logic does not really care about the content of the statements. For example, in terms of propositional logic, the claims, “if the moon is made of cheese then basketballs are round,” and “if spiders have eight legs then Sam walks with a limp” are exactly the same. They are both implications: statements of the form, \(P \imp Q\text{.}\)

###  Subsection Truth Tables

Here’s a question about playing Monopoly:

> If you get more doubles than any other player then you will lose, or if you lose then you must have bought the most properties.

True or false? We will answer this question, and won’t need to know anything about Monopoly. Instead we will look at the logical _form_ of the statement.

We need to decide when the statement \((P \imp Q) \vee (Q \imp R)\) is true. Using the definitions of the connectives in Section 0.2, we see that for this to be true, either \(P \imp Q\) must be true or \(Q \imp R\) must be true (or both). Those are true if either \(P\) is false or \(Q\) is true (in the first case) and \(Q\) is false or \(R\) is true (in the second case). So—yeah, it gets kind of messy. Luckily, we can make a chart to keep track of all the possibilities. Enter truth tables. The idea is this: on each row, we list a possible combination of T’s and F’s (for true and false) for each of the sentential variables, and then mark down whether the statement in question is true or false in that case. We do this for every possible combination of T’s and F’s. Then we can clearly see in which cases the statement is true or false. For complicated statements, we will first fill in values for each part of the statement, as a way of breaking up our task into smaller, more manageable pieces.

Since the truth value of a statement is completely determined by the truth values of its parts and how they are connected, all you really need to know is the truth tables for each of the logical connectives. Here they are:

\(P\) | \(Q\) | \(P\wedge Q\)
---|---|---
T | T | T
T | F | F
F | T | F
F | F | F

\(P\) | \(Q\) | \(P\vee Q\)
---|---|---
T | T | T
T | F | T
F | T | T
F | F | F

\(P\) | \(Q\) | \(P\imp Q\)
---|---|---
T | T | T
T | F | F
F | T | T
F | F | T

\(P\) | \(Q\) | \(P\iff Q\)
---|---|---
T | T | T
T | F | F
F | T | F
F | F | T

The truth table for negation looks like this:

\(P\) | \(\neg P\)
---|---
T | F
F | T

None of these truth tables should come as a surprise; they are all just restating the definitions of the connectives. Let’s try another one.

####  Example 3.1.1.

Make a truth table for the statement \(\neg P \vee Q\text{.}\)

Now let’s answer our question about monopoly:

####  Example 3.1.2.

Analyze the statement, “if you get more doubles than any other player you will lose, or that if you lose you must have bought the most properties,” using truth tables.

The statement about monopoly is an example of a tautology, a statement which is true on the basis of its logical form alone. Tautologies are always true but they don’t tell us much about the world. No knowledge about monopoly was required to determine that the statement was true. In fact, it is equally true that “If the moon is made of cheese, then Elvis is still alive, or if Elvis is still alive, then unicorns have 5 legs.”

###  Subsection Logical Equivalence

You might have noticed in Example 3.1.1 that the final column in the truth table for \(\neg P \vee Q\) is identical to the final column in the truth table for \(P \imp Q\text{:}\)

\(P\) | \(Q\) | \(P \imp Q\) | \(\neg P \vee Q\)
---|---|---|---
T | T | T | T
T | F | F | F
F | T | T | T
F | F | T | T

This says that no matter what \(P\) and \(Q\) are, the statements \(\neg P \vee Q\) and \(P \imp Q\) either both true or both false. We therefore say these statements are logically equivalent.

#### Logical Equivalence.

Two (molecular) statements \(P\) and \(Q\) are logically equivalent provided \(P\) is true precisely when \(Q\) is true. That is, \(P\) and \(Q\) have the same truth value under any assignment of truth values to their atomic parts.

To verify that two statements are logically equivalent, you can make a truth table for each and check whether the columns for the two statements are identical.

Recognizing two statements as logically equivalent can be very helpful. Rephrasing a mathematical statement can often lend insight into what it is saying, or how to prove or refute it. By using truth tables we can systematically verify that two statements are indeed logically equivalent.

####  Example 3.1.3.

Are the statements, “it will not rain or snow” and “it will not rain and it will not snow” logically equivalent?

Notice that this example gives us a way to “distribute” a negation over a disjunction (an “or”). We have a similar rule for distributing over conjunctions (“and”s):

#### De Morgan’s Laws.

$$
\neg(P \wedge Q) \text{ is logically equivalent to } \neg P \vee \neg Q\text{.}
$$ $$
\neg(P \vee Q) \text{ is logically equivalent to } \neg P \wedge \neg Q\text{.}
$$

This suggests there might be a sort of “algebra” you could apply to statements (okay, there is: it is called _Boolean algebra_) to transform one statement into another. We can start collecting useful examples of logical equivalence, and apply them in succession to a statement, instead of writing out a complicated truth table.

De Morgan’s laws do not do not directly help us with implications, but as we saw above, every implication can be written as a disjunction:

#### Implications are Disjunctions.

$$
P \imp Q \text{ is logically equivalent to } \neg P \vee Q\text{.}
$$

Example: “If a number is a multiple of 4, then it is even” is equivalent to, “a number is not a multiple of 4 or (else) it is even.”

With this and De Morgan’s laws, you can take any statement and _simplify_ it to the point where negations are only being applied to atomic propositions. Well, actually not, because you could get multiple negations stacked up. But this can be easily dealt with:

#### Double Negation.

$$
\neg \neg P \text{ is logically equivalent to } P\text{.}
$$

Example: “It is not the case that \(c\) is not odd” means “\(c\) is odd.”

Let’s see how we can apply the equivalences we have encountered so far.

####  Example 3.1.4.

Prove that the statements \(\neg(P \imp Q)\) and \(P\wedge \neg Q\) are logically equivalent without using truth tables.

Notice that the above example illustrates that the negation of an implication is NOT an implication: it is a conjunction! We saw this before, in Section 0.2, but it is so important and useful, it warants a second blue box here:

#### Negation of an Implication.

The negation of an implication is a conjunction:

$$
\neg(P \imp Q) \text{ is logically equivalent to } P \wedge \neg Q\text{.}
$$

That is, the only way for an implication to be false is for the hypothesis to be true _AND_ the conclusion to be false.

To verify that two statements are logically equivalent, you can use truth tables or a sequence of logically equivalent replacements. The truth table method, although cumbersome, has the advantage that it can verify that two statements are NOT logically equivalent.

####  Example 3.1.5.

Are the statements \((P \vee Q) \imp R\) and \((P \imp R) \vee (Q \imp R)\) logically equivalent?

###  Subsection Deductions

######  Investigate!

Holmes owns two suits: one black and one tweed. He always wears either a tweed suit or sandals. Whenever he wears his tweed suit and a purple shirt, he chooses to not wear a tie. He never wears the tweed suit unless he is also wearing either a purple shirt or sandals. Whenever he wears sandals, he also wears a purple shirt. Yesterday, Holmes wore a bow tie. What else did he wear?

Earlier we claimed that the following was a valid argument:

> If Edith eats her vegetables, then she can have a cookie. Edith ate her vegetables. Therefore Edith gets a cookie.

How do we know this is valid? Let’s look at the form of the statements. Let \(P\) denote “Edith eats her vegetables” and \(Q\) denote “Edith can have a cookie.” The logical form of the argument is then:

| \(P \imp Q\)
---|---
| \(P\)
\(\therefore\) | \(Q\)

This is an example of a deduction rule, an argument form which is always valid. This one is a particularly famous rule called _modus ponens_. Are you convinced that it is a valid deduction rule? If not, consider the following truth table:

\(P\) | \(Q\) | \(P\imp Q\)
---|---|---
T | T | T
T | F | F
F | T | T
F | F | T

This is just the truth table for \(P \imp Q\text{,}\) but what matters here is that all the lines in the deduction rule have their own column in the truth table. Remember that an argument is valid provided the conclusion must be true given that the premises are true. The premises in this case are \(P \imp Q\) and \(P\text{.}\) Which _rows_ of the truth table correspond to both of these being true? \(P\) is true in the first two rows, and of those, only the first row has \(P \imp Q\) true as well. And lo-and-behold, in this one case, \(Q\) is also true. So if \(P\imp Q\) and \(P\) are both true, we see that \(Q\) must be true as well.

Here are a few more examples.

####  Example 3.1.6.

Show that

| \(P \imp Q\)
---|---
| \(\neg P \imp Q\)
\(\therefore\) | \(Q\)

is a valid deduction rule.

####  Example 3.1.7.

Decide whether

| \(P \imp R\)
---|---
| \(Q \imp R\)
| \(R\)
\(\therefore\) | \(P \vee Q\)

is a valid deduction rule.

###  Subsection Beyond Propositions

As we saw in Section 0.2, not every statement can be analyzed using logical connectives alone. For example, we might want to work with the statement:

> All primes greater than 2 are odd.

To write this statement symbolically, we must use quantifiers. We can translate as follows:

$$
\forall x ((P(x) \wedge x \gt 2) \imp O(x))\text{.}
$$

In this case, we are using \(P(x)\) to denote “\(x\) is prime” and \(O(x)\) to denote “\(x\) is odd.” These are not propositions, since their truth value depends on the input \(x\text{.}\) Better to think of \(P\) and \(O\) as denoting _properties_ of their input. The technical term for these is predicates and when we study them in logic, we need to use predicate logic.

It is important to stress that predicate logic _extends_ propositional logic (much in the way quantum mechanics extends classical mechanics). You will notice that our statement above still used the (propositional) logical connectives. Everything that we learned about logical equivalence and deductions still applies. However, predicate logic allows us to analyze statements at a higher resolution, digging down into the individual propositions \(P\text{,}\) \(Q\text{,}\) etc.

A full treatment of predicate logic is beyond the scope of this text. One reason is that there is no systematic procedure for deciding whether two statements in predicate logic are logically equivalent (i.e., there is no analogue to truth tables here). Rather, we end with a two examples of logical equivalence and deduction, to pique your interest.

####  Example 3.1.8.

Suppose we claim that there is no smallest number. We can translate this into symbols as

$$
\neg \exists x \forall y (x \le y)
$$

(literally, “it is not true that there is a number \(x\) such that for all numbers \(y\text{,}\) \(x\) is less than or equal to \(y\)”).

However, we know how negation interacts with quantifiers: we can pass a negation over a quantifier by switching the quantifier type (between universal and existential). So the statement above should be _logically equivalent_ to

$$
\forall x \exists y (y \lt x)\text{.}
$$

Notice that \(y \lt x\) is the negation of \(x \le y\text{.}\) This literally says, “for every number \(x\) there is a number \(y\) which is smaller than \(x\text{.}\)” We see that this is another way to make our original claim.

####  Example 3.1.9.

Can you switch the order of quantifiers? For example, consider the two statements:

$$
\forall x \exists y P(x,y) \qquad \mathrm{ and } \qquad \exists y \forall x P(x,y)\text{.}
$$

Are these logically equivalent?
