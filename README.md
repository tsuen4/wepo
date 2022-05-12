# wepo

POST contents to the webhook URL set in config.ini.

## Installation

```sh
go install github.com/tsuen4/wepo@latest
```

## Settings

- Create `config.ini` to `$HOME/.config/wepo` directory

```ini
; default destinations
webhook_url=https://[webhook_url]

; optional
payload={"content": "{input}"}
char_limit=1024

; ------------------------
; add other destinations
[sec1]
webhook_url=https://[sec1]/[webhook_url]
; optional
payload={"content": "prefix {input}"}

; ------------------------
; add other destinations
[sec2]
webhook_url=https://[sec2]/[webhook_url]
; use default payload settings
; payload=
```

- Set the URL of the Webhook in `webhook_url` in `config.ini`
  - Multiple destinations can be set by adding a section (e.g. `[sec1]`)
- The following keys have a fixed priority and are loaded in the order of section, default, initial settings
  - `payload`: JSON format to be sent(initial settings: `{"content": "{input}"}`)
    - The `{input}` part will be replaced by the value passed.
  - `char_limit`: Character limit to be sent(initial settings: `1024`)
- The arguments specified at run time or the values of the standard input are sent.

## Usage: Shell mode

```sh
# Use args
wepo hello

# Use stdin
cat example.txt | wepo

# Use other destinations
wepo -s sec1 example
```

## Usage: Run in TUI mode

```sh
# Run in TUI mode
wepo -t

# Run in TUI mode with other destinations
wepo -t -s sec1

# Press Ctrl + C to exit
```
