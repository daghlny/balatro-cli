# balatro-cli

A terminal-UI clone of [Balatro](https://www.playbalatro.com/) — the poker-themed roguelike — written in Go on the Charmbracelet stack (Bubble Tea + Lip Gloss).

[![English](https://img.shields.io/badge/lang-English-blue?style=for-the-badge)](./README.md)
[![中文](https://img.shields.io/badge/lang-中文-lightgrey?style=for-the-badge)](./README.zh-CN.md)

> **Note**: This is a fan project for learning purposes. All credit for the original game design goes to LocalThunk / Playstack.

---

## Features

- Full 8-ante run loop — small blind → big blind → boss blind, shops in between.
- Poker hand detection with per-hand-type levelling (high card → flush five).
- Joker shop, held jokers with left-to-right scoring effects, sell for money.
- Boss blind effects (debuffs, skip-conditions, chip multipliers).
- `zh-CN` default locale with an `en-US` table built in.
- Keyboard-only TUI — no mouse required.

## Install

One-liner (macOS / Linux, amd64 / arm64):

```sh
curl -fsSL https://raw.githubusercontent.com/daghlny/balatro-cli/main/install.sh | sh
```

This pulls the latest release archive from GitHub, extracts the binary, and installs it to `~/.local/bin/balatro-cli`. Override with `INSTALL_DIR=/usr/local/bin` if you prefer a different location.

Then run:

```sh
balatro-cli
```

Make sure `~/.local/bin` is on your `PATH` (`export PATH="$HOME/.local/bin:$PATH"` in your shell rc).

### Other ways to install

- **Homebrew / Scoop**: not yet published.
- **Windows**: download the `.zip` from [Releases](https://github.com/daghlny/balatro-cli/releases) and put `balatro-cli.exe` somewhere on your `PATH`.
- **From source** (Go 1.26+):
  ```sh
  git clone https://github.com/daghlny/balatro-cli
  cd balatro-cli
  go build -o balatro-cli .
  ./balatro-cli
  ```

## How to play

Start a run from the title screen, then cycle through each ante:

1. **Blind select** — pick the current blind or spend a skip to jump it.
2. **Play** — highlight up to 5 cards from your hand, then:
   - `Enter` → play them (scored against the blind target).
   - `x` → discard them and draw replacements (limited discards per round).
   - `←/→` or `h/l` → move the cursor.
   - `Space` → toggle selection.
3. **Scoring** — watch chips × mult resolve, then continue.
4. **Shop** — spend money on jokers (left to right = scoring order). Reroll or leave.
5. **Repeat** until ante 8 boss falls, or your score misses the target.

Controls cheat-sheet is always visible in the status bar at the bottom.

## Project layout

Three layers with a strict one-way dependency: **`engine` (pure logic) ← `jokers` (content) ← `ui` (Bubble Tea)**. See [`CLAUDE.md`](./CLAUDE.md) for details on adding a joker, boss, or screen.

## License

MIT. See [`LICENSE`](./LICENSE) if present.
