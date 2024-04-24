# gusset

Gusset is a compile-to-JS language with syntax akin to Go. It shares many features with Go, making it familiar to programmers versed in both languages.

<img style="text-align:center; margin-bottom: 20px;" src="https://avatars.githubusercontent.com/u/166648156?s=200" />

## Goals

1. Facilitate writing clear and correct code for JavaScript applications.
2. Improve on JavaScript's lack of runtime guarantees with strict types and compile-time checks.
3. Expand upon Go's functionality for frontend development with features like list/map comprehension, enums, pattern matching, and native JSON and JSX.

To achieve this, Gusset employs Go's type system and syntax supplemented with additions (*gussets*) from Rust, Ruby, and JavaScript itself.

## Getting Started

Install the `gusset` and `gus` CLIs to start writing and building a Gusset program. The `gusset` CLI manages versions of `gus` on your machine. When you install `gusset`, the latest version of the language and corresponding `gus` CLI is installed.

### homebrew (macOS or Linux)

```sh
brew install gusset
```

### Windows

```
```

### Install script

```
```

### Create a Module

Use the `gus` CLI to initialize a project with some default files.

```sh
mkdir my-app && cd my-app
gus init github.com/myuser/project
```

The init command creates a `main.gus` file with a hello world program. Build the program using the `build` command:

```sh
gus build
```

The `build` command by default outputs a set of JavaScript files to the `dist` folder in the project directory.

### Running a Gusset Program

Gusset does not include a way to execute JavaScript. Running your program will depend on the execution environment for the type of JS app you're building.

For more details on common execution setups, check out the [website docs](https://gusset.dev/docs).

## Learning Gusset

The most up-to-date way to learn Gusset is to read the [Language Guide](https://gusset.dev/docs). The guide stays up to date as the language is evolved.

In the future, I plan on making a video series on writing Gusset programs.

## Contributing

Interested in helping build Gusset? Check out [the Contributing guide](/CONTRIBUTING.md).

