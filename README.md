# shelly

Store, organize, and quickly execute your favorite shell commands through a beautiful and intuitive TUI.

## Why shelly?

There's alot of reasons why you'd want to use shelly. Perhaps you have too many aliases and need a simple way to search through them, perhaps you frequently forget them, perhaps you want to store some aliases seprately from your main aliases, etc. Whatever the reason, shelly is here for it.

## Supported Commands

- `shelly`: Open the shelly TUI where you can view your commands in a table, search through them, and execute them.
- `shelly init`: Initalize shelly by creating config.json and commands.json files.
- `shelly add [ALIAS] [DESCRIPTION] [CMD]`: Add an alias with its associated description, and command.
- `shelly del [ALIAS]`: Delete an alias.
- `shelly exec [ALIAS]`: Execute the command associated with the alias.
