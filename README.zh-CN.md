# balatro-cli

[Balatro](https://www.playbalatro.com/)（小丑牌）的终端 UI 复刻版。基于 Go + Charmbracelet（Bubble Tea / Lip Gloss）实现。

[![English](https://img.shields.io/badge/lang-English-lightgrey?style=for-the-badge)](./README.md)
[![中文](https://img.shields.io/badge/lang-中文-blue?style=for-the-badge)](./README.zh-CN.md)

> **说明**：本项目是学习性质的爱好者作品，原版游戏设计版权归 LocalThunk / Playstack 所有。

---

## 功能

- 完整的 8 轮 Ante 循环：小盲 → 大盲 → Boss 盲注，之间夹带商店。
- 扑克牌型识别 + 牌型升级（从高牌一路到同花五）。
- 小丑商店，手持小丑从左到右按顺序结算效果，可卖出回血。
- Boss 盲注效果（减益、出牌限制、倍率加成等）。
- 默认简体中文，内置英文语言表，可切换。
- 纯键盘操作，不需要鼠标。

## 安装

一行安装（macOS / Linux，支持 amd64 / arm64）：

```sh
curl -fsSL https://raw.githubusercontent.com/daghlny/balatro-cli/main/install.sh | sh
```

脚本会从 GitHub 拉取最新 Release 包，解压后安装到 `~/.local/bin/balatro-cli`。如果想装到别处，可以设置环境变量：`INSTALL_DIR=/usr/local/bin`。

然后运行：

```sh
balatro-cli
```

注意把 `~/.local/bin` 加到 `PATH`（在 shell rc 里加 `export PATH="$HOME/.local/bin:$PATH"`）。

### 其他安装方式

- **Homebrew / Scoop**：尚未发布。
- **Windows**：去 [Releases 页面](https://github.com/daghlny/balatro-cli/releases) 下载 `.zip`，把 `balatro-cli.exe` 放到 `PATH` 下任意目录。
- **源码编译**（需要 Go 1.26+）：
  ```sh
  git clone https://github.com/daghlny/balatro-cli
  cd balatro-cli
  go build -o balatro-cli .
  ./balatro-cli
  ```

## 玩法

从标题界面开始一局，然后循环每个 Ante：

1. **选择盲注** — 打当前盲，或花一次跳过次数跳过。
2. **出牌阶段** — 从手牌里最多选 5 张：
   - `Enter` → 打出（跟盲注目标分比）。
   - `x` → 弃掉并补抽（每轮弃牌次数有限）。
   - `←/→` 或 `h/l` → 移动光标。
   - `Space` → 选中 / 取消。
3. **结算** — 看 筹码 × 倍率 滚动完成，按任意键继续。
4. **商店** — 花钱买小丑（从左到右就是结算顺序），可以刷新或直接离开。
5. **循环**，直到第 8 轮 Boss 被击败，或者分数没达到盲注目标就结束。

快捷键提示常驻在底部状态栏。

## 项目结构

三层，单向依赖：**`engine`（纯逻辑）← `jokers`（内容）← `ui`（Bubble Tea）**。新增小丑、Boss 或界面的细节见 [`CLAUDE.md`](./CLAUDE.md)。

## 许可证

MIT。如有 [`LICENSE`](./LICENSE) 文件，以其为准。
