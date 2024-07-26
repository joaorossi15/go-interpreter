# Lexical Analysis

## Lexer
    - transform source code to be able to parse it
    - first transformation: source code to tokens
    - called lexical analysis or lexing
    - done by a lexer (tokenizer)
    - some tokens are reserved keywords, like let or fn
    - some tokens are symbols, like = or +
    - some tokens have values assigned as well, like INTEGER(5)
    - example
    ```
    let x = 5 + 5;
    // result passing through a lexer
    [
        LET,
        IDENTIFIER(x),
        ASSIGN,
        INTEGER(5),
        PLUS,
        INTEGER(5),
        SEMICOLON
    ]
    ```
