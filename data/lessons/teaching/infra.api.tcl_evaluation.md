# Tcl evaluation and substitution

Most EDA tools, OpenROAD and OpenSTA included, are driven by Tcl, and SDC constraint files are Tcl programs. Tcl's evaluation rules are short but unusual: everything is a string, a command is a list of words, and substitution happens once, before the command runs, in ways that braces and quotes control. Tool builders must know these rules to write command parsers, to read SDC correctly, and to avoid subtle scripting bugs such as integer division that differs from C.

## Words and lists

A Tcl command is a sequence of words separated by white space; the first word names the command. A Tcl list has the same syntax: elements separated by spaces, with braces grouping several words into one element. The command llength counts top-level elements. In the list u1 {a b c} u2 {} {x {y z}}, written inside outer braces, there are five elements: u1, the group a b c, u2, an empty element, and the group x {y z}. Nesting does not matter to llength; a braced group is one element however much it contains, and an empty pair of braces is a legal empty element. Tool commands that return collections of objects, such as the get commands of SDC, return lists or list-like handles, and scripts iterate them with foreach. A list built by string concatenation can split wrongly when names contain spaces, so scripts should build lists with the list command.

## Substitution in quotes and braces

Before a command runs, Tcl performs one round of substitution on each word: a dollar sign followed by a variable name is replaced by the variable's value, and a bracketed command such as [set w] is replaced by that command's result. Double quotes group a word but still allow substitution. Braces group a word and suppress all substitution. So after set w U42, the command set s "inst_[set w]/Z" stores inst_U42/Z, 10 characters, while set s {inst_[set w]/Z} stores the literal text inst_[set w]/Z, 14 characters. Substitution happens exactly once: if a value contains brackets, they are not evaluated again. This matters for EDA names, which often contain brackets for bus bits such as data[3]: a script must brace or escape them, or Tcl will try to run a command named 3.

## Integer arithmetic in expr

The expr command evaluates arithmetic. Its argument should be braced, as in expr {a + b} with variables, so that expr itself performs substitution once and can compile the expression; unbraced expressions are substituted twice and are slower and less safe. Integer division in Tcl rounds toward negative infinity, and the remainder takes the sign of the divisor, so that $a = b \cdot q + r$ always holds. Thus expr {17 / 5} is 3 and expr {17 % 5} is 2, but expr {-17 / 5} is -4 and expr {-17 % 5} is 3, and expr {17 % -5} is -3. C and C++ truncate toward zero instead, giving -3 and -2 for the negative cases. A grid-snapping formula ported from C++ to a Tcl script can therefore snap negative coordinates differently. Mixing in a floating-point operand, as in 7 / 2.0, gives 3.5.

## Going deeper

John Ousterhout and Ken Jones, "Tcl and the Tk Toolkit" (2nd ed., 2009), the chapters on syntax and expressions; the Tcl manual page Tcl(n), which states the evaluation rules.
