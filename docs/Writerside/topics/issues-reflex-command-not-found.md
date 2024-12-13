# Reflex command not found

**Reflex command not found (`zsh: command not found: reflex`)**

* open your terminal and run `ls $HOME/go/bin/reflex` to check if reflex is installed
* if it is installed run `open ~/.zprofile` or `~/.zshrc` to open your shell configuration
* your shell config may have different name
* add these lines to your config

```bash
# go path
export PATH="$PATH:$HOME/go/bin"
```