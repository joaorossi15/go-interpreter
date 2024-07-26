# Parser

- a parser is a software component that takes input data (frequently text) and builds a data structure
- normally builds a parse tree or abstract syntax tree
- this gives structural representation of the input
- so, summarizing, a parser turns input data into a data structure that represents the input

## Abstract Syntax Tree (AST)
    - normally newlines, braces, semicolons, brackets and parentheses only guide the parser but do not appear
    - but there is no universal format to an AST
    - example using JS
    ```
    if (3 * 5 > 10) {
        return "hello";
        } else {
        return "goodbye";
    }

    // parsed result
    {
        type: "if-statement",
        condition: {
            type: "operator-expression",
            operator: ">",
            left: {
                type: "operator-expression",
                operator: "*",
                left: { type: "integer-literal", value: 3 },
                right: { type: "integer-literal", value: 5 }
            },
            right: { type: "integer-literal", value: 10 }
        },
        consequence: {
            type: "return-statement",
            returnValue: { type: "string-literal", value: "hello" }
        },
        alternative: {
            type: "return-statement",
            returnValue: { type: "string-literal", value: "goodbye" }
        }
    }
    ```

## Types
    - top-down or bottom-up
    - recursive descent, early passing, predictive parsing, etc

## Monkey Parser
    - let statement: let <identifier> = <expression>
    - statement doenst produce value, expression does (return 5 doesnt produce values, but add(2, 5) does)
    - in Monkey a lot of things are expressions, including function literals
