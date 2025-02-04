- **R**                 A register
- **R/N**               A register, or a number between -9999 and 9999
- **L**                 A label defined by a **MARK** pseudo-instruction

# REGISTERS
### X
The **X** register is a general-purpose storage register and can store a number or a keyword.
### T
The **T** register is a general-purpose storage register and can store a number or a keyword. It is also the destination for **TEST** instructions, and is the criterion for conditional jumps (TJMP and FJMP).
### F
The **F** register allows an EXA to read and write the contents of a held file. When an EXA grabs a file, its "file cursor" will be set to the first value in hte file. Reading from the **F** register will read this value; writing to the **F** register will overwrite this value. After reading or writing the **F** register, the file cursor will automatically advance. Writing to the end of the file will append a new value instead of overwriting. 
### M
The **M** register controls an EXA's message-passing functionality. When an EXA writes to the **M** register the value will be stored in that EXA's outgoing message slot until another EXA reads from the **M** register and receives the previously written value. Both numbers and keywords can be transferred in this way.

If an EXA writes to the **M** register, it will pause execution until that value is read by another EXA. If an EXA reads from the **M** register, it will pause execution until a value is available to be read. If two or more EXAs attempt to read from another EXA at the same time (or vice versa), one will succeed but which one succeeds will be unpredictable.

By default, an EXA can communicate with any other EXA in the same network. This can be restircted to EXAs in the same network. This can be restricted to EXAs in the same host by executing a **MODE** instruction. An EXA in global mode cannot communicate with an EXA in local mode, even if they are in the same host.

# INSTRUCTIONS
## Manipulating Values

- **COPY R/N R**        Copy the value of the first operand into the second operand.
- **ADDI R/N R/N R**    Add the value of the first operand to the value of the second operand and store the result in the third operand.
- **SUBI R/N R/N R**    Subtract the value of the second operand from the value of the first operand and store the result in the third operand.
- **MULI R/N R/N R**    Multiply the value of the first operand to the value of the second operand and store the result in the third operand.
- **DIVI R/N R/N R**    Divide the value of the first operand by the value of the second operand and store the result in the third operand.
- **MODI R/N R/N R**    Modulo the value of the first operand by the value of the second operand and store the result in the third operand.

- **SWIZ R/N R/N R**    Swizzle the value of the first operand using the value of the second operand as a swizzle mask and store the result in the third operand. The swizzle insturction can be used to rearrange and/or extract the digits in a number as shown:

| Input | Mask | Result |
|-------|------|--------|
| 6789  | 4321 | 6789   |
| 6789  | 1234 | 9876   |
| 6789  | 3333 | 7777   |
| 6789  | 1211 | 9899   |
| 6789  | -4321| -6789  |
| -6789 | -4321| 6789   |
| 6789  | 2000 | 8000   |
| 6789  | 0001 | 0009   |

## Branching

- **MARK L**            Mark this line with the specified label. **MARK** is a pseudo-instruction and is not executed.
- **JUMP L**            Jump to the specified label.
- **TJMP L**            Jump to the specified label if the **T** register equals 1 (or any value other than 0). This corresponds to a **TEST** result that was true.
- **FJMP L**            Jump to the specified label if the **T** register equals 0. This corresponds to a **TEST** result that was false.

## Testing Values

**TEST R/N = R/N**    Compare the value of the first operand to the value of the second operand. If they are equal, set the **T** register to 1, otherwise set the **T** register to 0.

|                 | TEST R/N = R/N | TEST R/N < R/N | TEST R/N > R/N |
|-----------------|----------------|----------------|----------------|
| Number/Number   | Test equality  |Test numerical order|Test numerical order|
| Keyword/Keyword | Test equality  |Test alphabetical order|Test alphabetical order|
| Number/Keyword  | Always false   | Always false   | Always false   |

## Lifecycle

- **REPL L**            Create a copy of this EXA and jump to the specified label in the copy. If an EXA is holding a file when executing a **REPL** instruction the file will not be copied and will remain held by the original EXA.
- **HALT**              Terminate this EXA. If it was holding a file, the file is dropped.
- **KILL**              Terminate another EXA in the same host as this EXA, prioritizing EXAs created by the same user. If there is more than one possible target, the target wil be chosen in a unpredictable manner.

## Movement

- **LINK R/N**          Traverse the link with the specified ID.
- **HOST R**            Copy the name of the current host into the specified register.

## Communication

- **MODE**              Toggle the **M** register between global and local mode.
- **VOID M**            Read and discard a value from the **M** register.
- **TEST MRD**          If this EXA could read from another EXA without pausing, set the **T** register to 1, otherwise set the **T** register to 0.

## File Manipulation

- **MAKE**              Create and grab a new file.
- **GRAB R/N**          Grab the file with the specified ID.
- **FILE R**            Copy the ID of the held file into the specified register. 
- **SEEK R/N**          Move the file cursor forward (positive) or backward (negative) by the specified number of values.
- **VOID F**            Remove the value at the file cursor from the currently held file.
- **DROP**              Drop the currently held file.
- **WIPE**              Delete the currently held file.
- **TEST EOF**          If the file pointer is currently at the end of the held file, set the **T** register to 1, otherwise set the **T** register to 0.

## Miscellaneous

- **NOTE**              Any text following the **NOTE** pseudo-instruction will be discar4ed when compiling, allowing it to be used to write "comments" to document the code. Any test following a semicolon, anywhere on a line, will also be discarded.
- **NOOP**              Do nothing for one cycle
- **RAND R/N R/N R**    Generate a random number between the first and second operand (inclusive) and store the result in the third operand.
