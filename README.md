# pj

[WIP] JSON Parser written in Go

## BNF

```
json     = value
value    = object | array | string | number | boolean | 'null'
object   = '{' property (', ' property)* '}'
array    = '[' value (', ' value)* ']'
property = string ":" value

string   = '"' characters '"'
number   = [0-9]*
boolean  = 'true' | 'false'
```
