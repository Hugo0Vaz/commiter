# commiter
Commiter is a simple agentic commit generator.

# Installation

1. Clone the repository.
2. Install dependencies: go, make.
3. Run `make install`.

# Usage

The application has two main ways to work:

- `commiter`: running `commiter` all files staged are analyzed one by one, each generating a single commit.
- `commiter all`: runnig `commiter all` all files staged are analyzed as whole, generating a sigle commit with  _all_ changes.

## Options

- `--output` flag controls when passed, outputs the commit message instead of running, the available output types are: `json` and `cmd`.
- `--lang` flag controls the language of the commit message, the available languages are: `pt-br` and `en`.
- `--y` flag makes the program not ask for confirmation to commit.
