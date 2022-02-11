# Define

A CLI for finding the definition of a word.

## Requirements

You'll need to get an API key for the [Lingua Robot
API](https://rapidapi.com/rokish/api/lingua-robot/endpoints). That needs to be
in your environment variable as `LINGUA_ROBOT_API_KEY`.

## Installing

```
# compiles the CLI into ./bin/go-define
make
```

If you'd like to be able to run this from any place, move the binary to your
$PATH. It's recommended that you use `$HOME/.bin/`, if your shell supports it.

## Testing

```
make test
```
