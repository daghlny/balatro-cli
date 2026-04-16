# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Terminal UI clone of Balatro (小丑牌), a poker-themed roguelike. Written in Go using the Charmbracelet stack (Bubble Tea / Lip Gloss). The default locale is `zh-CN`.

The authoritative ruleset being implemented is `simple_game_rules.md` (a simplified subset of `game_rules.md` — no consumables, enhancements, stakes, or stickers in this build).

## Build / run / test

```bash
go build -o balatro-cli .         # produces ./balatro-cli
go run .                           # run directly
go test ./...                      # run full test suite
go test ./engine -run TestHand     # run one package / match a test name
go test ./engine -v -run TestPlayHand   # single test, verbose
```

There is no separate lint step configured — `go vet ./...` and `go build` are the available checks.

## Architecture

The code is split into three layers with a one-way dependency direction: **`engine` (pure logic) ← `jokers` (content registration) ← `ui` (Bubble Tea presentation)**. `main.go` wires them together and calls `tea.NewProgram(model, tea.WithAltScreen())`.

### `engine/` — pure game logic, no I/O, no UI

`GameState` (`engine/game.go`) is the single source of truth for an active run. Its `Phase` field (`PhaseTitle`, `PhaseBlindSelect`, `PhaseBossInfo`, `PhasePlay`, `PhaseScoring`, `PhaseSettlement`, `PhaseShop`, `PhaseGameOver`) drives UI routing. State transitions are methods on `GameState` (`StartGame`, `EnterBlind`, `PlayHand`, `DiscardSelected`, `FinishScoring`, `FinishSettlement`, `LeaveShop`, `AdvanceBlind`, `SkipBlind`, `SellJoker`). The UI calls these; it must not mutate `GameState` fields directly except for the documented cases (the play screen sets `Phase = PhaseSettlement/PhaseGameOver` inline after `PlayHand` — see `ui/screens/play.go`).

Scoring runs a **3-phase pipeline** in `engine/score.go::CalculateScore`:
1. Base chips/mult from `HandTypeDataTable[ht]` adjusted by per-hand-type level.
2. Per-scoring-card chips (`Rank.Chips()`, zero if `Debuffed`).
3. Per-joker effects left-to-right: for each joker, accumulate `OnCardScore` for every non-debuffed scoring card then `OnScore` once, then apply `PlusChips`, `PlusMult`, and `MultMult` to the running totals.

Randomness flows through `engine.RNG` (`engine/rng.go`). Use `NewRNG()` in production and `NewSeededRNG(seed)` in tests — every stochastic engine path takes `*RNG` as a parameter so tests are deterministic.

Boss effects (`engine/boss.go`) implement the `BossEffect` interface and hook into `Apply` (mutate `RoundState` at round start), `TargetMultiplier`, `IsDebuffed(card)`, `ValidatePlay`, and `OnPlay`. `RandomBoss(ante, rng)` picks from `AllRegularBosses` for ante 1–7 and `AllFinaleBosses` at ante 8.

Hand detection (`engine/hand.go::DetectHand`) takes `HandDetectOpts` carrying the aggregate of passive joker flags (`FourFingers`, `SmearedJoker`, `Splash`) — these are collected by `GetPassiveEffects(jokers)` in `engine/joker.go` and applied both at detection time and when computing hand size / discards at round start.

### `jokers/` — joker content, registered via `init()`

Each joker is a `*engine.JokerDef` declared as a package-level var plus a `JokerEffect` implementation that embeds `engine.BaseJokerEffect` and overrides only the hooks it needs (`OnCardScore`, `OnScore`, `OnHandPlayed`, `OnDiscard`, `OnRoundEnd`, `OnSell`). Files are grouped by rarity (`common.go`, `uncommon.go`, `rare.go`, `legendary.go`). **All defs must be registered in `jokers/registry.go::init()`** — the package is imported for its side effects via `_ "github.com/yinuoli/balatro-cli/jokers"` in `main.go`. Without that blank import the registry is empty and the shop has nothing to sell.

