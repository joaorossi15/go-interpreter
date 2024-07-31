# Expressions in Monkey

- expressions can have prefix operators
    - `-5`
    - `!true`
    - `!false`

- expressions can have infix operators (binary)
    - `5 + 5`
    - `5 * 5`

- there is also comparisons operators
    - `x == y`
    - `x >= y`
    - `x != y`

- there is also grouping and order of evaluation by using parentheses
    - `5 * (5 + 5)`
    - `((5 + 5) * 5) * 5`

- call expressions
    - `add(5, 5)`
    - `max(5, add(5, (5 * 5)))`

- identifiers are also expressions
    - `x * y / z`
    - `add(x, y)`

- since functions are first-class citizes, their literals are also expressions
    - `let add = fn(x, y) { return x + y }`
    - `fn(x, y) {return x + y}(5, 5)`
    - `(fn (x) {return x}(5) + 10) * 10`

- finally we also have if expressions
    - `let res = if(10 > 5) {true} else {false} // res = true`

