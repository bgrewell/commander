# Commander

Commander is a tool to ...

## Installation

```bash
install commander
```

## Usage

```bash
commander
```

## Release To Do

- [ ] Create configuration file system to allow option defaults to be set via a file instead of being passed on every command
- [ ] Implement an internal commands system to allow commands to be sent to commander. i.e. 'set_aliases'
- [x] Implement internal 'processor' to handle construction of messages to the LLM, output, etc. instead of having it in the main file
- [ ] Clean up build code
- [ ] Create integration tests or other way to come up with a way to test conformance of the LLM to the expected behavior
- [ ] Add ability to perform multiple parallel queries and return the most common response or have another LLM build the response based on those responses
- [ ] Add check for command, if they aren't on the system provide instructions on installing/offer to install them or suggest another command
  - [ ] Add AI function to search for command presence on the system
- [ ] Related to above suggest commands that are present on the system unless too complicated then offer new tools