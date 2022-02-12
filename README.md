# `go-define`

A CLI for finding the definition of a word.

## Requirements

Before starting, get a **[Lingua Robot API key](https://rapidapi.com/rokish/api/lingua-robot/endpoints)**.

By default, the CLI will look for a config file in `$HOME/.config/go-define/config.yml` with the key
**LINGUA_ROBOT_API_KEY**.

The config file should look like:

```yml
LINGUA_ROBOT_API_KEY: <Key>
```

If you want to specify a custom config file, pass the config path as a global flag:

```
go-define -c path/to/custom/config.yml copacetic
```

## Installing

```
# compiles the CLI into ./bin/go-define
make
```

If you'd like to be able to run this from any place, move the binary to your
`$PATH`. It's recommended that you use `$HOME/.bin/`, if your shell supports it.

## Testing

```
make test
```
