# Inferring carry chains, DSPs and RAMs

Inference recognizes RTL patterns that belong on hard resources: additions on carry chains, multiplications on DSP blocks, arrays on block RAM or LUT RAM. A good inference pass saves orders of magnitude in area and delay; a missed one silently builds the function in soft logic. Its rules come from the hard blocks' exact behavior, so RTL style decides what is inferred.

## What maps where

An addition or comparison becomes LUTs that compute per-bit propagate and generate signals feeding carry-chain primitives. A multiplication, especially one with registered inputs and an accumulation, becomes one or more DSP blocks. An array becomes memory: large ones with synchronous reads go to block RAM, small ones with asynchronous reads go to LUT RAM, and very small ones may stay as flip-flops. In Yosys these steps are passes such as alumacc and techmap for arithmetic, per-family DSP passes, and memory mapping against a description of the family's RAM primitives. The rest of the logic then goes to the LUT mapper. Getting inference right first matters, because LUT mapping cannot later turn soft logic back into a hard block.

## Coding style decides inference

Hard blocks implement specific behavior, and the RTL must match it. Block RAM reads are synchronous, so a memory whose read data is used in the same cycle as its address cannot use block RAM; registering the address or the output, at the cost of one cycle of latency, makes it legal. A memory that clears every word on reset cannot be a RAM, because a RAM writes at most a port's worth of words per cycle. Registers around a multiplier can be absorbed into the DSP block's input, product and output registers, which makes the multiplier faster without fabric flip-flops. A multiplication wider than the DSP's ports is split into partial products across several blocks and summed through their post-adders and cascades. Knowing these rules is how a designer gets the inference they want.

## Consequences downstream

Inference decisions constrain the rest of the flow. Carry chains must be placed as vertical macros, so the placer moves them as rigid shapes. Hard blocks sit in fixed columns, so the placer legalizes them separately from LUTs. Hard resources are a budget: once the DSP blocks are used up, the remaining multiplies fall back to LUTs and carry chains, larger and slower. Tiles are allocated whole, so a 512 by 40 memory in a 512 by 72 tile wastes almost half of it. Disabling the carry chain for an adder frees its placement but slows every carry through LUTs and routing. A tool designer should expose these choices as attributes and report utilization per resource, because they are design decisions, not mapping details.

## Going deeper

The Yosys manual, chapters on memory mapping and technology mapping; AMD's synthesis user guide for Vivado (UG901), sections on RAM and DSP inference.
