> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed](https://discrete.openmathbooks.org/dmoi3/sec_intro-sets.html) by Oscar Levin, CC BY-SA 4.0

# Venn Diagrams

##  Section 0.3 Sets

The most fundamental objects we will use in our studies (and really in all of math) are _sets_. Much of what follows might be review, but it is very important that you are fluent in the language of set theory. Most of the notation we use below is standard, although some might be a little different than what you have seen before.

For us, a set will simply be an unordered collection of objects. Two examples: we could consider the set of all actors who have played _The Doctor_ on _Doctor Who_ , or the set of natural numbers between 1 and 10 inclusive. In the first case, Tom Baker is an element (or member) of the set, while Idris Elba, among many others, is not an element of the set. Also, the two examples are of different sets. Two sets are equal exactly if they contain the exact same elements. For example, the set containing all of the vowels in the declaration of independence is precisely the same set as the set of vowels in the word “questionably” (namely, all of them); we do not care about order or repetitions, just whether the element is in the set or not.

###  Subsection Notation

We need some notation to make talking about sets easier. Consider,

$$
A = \{1, 2, 3\}\text{.}
$$

This is read, “\(A\) is the set containing the elements 1, 2 and 3.” We use curly braces “\(\{,~~ \}\)” to enclose elements of a set. Some more notation:

$$
a \in \{a, b, c\}\text{.}
$$

The symbol “\(\in\)” is read “is in” or “is an element of.” Thus the above means that \(a\) is an element of the set containing the letters \(a\text{,}\) \(b\text{,}\) and \(c\text{.}\) Note that this is a true statement. It would also be true to say that \(d\) is not in that set:

$$
d \not\in \{a, b, c\}\text{.}
$$

Be warned: we write “\(x \in A\)” when we wish to express that one of the elements of the set \(A\) is \(x\text{.}\) For example, consider the set,

$$
A = \{1, b, \{x, y, z\}, \emptyset\}\text{.}
$$

This is a strange set, to be sure. It contains four elements: the number 1, the letter b, the set \(\{x,y,z\}\text{,}\) and the empty set \(\emptyset = \{ \}\text{,}\) the set containing no elements. Is \(x\) in \(A\text{?}\) The answer is no. None of the four elements in \(A\) are the letter \(x\text{,}\) so we must conclude that \(x \notin A\text{.}\) Similarly, consider the set \(B = \{1,b\}\text{.}\) Even though the elements of \(B\) are elements of \(A\text{,}\) we cannot say that the _set_ \(B\) is one of the elements of \(A\text{.}\) Therefore \(B \notin A\text{.}\) (Soon we will see that \(B\) is a _subset_ of \(A\text{,}\) but this is different from being an _element_ of \(A\text{.}\))

We have described the sets above by listing their elements. Sometimes this is hard to do, especially when there are a lot of elements in the set (perhaps infinitely many). For instance, if we want \(A\) to be the set of all even natural numbers, would could write,

$$
A = \{0, 2, 4, 6, \ldots\}\text{,}
$$

but this is a little imprecise. A better way would be

$$
A = \{x \in \N \st \exists n\in \N ( x = 2 n)\}\text{.}
$$

Let’s look at this carefully. First, there are some new symbols to digest: “\(\N\)” is the symbol usually used to denote that natural numbers, which we will take to be the set \(\{0, 1, 2, 3, \ldots\}\text{.}\) Next, the colon, “:”, is read _such that_ ; it separates the elements that are in the set from the condition that the elements in the set must satisfy. So putting this all together, we would read the set as, “the set of all \(x\) in the natural numbers, such that there exists some \(n\) in the natural numbers for which \(x\) is twice \(n\text{.}\)” In other words, the set of all natural numbers, that are even. Here is another way to write the same set.

$$
A = \{x \in \N \st x\text{ is even} \}\text{.}
$$

Note: Sometimes mathematicians use \(|\) or \(\backepsilon\) for the “such that” symbol instead of the colon. Also, there is a fairly even split between mathematicians about whether 0 is an element of the natural numbers, so be careful there.

This notation is usually called set builder notation. It tells us how to _build_ a set by telling us precisely the condition elements must meet to gain access (the condition is the logical statement after the “\(\st\)” symbol). Reading and comprehending sets written in this way takes practice. Here are some more examples:

####  Example 0.3.1.

Describe each of the following sets both in words and by listing out enough elements to see the pattern.

  1. \(\{x \st x + 3 \in \N\}\text{.}\)

  2. \(\{x \in \N \st x + 3 \in \N\}\text{.}\)

  3. \(\{x \st x \in \N \vee -x \in \N\}\text{.}\)

  4. \(\{x \st x \in \N \wedge -x \in \N\}\text{.}\)

There is also a subtle variation on set builder notation. While the condition is generally given after the “such that”, sometimes it is hidden in the first part. Here is an example.

####  Example 0.3.2.

List a few elements in the sets below and describe them in words. The set \(\Z\) is the set of integers; positive and negative whole numbers.

  1. \(\displaystyle A = \{x \in \Z \st x^2 \in \N\}\)

  2. \(\displaystyle B = \{x^2 \st x \in \N\}\)

We already have a lot of notation, and there is more yet. Below is a handy chart of symbols. Some of these will be discussed in greater detail as we move forward.

#### Special sets.

\(\emptyset\)

The empty set is the set which contains no elements.

\(\mathcal{U}\)

The universe set is the set of all elements.

\(\N\)

The set of natural numbers. That is, \(\N = \{0, 1, 2, 3\ldots\}\text{.}\)

\(\Z\)

The set of integers. That is, \(\Z = \{\ldots, -2, -1, 0, 1, 2, 3, \ldots\}\text{.}\)

\(\Q\)

The set of rational numbers.

\(\R\)

The set of real numbers.

\(\pow(A)\)

The power set of any set \(A\) is the set of all subsets of \(A\text{.}\)

#### Set Theory Notation.

\(\{, \}\)

We use these braces to enclose the elements of a set. So \(\{1,2,3\}\) is the set containing 1, 2, and 3.

\(\st\)

\(\{x \st x > 2\}\) is the set of all \(x\) such that \(x\) is greater than 2.

\(\in\)

\(2 \in \{1,2,3\}\) asserts that 2 is an element of the set \(\{1,2,3\}\text{.}\)

\(\not\in\)

\(4 \notin \{1,2,3\}\) because 4 is not an element of the set \(\{1,2,3\}\text{.}\)

\(\subseteq\)

\(A \subseteq B\) asserts that \(A\) is a subset of \(B\): every element of \(A\) is also an element of \(B\text{.}\)

\(\subset\)

\(A \subset B\) asserts that \(A\) is a proper subset of \(B\): every element of \(A\) is also an element of \(B\text{,}\) but \(A \ne B\text{.}\)

\(\cap\)

\(A \cap B\) is the intersection of \(A\) and \(B\): the set containing all elements which are elements of both \(A\) and \(B\text{.}\)

\(\cup\)

\(A \cup B\) is the union of \(A\) and \(B\): is the set containing all elements which are elements of \(A\) or \(B\) or both.

\(\times\)

\(A \times B\) is the Cartesian product of \(A\) and \(B\): the set of all ordered pairs \((a,b)\) with \(a \in A\) and \(b \in B\text{.}\)

\(\setminus\)

\(A \setminus B\) is set difference between \(A\) and \(B\): the set containing all elements of \(A\) which are not elements of \(B\text{.}\)

\(\bar{A}\)

The complement of \(A\) is the set of everything which is not an element of \(A\text{.}\)

\(\card{A}\)

The cardinality (or size) of \(A\) is the number of elements in \(A\text{.}\)

######  Investigate!

  1. Find the cardinality of each set below.

     1. \(A = \{3,4,\ldots, 15\}\text{.}\)

     2. \(B = \{n \in \N \st 2 \lt n \le 200\}\text{.}\)

     3. \(C = \{n \le 100 \st n \in \N \wedge \exists m \in \N (n = 2m+1)\}\text{.}\)

  2. Find two sets \(A\) and \(B\) for which \(|A| = 5\text{,}\) \(|B| = 6\text{,}\) and \(|A\cup B| = 9\text{.}\) What is \(|A \cap B|\text{?}\)

  3. Find sets \(A\) and \(B\) with \(|A| = |B|\) such that \(|A\cup B| = 7\) and \(|A \cap B| = 3\text{.}\) What is \(|A|\text{?}\)

  4. Let \(A = \{1,2,\ldots, 10\}\text{.}\) Define \(\mathcal{B}_2 = \{B \subseteq A \st |B| = 2\}\text{.}\) Find \(|\mathcal{B}_2|\text{.}\)

  5. For any sets \(A\) and \(B\text{,}\) define \(AB = \{ab \st a\in A \wedge b \in B\}\text{.}\) If \(A = \{1,2\}\) and \(B = \{2,3,4\}\text{,}\) what is \(|AB|\text{?}\) What is \(|A \times B|\text{?}\)

###  Subsection Relationships Between Sets

We have already said what it means for two sets to be equal: they have exactly the same elements. Thus, for example,

$$
\{1, 2, 3\} = \{2, 1, 3\}\text{.}
$$

(Remember, the order the elements are written down in does not matter.) Also,

$$
\{1, 2, 3\} = \{1, 1+1, 1+1+1\} = \{I, II, III\} = \{1, 2, 3, 1+2\}
$$

since these are all ways to write the set containing the first three positive integers (how we write them doesn’t matter, just what they are).

What about the sets \(A = \{1, 2, 3\}\) and \(B = \{1, 2, 3, 4\}\text{?}\) Clearly \(A \ne B\text{,}\) but notice that every element of \(A\) is also an element of \(B\text{.}\) Because of this we say that \(A\) is a _subset_ of \(B\text{,}\) or in symbols \(A \subset B\) or \(A \subseteq B\text{.}\) Both symbols are read “is a subset of.” The difference is that sometimes we want to say that \(A\) is either equal to or is a subset of \(B\text{,}\) in which case we use \(\subseteq\text{.}\) This is analogous to the difference between \(\lt\) and \(\le\text{.}\)

####  Example 0.3.3.

Let \(A = \{1, 2, 3, 4, 5, 6\}\text{,}\) \(B = \{2, 4, 6\}\text{,}\) \(C = \{1, 2, 3\}\) and \(D = \{7, 8, 9\}\text{.}\) Determine which of the following are true, false, or meaningless.

  1. \(A \subset B\text{.}\)

  2. \(B \subset A\text{.}\)

  3. \(B \in C\text{.}\)

  4. \(\emptyset \in A\text{.}\)

  5. \(\emptyset \subset A\text{.}\)

  6. \(A \lt D\text{.}\)

  7. \(3 \in C\text{.}\)

  8. \(3 \subset C\text{.}\)

  9. \(\{3\} \subset C\text{.}\)

In the example above, \(B\) is a subset of \(A\text{.}\) You might wonder what other sets are subsets of \(A\text{.}\) If you collect all these subsets of \(A\) into a new set, we get a set of sets. We call the set of all subsets of \(A\) the power set of \(A\text{,}\) and write it \(\pow(A)\text{.}\)

####  Example 0.3.4.

Let \(A = \{1,2,3\}\text{.}\) Find \(\pow(A)\text{.}\)

Another way to compare sets is by their _size_. Notice that in the example above, \(A\) has 6 elements and \(B\text{,}\) \(C\text{,}\) and \(D\) all have 3 elements. The size of a set is called the set’s cardinality . We would write \(|A| = 6\text{,}\) \(|B| = 3\text{,}\) and so on. For sets that have a finite number of elements, the cardinality of the set is simply the number of elements in the set. Note that the cardinality of \(\{ 1, 2, 3, 2, 1\}\) is 3. We do not count repeats (in fact, \(\{1, 2, 3, 2, 1\}\) is exactly the same set as \(\{1, 2, 3\}\)). There are sets with infinite cardinality, such as \(\N\text{,}\) the set of rational numbers (written \(\mathbb Q\)), the set of even natural numbers, and the set of real numbers (\(\mathbb R\)). It is possible to distinguish between different infinite cardinalities, but that is beyond the scope of this text. For us, a set will either be infinite, or finite; if it is finite, the we can determine its cardinality by counting elements.

####  Example 0.3.5.

  1. Find the cardinality of \(A = \{23, 24, \ldots, 37, 38\}\text{.}\)

  2. Find the cardinality of \(B = \{1, \{2, 3, 4\}, \emptyset\}\text{.}\)

  3. If \(C = \{1,2,3\}\text{,}\) what is the cardinality of \(\pow(C)\text{?}\)

###  Subsection Operations On Sets

Is it possible to add two sets? Not really, however there is something similar. If we want to combine two sets to get the collection of objects that are in either set, then we can take the union of the two sets. Symbolically,

$$
C = A \cup B\text{,}
$$

read, “\(C\) is the union of \(A\) and \(B\text{,}\)” means that the elements of \(C\) are exactly the elements which are either an element of \(A\) or an element of \(B\) (or an element of both). For example, if \(A = \{1, 2, 3\}\) and \(B = \{2, 3, 4\}\text{,}\) then \(A \cup B = \{1, 2, 3, 4\}\text{.}\)

The other common operation on sets is intersection. We write,

$$
C = A \cap B
$$

and say, “\(C\) is the intersection of \(A\) and \(B\text{,}\)” when the elements in \(C\) are precisely those both in \(A\) and in \(B\text{.}\) So if \(A = \{1, 2, 3\}\) and \(B = \{2, 3, 4\}\text{,}\) then \(A \cap B = \{2, 3\}\text{.}\)

Often when dealing with sets, we will have some understanding as to what “everything” is. Perhaps we are only concerned with natural numbers. In this case we would say that our universe is \(\N\text{.}\) Sometimes we denote this universe by \(\mathcal{U}\text{.}\) Given this context, we might wish to speak of all the elements which are _not_ in a particular set. We say \(B\) is the complement of \(A\text{,}\) and write,

$$
B = \bar A
$$

when \(B\) contains every element not contained in \(A\text{.}\) So, if our universe is \(\{1, 2,\ldots, 9, 10\}\text{,}\) and \(A = \{2, 3, 5, 7\}\text{,}\) then \(\bar A = \{1, 4, 6, 8, 9,10\}\text{.}\)

Of course we can perform more than one operation at a time. For example, consider

$$
A \cap \bar B\text{.}
$$

This is the set of all elements which are both elements of \(A\) and not elements of \(B\text{.}\) What have we done? We’ve started with \(A\) and removed all of the elements which were in \(B\text{.}\) Another way to write this is the set difference:

$$
A \cap \bar B = A \setminus B\text{.}
$$

It is important to remember that these operations (union, intersection, complement, and difference) on sets produce other sets. Don’t confuse these with the symbols from the previous section (element of and subset of). \(A \cap B\) is a set, while \(A \subseteq B\) is true or false. This is the same difference as between \(3 + 2\) (which is a number) and \(3 \le 2\) (which is false).

####  Example 0.3.6.

Let \(A = \{1, 2, 3, 4, 5, 6\}\text{,}\) \(B = \{2, 4, 6\}\text{,}\) \(C = \{1, 2, 3\}\) and \(D = \{7, 8, 9\}\text{.}\) If the universe is \(\mathcal{U} = \{1, 2, \ldots, 10\}\text{,}\) find:

  1. \(A \cup B\text{.}\)

  2. \(A \cap B\text{.}\)

  3. \(B \cap C\text{.}\)

  4. \(A \cap D\text{.}\)

  5. \(\bar{B \cup C}\text{.}\)

  6. \(A \setminus B\text{.}\)

  7. \((D \cap \bar C) \cup \bar{A \cap B}\text{.}\)

  8. \(\emptyset \cup C\text{.}\)

  9. \(\emptyset \cap C\text{.}\)

Having notation like this is useful. We will often want to add or remove elements from sets, and our notation allows us to do so precisely.

####  Example 0.3.7.

If \(A = \{1,2,3\}\text{,}\) then we can describe the set we get by adding the number 4 as \(A \cup \{4\}\text{.}\) If we want to express the set we get by removing the number 2 from \(A\) we can do so by writing \(A \setminus \{2\}\text{.}\)

Careful though. If you add an element to the set, you get a new set! So you would have \(B = A \cup \{4\}\) and then correctly say that \(B\) contains 4, but \(A\) does not.

You might notice that the symbols for union and intersection slightly resemble the logic symbols for “or” and “and.” This is no accident. What does it mean for \(x\) to be an element of \(A\cup B\text{?}\) It means that \(x\) is an element of \(A\) _or_ \(x\) is an element of \(B\) (or both). That is,

$$
x \in A \cup B \qquad \iff \qquad x \in A \vee x \in B\text{.}
$$

Similarly,

$$
x \in A \cap B \qquad \iff \qquad x \in A \wedge x \in B\text{.}
$$

Also,

$$
x \in \bar A \qquad \iff \qquad \neg (x \in A)\text{.}
$$

which says \(x\) is an element of the complement of \(A\) if \(x\) is not an element of \(A\text{.}\)

There is one more way to combine sets which will be useful for us: the Cartesian product, \(A \times B\text{.}\) This sounds fancy but is nothing you haven’t seen before. When you graph a function in calculus, you graph it in the Cartesian plane. This is the set of all ordered pairs of real numbers \((x,y)\text{.}\) We can do this for _any_ pair of sets, not just the real numbers with themselves.

Put another way, \(A \times B = \{(a,b) \st a \in A \wedge b \in B\}\text{.}\) The first coordinate comes from the first set and the second coordinate comes from the second set. Sometimes we will want to take the Cartesian product of a set with itself, and this is fine: \(A \times A = \{(a,b) \st a, b \in A\}\) (we might also write \(A^2\) for this set). Notice that in \(A \times A\text{,}\) we still want _all_ ordered pairs, not just the ones where the first and second coordinate are the same. We can also take products of 3 or more sets, getting ordered triples, or quadruples, and so on.

####  Example 0.3.8.

Let \(A = \{1,2\}\) and \(B = \{3,4,5\}\text{.}\) Find \(A \times B\) and \(A \times A\text{.}\) How many elements do you expect to be in \(B \times B\text{?}\)

###  Subsection Venn Diagrams

There is a very nice visual tool we can use to represent operations on sets. A Venn diagram displays sets as intersecting circles. We can shade the region we are talking about when we carry out an operation. We can also represent cardinality of a particular set by putting the number in the corresponding region.

Each circle represents a set. The rectangle containing the circles represents the universe. To represent combinations of these sets, we shade the corresponding region. For example, we could draw \(A \cap B\) as:

Here is a representation of \(A \cap \bar B\text{,}\) or equivalently \(A \setminus B\text{:}\)

A more complicated example is \((B \cap C) \cup (C \cap \bar A)\text{,}\) as seen below.

Notice that the shaded regions above could also be arrived at in another way. We could have started with all of \(C\text{,}\) then excluded the region where \(C\) and \(A\) overlap outside of \(B\text{.}\) That region is \((A \cap C) \cap \bar B\text{.}\) So the above Venn diagram also represents \(C \cap \bar{\left((A\cap C)\cap \bar B\right)}\text{.}\) So using just the picture, we have determined that

$$
(B \cap C) \cup (C \cap \bar A) = C \cap \bar{\left((A\cap C)\cap \bar B\right)}\text{.}
$$
