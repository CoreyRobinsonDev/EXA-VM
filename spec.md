# Keywords
#### Host
- 1st line must be the `HOST` name
`HOST <name>`

#### LINK
- `LINK`s are optional
- `LINK`s must be -9999 and 9999
`LINK <internal ID> <external ID>`

#### COPY
`COPY R/N R`

#### MAKE
- create and grab a file with `MAKE`
- `MAKE` takes an optional filename as it's only argument
```exa
MAKE <filename>
```

# Example
```exa
HOST main
Link 800 200

COPY 100 X
```
