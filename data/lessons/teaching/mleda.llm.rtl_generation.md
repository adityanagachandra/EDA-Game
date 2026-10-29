# LLMs for RTL generation

Large language models can write Verilog from a natural-language specification. The output is a sample from a learned distribution over tokens, which may compile, may behave as intended, and may synthesize to good hardware, three separate questions. A tool designer treats generated RTL as a proposal and builds the checks that answer each question.

## From specification to Verilog

A transformer conditioned on a prompt, the specification plus any context, samples a module token by token. Nothing in that process simulates or synthesizes the code. Sampling at nonzero temperature gives different completions on each try, and different samples succeed on different problems, so evaluations draw several samples per problem and a checker picks among them. Fine-tuning on open-source Verilog shifts the output toward that corpus's idioms, including its bugs. An ambiguous specification, such as a FIFO with no stated reset style, becomes an arbitrary choice in the output, so prompts for generation need the precision of a real design specification.

## Grounding in tools

Each tool in the loop checks a different property. A compiler or linter, with its errors fed back to the model, improves syntactic and elaboration correctness, and nothing else. Functional correctness needs an independent statement of intent: a testbench with reference outputs, or an equivalence check against a reference model. Retrieving the real port lists and parameter names of neighboring modules into the prompt prevents interface mismatches, a common failure when the model guesses names and widths. Design quality, meaning area, timing and power, appears only after synthesis, so generated code that passes simulation but synthesizes to three times the area has failed a check no testbench runs.

## Where generated RTL fails

The costly failures compile cleanly. Blocking assignments in a clocked block that models a shift register can make simulation shift by several stages per clock. A combinational case statement without a default infers a latch. A module tested at one parameter width hides bugs at others. Reset polarity and clock-domain assumptions drift from the intended interface. Lint, synthesis reports and tests across the parameter range catch these, and code review catches what tests miss. No generated RTL should enter a design branch unchecked: it enters after lint, simulation or formal checks against an independent specification, like any other contribution.

## Going deeper

Liu et al., "VerilogEval: Evaluating Large Language Models for Verilog Code Generation" (ICCAD 2023); Cummings, "Nonblocking Assignments in Verilog Synthesis, Coding Styles That Kill!" (SNUG 2000).