Passive-only jokers (e.g. Four Fingers, Juggler) leave `EffectFactory` nil and declare their behavior via the `Passive PassiveEffect` field on `JokerDef`.

### `ui/` — Bubble Tea presentation

`ui.AppModel` (`ui/app.go`) is the root `tea.Model`. It holds a `map[engine.Phase]screens.Screen` and dispatches `Update`/`View` to whichever screen matches `game.Phase`. Each `Screen` (`ui/screens/screen.go`) implements:
- `Update(msg, *GameState) tea.Cmd` — handles input and mutates `GameState` via engine methods.
- `View(*GameState, width, height) string` — renders the screen body (status bar is appended by `AppModel.View`).
- `Focus() components.FocusTarget` — describes what the cursor is pointing at so the shared status bar can render context.

Screens are constructed once and reused across rounds. `PlayScreen` caches `tableCards`/`tableBD` for the "last played hand" middle panel; it clears them when `game.LastBreakdown` is nil (see the comment in `ui/screens/play.go::View` — the engine resets `LastBreakdown` in `startRound`, and the screen mirrors that).

`ui/components/` holds reusable widgets: `card.go` (playing-card rendering), `joker_card.go`, `sidebar.go` (full left rail with logo, ante, target, score, hands/discards, money, hand-type reference table), `statusbar.go` (context + shortcut cheatsheet), `bignum.go` (big Chips × Mult digits), `theme.go` (color constants — `ColorGold`, `ColorRed`, `ColorBlue`, `ColorCyan`, `ColorDim`, `ColorWhite`).

### `i18n/` — localization

`i18n.T(key)` is the only lookup function; `i18n.SetLocale(code)` switches the active table. Locale tables are registered in `init()` (`i18n/zh_cn.go`, `i18n/en_us.go`) and loaded via a blank import in `main.go`. Default locale is `zh-CN`. All user-visible strings in the engine and UI use i18n keys (hand types as `hand.*`, suits as `suit.*`, bosses as `boss.<id>.name`/`.desc`, jokers as `joker.<id>.name`/`.desc`, UI chrome as `ui.*`).

### `testutil/` — test helpers

`testutil.MakeCard("Ks")` / `MakeCards("Ah","10d",...)` build cards from shortcodes (rank + suit letter `s/h/c/d`). `testutil.NewTestRNG(seed)` returns a deterministic `*engine.RNG`. Prefer these over constructing cards by hand in tests.

## Adding content

- **New joker**: add a `JokerDef` + effect in the appropriate `jokers/<rarity>.go`, register it in `jokers/registry.go`, add `joker.<id>.name` / `.desc` entries to every file in `i18n/`.
- **New boss**: add a type implementing `engine.BossEffect` in `engine/boss.go`, append its constructor to `AllRegularBosses` or `AllFinaleBosses`, add `boss.<id>.name` / `.desc` entries to the i18n tables.
- **New screen/phase**: add a `Phase` constant in `engine/game.go`, implement `screens.Screen` under `ui/screens/`, register it in `ui.NewAppModel`'s screen map. Engine transitions should set `game.Phase` to that constant; avoid routing decisions in the UI layer.

## Change log

After completing a non-trivial feature, fix, or refactor, append an entry to `.claude_change_log.md` at the repo root. One entry per topic, newest on top. Use this template:

```markdown
## YYYY-MM-DD — Short title (verb + object)

- One-line summary of what changed and why it matters.
- Follow-up bullets for non-obvious details, gotchas, or related work.
- **Touched:** comma-separated list of files or directories.
```

Skip the log for tiny edits (typos, whitespace, single-line tweaks) — git history is enough for those. Log anything that another contributor would benefit from seeing at a glance: new systems, UI redesigns, algorithm changes, content packs, dependency upgrades, build/release tooling.
