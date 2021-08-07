# PoE Macro

I wanted a /hideout macro when playing on Linux, so I made this super basic one. It needs sudo privileges to work, as required by github.com/MarinX/keylogger. L_CTRL + D is added as a way to pause the program so that pressing F5 doesn't trigger the macro, used when alt-tabbing out of PoE but leaving the program running.

The macro is set up for Colemak keyboard layout as github.com/MarinX/keylogger reads/writes keys based on keycodes rather than the actual value.

# Compiling and running

```bash
go build
sudo ./poe-macro
```