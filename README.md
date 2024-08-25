# bnlin - AI-Powered Bash Script Generator

bnlin (short of Bagaking Nature Language bINary), is an innovative command-line tool that leverages 
AI to automatically generate and execute bash scripts based on natural language input. It seamlessly 
translates user requests into precise bash commands, making complex system operations accessible to 
users of all skill levels.

## Features

- **Natural Language Processing**: Convert plain English instructions into executable bash scripts.
- **Cross-Platform Compatibility**: Supports Windows, Linux, and MacOS.
- **Flexible Configuration**: Set up via command-line flags or environment variables.
- **Intelligent OS Detection**: Adapts commands to the specific operating system environment.
- **Real-time Execution**: Generates and runs scripts on-the-fly.
- **User-Friendly Output**: Provides clear, formatted results for easy understanding.

## Installation

To install bnlin, ensure you have Go installed, then run:

```bash
go install github.com/bagaking/bnlin@latest
```

## Usage

The basic syntax for using bnlin is:

```bash
bnlin run "<your command in natural language>"
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

bnlin can be configured using command-line flags or environment variables:

- **Access Key**:
    - Flag: `-ak` or `--access_key`
    - Env: `VOLC_ACCESS_KEY`

- **Secret Key**:
    - Flag: `-sk` or `--secret_key`
    - Env: `VOLC_SECRET_KEY`

- **API Endpoint**:
    - Flag: `-e` or `--endpoint`
    - Env: `DOU_BAO_ENDPOINT`

Example with flags:
```bash
bnlin run -ak your_access_key -sk your_secret_key -e your_endpoint "your command here"
```

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