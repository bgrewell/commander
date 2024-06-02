## 0.3.0 (2024-06-02)

### Feat

- **main.go**: add --install and --alias options to install and alias commander

### Fix

- **injection.go**: add injection of preferred editor
- **main.go**: change default behavior of update
- fixed update code

## 0.3.0 (2024-06-01)

### Feat

- add ability for llm to use tools and basic tool to check if a command exists on the system

## 0.2.1 (2024-06-01)

### Fix

- set up update url using .env file

### Refactor

- fix version in .cz.json

## 0.2.0 (2024-06-01)

### Fix

- modify update code to show new version

## v0.0.1 (2024-06-01)

## 0.1.0 (2024-06-01)

### Feat

- add exec mode
- initial interactivity and explanations

### Fix

- tweak prompt to increase accuracy
- add temperature to reduce unexpected output format

### Refactor

- move code out of main into a processor struct
- add .env to .gitignore
