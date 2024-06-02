# Commander

Commander is a command-line interface (CLI) tool built in Golang that interprets user questions about command-line tasks and provides the necessary commands and arguments to accomplish them. It supports various large language models (LLMs) to generate and explain commands, enhancing user understanding and efficiency in the CLI environment.

## Features

- **Command Generation:** Generate commands based on natural language input.
- **Command Explanation:** Explain what each part of the command does.
- **Execution Option:** Directly execute the command from the tool (use with caution).
- **Clipboard Support:** Copy commands directly to the clipboard.
- **Automatic Updates:** Check for the latest updates automatically.

## Installation

```bash
# Clone the repository
git clone https://github.com/bgrewell/commander.git

# Go into the repository
cd commander

# Build the project
make build
```

## Usage

To use Commander, simply run the executable and pass your question as an argument:

```bash
commander "find all files in my home directory that contain the word 'super'"
```

### Options

- `-e, --explain` - Provide an explanation of the command.
- `-x, --exec` - Execute the command directly.
- `-c, --clip` - Copy the command to the clipboard.
- `-u, --update` - Automatically check for updates.

### Examples

**Basic Command:**

```bash
commander "find all files in my home directory that contain the word 'super'"
```

**Command with Explanation:**

```bash
commander -e "find all files in my home directory that contain the word 'super'"
```

## Contributing

We welcome contributions! Please open an issue or submit a pull request with your improvements.

## License

This project is licensed under the Creative Commons Attribution-NonCommercial 4.0 International License - see the [LICENSE](https://creativecommons.org/licenses/by-nc/4.0/legalcode) file for details.

## Authors

- **Benjamin Grewell** - *Initial work* - [bgrewell](https://github.com/bgrewell)

## Acknowledgments

- Thanks to all contributors of open source LLMs.
- Hat tip to anyone whose code was used.
- Inspiration, etc.


## Pre-Release To Do

- [ ] Create configuration file system to allow option defaults to be set via a file instead of being passed on every command
- [x] Implement an internal commands system to allow commands to be sent to commander. i.e. 'set_aliases'
- [x] Implement internal 'processor' to handle construction of messages to the LLM, output, etc. instead of having it in the main file
- [x] Clean up build code
- [ ] Create integration tests or other way to come up with a way to test conformance of the LLM to the expected behavior
- [ ] Add ability to perform multiple parallel queries and return the most common response or have another LLM build the response based on those responses
- [x] Add check for command, if they aren't on the system provide instructions on installing/offer to install them or suggest another command
  - [x] Add AI function to search for command presence on the system
- [ ] Related to above suggest commands that are present on the system unless too complicated then offer new tools
- [ ] Refactor much of the assistants code into a base class since per LLM changes are minor
- [ ] Setup update website
- [ ] Give LLM access to common ENV variables like EDITOR, HOME, VISUAL, SHELL, etc. in a more efficent way then the current method