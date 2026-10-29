# Canonical forms and matching signatures

To decide quickly whether a cut function matches a cell, the matcher maps every function to a canonical representative of its class and compares representatives with a hash lookup. Computing the canonical form exactly costs up to $2^{n+1} n!$ transformations, so practical matchers combine it with signatures that prune most candidates cheaply. These are the data structures behind ABC's Boolean matching.

## Applying a transformation to a truth table

Write a truth table as an integer whose bit m is f at the assignment with binary code m, x1 the least significant bit. To compute g(x1, x2, x3) = f(NOT x3, x2, NOT x1), take each assignment, evaluate the arguments, and copy the bit of f: for x = (0, 1, 0), m = 2, the arguments are (1, 1, 1), so g's bit 2 equals f's bit 7. Repeating for all 8 assignments gives g. In software this is done with bit-parallel swaps and flips of truth-table words, so a transformation costs a few machine operations for n <= 6.

## The canonical representative

A common canonical form is the smallest truth-table integer over all transformations in the group. For the permutation group of 3 inputs, compute the 6 permuted tables and take the minimum: a function with tables 169, 201 and 225 over its permutations has P-canonical form 169. For full NPN on 2 inputs, minimize over 16 transformations: every AND-like function (one or three ones) has canonical form 1, every XOR-like function 6. Two functions match exactly when their canonical forms are equal, so the library is stored as a hash table from canonical form to the cells and transformations that realize it.

## Signatures that prune the search

Exhaustive canonicalization is affordable for small n but wasteful. Signatures are cheap invariants: the number of ones of f, and of each cofactor $f_{x_i}$ and $f_{\bar{x}_i}$. A permutation reorders the cofactor counts but does not change their multiset, so two functions with different multisets cannot be P-equivalent. Input negation swaps $f_{x_i}$ and $f_{\bar{x}_i}$, and output negation complements counts, so NPN signatures use unordered pairs and polarity-normalized counts. Signatures also fix variable orders: a canonicalizer sorts inputs by cofactor counts and tries only permutations among inputs with equal counts, which is how semi-canonical forms run in microseconds.

## Going deeper

Debnath and Sasao, "Efficient computation of canonical form for Boolean matching in large libraries", ASP-DAC 2004; Huang, Wang, Mishchenko et al., "Fast Boolean matching based on NPN classification", ICFPT 2013.
