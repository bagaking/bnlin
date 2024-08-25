# bnlin - Natural Language Bash Script Helper

bnlin (short for Bagaking Natural Language bINary) is a command-line tool that
uses an AI driver to draft and execute bash scripts from natural language input.
It is intended for users who want a faster starting point for shell operations
while still reviewing what will run in their own environment.

## Safety Model

bnlin is an execution helper, not a sandbox. It turns natural language into a
bash script candidate and runs that script with the permissions and environment
of the current shell session.

- **Generated bash can be destructive**: Review the generated execution plan and
  command intent before using bnlin on important files, production systems, or
  shared machines.
- **Use a controlled environment**: Prefer a disposable checkout, temporary
  directory, container, VM, or other limited workspace when the requested task
  may modify files or system state.
- **Review command boundaries**: Treat filesystem writes, deletes, network
  calls, package installation, service control, and permission changes as
  high-risk operations that need manual inspection.
- **Protect secrets**: Do not put access keys, tokens, passwords, private file
  contents, or other sensitive values into natural-language prompts, shell
  history, logs, issue reports, or screenshots.
- **Driver output is not trusted**: The AI driver may produce incorrect,
  overbroad, or platform-incompatible commands. The user remains responsible
  for deciding whether the generated script should run.

## Features

- **Natural Language Input**: Convert plain English instructions into bash script candidates.
- **Platform Awareness**: Detects the operating system and includes that context when generating scripts.
- **Flexible Configuration**: Set up via command-line flags or environment variables.
- **Script Execution**: Generates and runs scripts from the command line.
- **Readable Output**: Prints command results in a formatted terminal view.

## Platform support

bnlin generates bash scripts and executes them with bash semantics. Linux and
macOS are the primary supported environments. On Windows, use bnlin from an
environment that provides bash-compatible script execution, such as WSL or Git
Bash. Native Windows shell semantics such as PowerShell and `cmd.exe` are not
the execution target.

## Installation

To install bnlin, ensure you have Go installed, then run:

```bash
go install github.com/bagaking/bnlin@latest
```

## Local validation

Run the project test suite with:

```bash
make test
```

This target runs:

```bash
go test ./...
```

## Usage

The basic syntax for using bnlin is:

```bash
bnlin run <your command in natural language>
```

### Examples

1. List uncommitted files and their line counts:

```bash
   bnlin run find all uncommitted files and list their line counts
```

2. View folders in the parent directory:

```bash
   bnlin run "show me all folders in the parent directory"
```

### Configuration

The driver can be configured using command-line flags:

```bash
bnlin run --driver doubao <your command in natural language>
```

By default, the driver is `doubao`; therefore, the following command is equivalent to the previous one:

```bash
bnlin run <your command in natural language>
```

#### Ollama

When using the `ollama` driver, an endpoint is required. It identifies the model to use;
for example, it can be `llama3.1`.

You can configure the endpoint using command-line flags:

```bash
bnlin run --driver ollama -e llama3.1 "your command here"
```

#### Doubao

bnlin can be configured using command-line flags or environment variables:

- **Access Key**:
    - Flag: `-ak` or `--access_key`
    - Env: `VOLC_ACCESSKEY`

- **Secret Key**:
    - Flag: `-sk` or `--secret_key`
    - Env: `VOLC_SECRETKEY`

- **API Endpoint**:
    - Flag: `-e` or `--endpoint`
    - Env: `DOUBAO_ENDPOINT`

Recommended environment-based configuration keeps credentials out of the
interactive command line:

```bash
export VOLC_ACCESSKEY=<access-key>
export VOLC_SECRETKEY=<secret-key>
export DOUBAO_ENDPOINT=<endpoint>
bnlin run "your command here"
```

The credential flags remain available for controlled, non-interactive use, but
environment variables are the safer default for regular shell sessions.

For more detailed information and advanced usage, run:

```bash
bnlin --help
```

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

bnlin is built with the following excellent libraries:

- [github.com/bagaking/botheater](https://github.com/bagaking/botheater)
- [github.com/bagaking/easycmd](https://github.com/bagaking/easycmd)

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/bagaking/bnlin/issues) on our GitHub repository.
