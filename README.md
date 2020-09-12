# Yaml to Erlang config converter.

A YAML file
```
array:
    - "string"
    - !!binary "binary_string"
    - !atom "atom"
    - 123
    - 456.7
    - true
    - false
    - null
    - {"foo": "fooVal", "bar": !atom "barVal" }
    -
        foo: fooVal
        bar: !atom barVal

object:
    a_string: "string"
    a_binray_string: !!binary "binary_string"
    an_atom: !atom "atom"
    an_int: 123
    a_float: 456.7
    a_bool1: true
    a_bool2: false
    a_nil: null
    an_array: ["foo", !!binary "baz", !atom "bar" ]
    a_tuple: !tuple ["foo", !!binary "baz", !atom "bar" ]
```

becomes Erlang config:
```
[
  {array, [
    "string",
    <<"binary_string">>,
    atom,
    123,
    456.7,
    true,
    false,
    nil,
    [{foo, "fooVal"}, {bar, barVal}],
    [
      {foo, "fooVal"},
      {bar, barVal}
    ]
  ]},
  {object, [
    {a_string, "string"},
    {a_binray_string, <<"binary_string">>},
    {an_atom, atom},
    {an_int, 123},
    {a_float, 456.7},
    {a_bool1, true},
    {a_bool2, false},
    {a_nil, nil},
    {an_array, ["foo", <<"baz">>, bar]},
    {a_tuple, {"foo", <<"baz">>, bar}}
  ]}
].
```

## Installation

* Go get
```bash
go get github.com/makasim/erltoyml/main
```

## Usage

The command outputs erlang config version to its stdout.
```bash
ymltoerl path/to/file.yaml
```

