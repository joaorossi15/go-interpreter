# Monkey Language

## Features
    - c like syntax
    - variable bindings
    - int and bool
    - arithmetic expressions
    - built-in functions
    - first-class and higher-order functions
    - closures
    - string data structure
    - array
    - hash map

## Var bind
    ```
    let age = 11;
    let name = "monkey";
    let result = 10 * (20 / 2);
    let arr = [1, 2, 3, 4];
    let dict = {"name": "biel", "age": 25};
    ```

## Access array and dict
    ```
    arr[0] // -> 1
    dict["name"] // -> biel
    ```

## Functions
    ```
    f(a, b); // call
    
    // function as first class
    let fibonacci = fn(x) {
        if (x == 0) {
            0
        } else {
            if (x == 1) {
                1
            } else {
                fibonacci(x - 1) + fibonacci(x - 2);
            }
        }
    };

    // higher order

    let twice = fn(f, x) {
        return f(f(x))
    };

    let add = fn(x) {
        return x + 2;
    };

    twice(add, 2); // -> 6
    ```
