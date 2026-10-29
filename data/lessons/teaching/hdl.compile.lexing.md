# Lexical analysis of HDL source

Every HDL tool starts by turning characters into tokens: keywords, identifiers, numbers, operators and punctuation. The lexer is the simplest phase of a compiler, but Verilog's operator set (=== , >>>, ~&, <= used both as assignment and comparison), its preprocessor and its comments make it easy to get subtly wrong. Counting tokens by hand is the fastest way to internalize the rules a lexer implements.

## Tokens and what produces none

A token is the smallest meaningful unit of source: a keyword such as assign or posedge, an identifier such as cnt, a number such as 47, an operator such as & or <=, or punctuation such as ( ) [ ] ; and @. Whitespace and comments produce no tokens; they only separate tokens. So assign y = a & b; is seven tokens, and so is assign y=a&b;.

Before lexing proper, Verilog runs a preprocessor that expands macros, includes files and evaluates conditional compilation (ifdef). Because macros can produce arbitrary token sequences, a front end must lex the preprocessed text, and good tools keep a map from each token back to its original file and line so that error messages point at the source the user wrote.

## Maximal munch

A lexer reads greedily: at each position it takes the longest character sequence that forms a valid token. This is the maximal-munch rule. In a<=b the characters <= form one token, giving three tokens a, <=, b. In a===b the lexer takes ===, not == followed by =. Other multi-character operators include !==, <<, >>, <<<, >>>, &&, ||, **, >= and the reduction operators ~&, ~| and ~^.

The parser, not the lexer, decides what <= means. In q <= d; after an always block it is a nonblocking assignment; inside if (a <= b) it is a comparison. The lexer emits the same token in both cases, and the grammar resolves it from context. Counting tokens in if ( a <= b ) y = 1 ; gives 10, whatever the meaning.

## Comments, selects and literals

A line comment starting with // runs to the end of the line, and a block comment between /* and */ can sit anywhere between tokens; both are discarded, even when they contain text that looks like code, such as a <= b. A bit select rdy[0] is four tokens: rdy, [, 0 and ]. A plain decimal number is one token.

Sized literals such as 8'hFF are a classic lexer subtlety, because the size, the base and the digits may legally be separated by whitespace, so lexers differ in whether they treat them as one token or several; the parser then reassembles them. A lexer must also handle escaped identifiers, which begin with a backslash and end at whitespace, and keywords that depend on the language version selected, which is why tools take a language-version option before they read any file.

## Going deeper

Aho, Lam, Sethi and Ullman, "Compilers: Principles, Techniques, and Tools" (2nd ed.), chapter 3; IEEE Std 1800-2017, clause 5 (lexical conventions).
