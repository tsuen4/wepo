# wepo

POST to the webhook URL set in config.ini.

## Usage

- Generate `config.ini`
  - `cp config.example.ini config.ini`
- Set the URL of the Webhook in `webhook_url` in `config.ini`
  - Multiple destinations can be set by adding a section (e.g. `[sec1]`)
- The following keys have a fixed priority and are loaded in the order of section, global, default
  - `payload`: JSON format to be sent(default: `{"content": "{input}"}`)
  - `char_limit`: Character limit to be sent(default: `1024`)
- The arguments specified at run time or the values ​​of the standard input are sent.

```sh
# Use args
./wepo example

# Use stdin
cat example.txt | ./wepo

# Use other destinations
./wepo -s sec1 example
```

- Run in TUI mode

```sh
# Run in TUI mode
./wepo -t

# Run in TUI mode with other destinations
./wepo -t -s sec1
```
