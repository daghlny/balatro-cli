# Balatro (小丑牌) 简化版游戏规则

> 本文档为 Terminal UI 第一版开发使用的精简规则。删除了消耗品（塔罗/星球/幽灵/券码）、卡牌强化/版本、进阶难度（赌注/贴纸）等系统，保留核心玩法循环、计分系统和小丑牌机制。

---

## 目录

1. [核心游戏循环](#1-核心游戏循环)
2. [扑克牌型系统](#2-扑克牌型系统)
3. [计分系统](#3-计分系统)
4. [卡牌系统](#4-卡牌系统)
5. [盲注系统](#5-盲注系统)
6. [商店系统](#6-商店系统)
7. [小丑牌列表](#7-小丑牌列表)
8. [数值速查表](#8-数值速查表)

---

## 1. 核心游戏循环

### 1.1 游戏概述

Balatro 是一款 roguelike 牌组构建扑克游戏。玩家使用标准扑克牌打出牌型获得分数，通过购买小丑牌强化构建，挑战不断升级的目标分数。

### 1.2 完整游戏流程

一局完整的游戏由 **8 个底注(Ante)** 组成，击败底注 8 即获胜。

```
GAME_INIT:
  money = $4
  deck = standard_52_cards (shuffled)
  joker_slots = [] (max 5)
  hand_type_levels = {each_type: 1}
  hand_type_play_counts = {each_type: 0}

for ante in 1..8:
  for blind in [SmallBlind, BigBlind, BossBlind]:

    // --- 盲注选择 ---
    if blind != BossBlind:
      player chooses: CHALLENGE or SKIP
      if SKIP:
        award random skip_reward (见 5.4)
        continue to next blind  // 跳过该盲注，放弃奖励和商店
    
    // --- 出牌阶段 (见 1.4) ---
    target = calc_blind_target(ante, blind)
    run_blind_round(target, blind.boss_effect)
    // 若回合失败 → GAME OVER

    // --- 结算阶段 ---
    money += blind_base_reward + remaining_hands * $1 + interest

    // --- 商店阶段 (见 第6节) ---
    run_shop()

GAME WIN  // 击败底注 8 的 Boss 盲注
```

### 1.3 单个底注内的盲注顺序

每个底注包含 3 个盲注，按顺序挑战：
- **小盲注 (Small Blind)**: 目标 = 1× 基础分。可跳过。
- **大盲注 (Big Blind)**: 目标 = 1.5× 基础分。可跳过。
- **Boss盲注 (Boss Blind)**: 目标 = 2× 基础分，附带特殊效果。不可跳过。

跳过代价：放弃该盲注的金钱奖励和商店机会。补偿：随机获得一项跳过奖励（见 5.4）。

### 1.4 出牌阶段详解

```
function run_blind_round(target_score, boss_effect):
  shuffle deck
  hand = draw(8) from deck
  remaining_hands = 4   // 出牌机会
  remaining_discards = 3 // 弃牌机会
  round_score = 0
  apply boss_effect (if any)

  loop:
    player chooses action:

      ACTION: PLAY (选择 1-5 张牌)
        remaining_hands -= 1
        played_cards = selected cards
        hand_type = identify_best_hand(played_cards)
        score = calculate_score(hand_type, played_cards, jokers)  // 见第3节
        round_score += score
        // 更新牌型出牌计数，每 3 次自动升 1 级
        hand_type_play_counts[hand_type] += 1
        if hand_type_play_counts[hand_type] % 3 == 0:
          hand_type_levels[hand_type] += 1
        move played_cards to discard_pile
        draw cards to refill hand to hand_size

      ACTION: DISCARD (选择最多 5 张牌)
        remaining_discards -= 1
        move selected cards to discard_pile
        draw equal number of new cards from deck
        // 不计分

    // 牌组抽空时，将弃牌堆洗牌作为新牌组
    if deck is empty: deck = shuffle(discard_pile)

    if round_score >= target_score: return SUCCESS
    if remaining_hands == 0: return GAME_OVER
```

**关键规则**：
- 每回合开始抽 **8 张**手牌
- **4 次出牌机会**，**3 次弃牌机会**
- 出牌选 **1-5 张**，系统识别牌型并计分
- 弃牌选**最多 5 张**，不计分，从牌组补充等量新牌
- 打出或弃掉的牌进入弃牌堆
- 牌组抽空时，弃牌堆洗牌成为新牌组
- 得分在回合内**累加**，达到目标即胜利，无需用完所有出牌机会

### 1.5 结算阶段

击败盲注后获得金钱奖励：

```
function calc_reward(blind_type, remaining_hands, current_money):
  base = {SmallBlind: 3, BigBlind: 4, BossBlind: 5}[blind_type]
  hand_bonus = remaining_hands * 1
  interest = min(floor(current_money / 5), 5)
  return base + hand_bonus + interest
```

**示例**：击败大盲注，剩余 2 次出牌，持有 $17 → $4 + $2 + floor(17/5) = $4 + $2 + $3 = **$9**

---

## 2. 扑克牌型系统

### 2.1 牌型表

| 牌型 | 英文名 | 基础筹码 | 基础倍率 | 说明 |
|------|--------|---------|---------|------|
| 高牌 | High Card | 5 | 1 | 无法构成其他牌型时，最大的一张牌 |
| 对子 | Pair | 10 | 2 | 2 张相同点数 |
| 两对 | Two Pair | 20 | 2 | 2 组不同的对子 |
| 三条 | Three of a Kind | 30 | 3 | 3 张相同点数 |
| 顺子 | Straight | 30 | 4 | 5 张连续点数 (A可作为1或14) |
| 同花 | Flush | 35 | 4 | 5 张相同花色 |
| 葫芦 | Full House | 40 | 4 | 三条 + 对子 |
| 四条 | Four of a Kind | 60 | 7 | 4 张相同点数 |
| 同花顺 | Straight Flush | 100 | 8 | 顺子 + 同花 |
| 皇家同花顺 | Royal Flush | 100 | 8 | 10-J-Q-K-A 同花顺 |

### 2.2 牌型升级

每种牌型都有**等级**，等级越高基础数值越大。

**升级方式**：每累计打出同一种牌型 **3 次**，该牌型自动升 1 级。

| 牌型 | 每级 +筹码 | 每级 +倍率 |
|------|-----------|-----------|
| 高牌 | +10 | +1 |
| 对子 | +15 | +1 |
| 两对 | +20 | +1 |
| 三条 | +20 | +2 |
| 顺子 | +30 | +3 |
| 同花 | +15 | +2 |
| 葫芦 | +25 | +2 |
| 四条 | +30 | +3 |
| 同花顺 | +40 | +4 |
| 皇家同花顺 | +40 | +4 |

**示例**：对子打出 6 次 → 等级 3 → 基础筹码 10+2×15=**40**，基础倍率 2+2×1=**4**

### 2.3 牌型识别规则

- 系统自动识别打出牌中**最高级别**的牌型
- 牌的打出顺序不影响识别
- 只有**构成牌型的牌**参与计分
  - 例：打出 K♠ K♥ 7♣ 3♦ 2♠ → 识别为对子，仅 K♠ K♥ 参与计分
  - 例外：拥有 Splash 小丑时，所有打出的牌都参与计分
- 顺子中 A 可以作为最高 (10-J-Q-K-A) 或最低 (A-2-3-4-5)

---

## 3. 计分系统

### 3.1 核心公式

`最终得分 = 总筹码(Chips) × 总倍率(Mult)`

### 3.2 计分流程

```
function calculate_score(hand_type, played_cards, jokers) -> int:
  level = hand_type_levels[hand_type]
  chips = hand_type.base_chips + (level - 1) * hand_type.level_chips
  mult  = hand_type.base_mult  + (level - 1) * hand_type.level_mult

  // 阶段一: 基础值已在上面初始化

  // 阶段二: 打出的牌计分 (从左到右)
  scoring_cards = get_scoring_cards(played_cards, hand_type)
  // 默认只有构成牌型的牌参与计分
  // 例外: Splash 小丑使所有打出的牌都参与计分
  for card in scoring_cards (left to right):
    chips += card.rank_chips   // 见 3.3 点数筹码表

  // 阶段三: 小丑牌效果 (从左到右，顺序影响结果)
  for joker in jokers (left to right):
    apply joker.effect to (chips, mult)
    // 效果类型: +chips, +mult, 或 ×mult
    // 建议排列顺序: +chips → +mult → ×mult

  return chips * mult
```

**重要**：小丑牌的排列顺序影响最终得分。`+chips` 和 `+mult` 是加法，`×mult` 是乘法，乘法类小丑放在最右侧收益最大。

### 3.3 卡牌点数筹码值

| 点数 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | J | Q | K | A |
|------|---|---|---|---|---|---|---|---|----|---|---|---|---|
| 筹码 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 10 | 10 | 10 | 11 |

### 3.4 计分示例

打出 **K♠ K♥ 5♣**，对子 (等级 1)，拥有 Jolly Joker (+8 mult if pair)：

```
chips = 10 (对子基础), mult = 2 (对子基础)
scoring_cards = [K♠, K♥]  // 5♣不是牌型构成牌，不计分
K♠: chips += 10 → chips = 20
K♥: chips += 10 → chips = 30
Jolly Joker: 包含对子 → mult += 8 → mult = 10
result = 30 × 10 = 300
```

---

## 4. 卡牌系统

### 4.1 标准扑克牌

标准牌组包含 **52 张牌**：4 种花色 × 13 种点数。

| 花色 | 符号 | 颜色 |
|------|------|------|
| 黑桃 Spades | ♠ | 黑 |
| 红心 Hearts | ♥ | 红 |
| 梅花 Clubs | ♣ | 黑 |
| 方块 Diamonds | ♦ | 红 |

**点数**: 2, 3, 4, 5, 6, 7, 8, 9, 10, J, Q, K, A

**人脸牌 (Face Cards)**: J, Q, K

### 4.2 被削弱 (Debuffed)

被削弱的牌（由 Boss 盲注效果造成）：
- **不提供筹码**
- 但**仍然可以构成牌型**（如被削弱的 K 仍可以组成对子）

---

## 5. 盲注系统

### 5.1 基础筹码需求

| 底注 | 基础分 | 小盲注 (×1) | 大盲注 (×1.5) | Boss盲注 (×2) |
|------|--------|------------|--------------|--------------|
| 1 | 300 | 300 | 450 | 600 |
| 2 | 800 | 800 | 1,200 | 1,600 |
| 3 | 2,000 | 2,000 | 3,000 | 4,000 |
| 4 | 5,000 | 5,000 | 7,500 | 10,000 |
| 5 | 11,000 | 11,000 | 16,500 | 22,000 |
| 6 | 20,000 | 20,000 | 30,000 | 40,000 |
| 7 | 35,000 | 35,000 | 52,500 | 70,000 |
| 8 | 50,000 | 50,000 | 75,000 | 100,000 |

### 5.2 Boss 盲注

每个底注的 Boss 盲注**随机**从以下列表中选取（底注 8 使用决战盲注）。

| Boss盲注 | 英文 | 效果 |
|---------|------|------|
| 棍棒 | The Club | 所有**梅花**牌被削弱 |
| 刺棒 | The Goad | 所有**黑桃**牌被削弱 |
| 头颅 | The Head | 所有**红心**牌被削弱 |
| 窗户 | The Window | 所有**方块**牌被削弱 |
| 植物 | The Plant | 所有**人脸牌**被削弱 |
| 灵媒 | The Psychic | 必须打出 **5 张牌**才能计分 |
| 镣铐 | The Manacle | 手牌上限 **-1** |
| 燧石 | The Flint | 基础筹码和基础倍率**减半** |
| 墙壁 | The Wall | 超大盲注（**4× 基础分**而非 2×） |
| 针 | The Needle | 本回合只能出 **1 手**牌 |
| 水 | The Water | 本回合**没有弃牌**机会 |
| 嘴巴 | The Mouth | 本回合只能打出**同一种牌型** |
| 眼睛 | The Eye | 每手牌必须是**不同的牌型** |
| 手臂 | The Arm | 打出的牌型**降低 1 级**（最低1级） |
| 牙齿 | The Tooth | 每打出 1 张牌 **-$1** |

### 5.3 决战盲注（底注 8 专属）

| 决战盲注 | 英文 | 效果 |
|---------|------|------|
| 紫罗兰之器 | Violet Vessel | 巨大盲注（**6× 基础分** = 300,000） |
| 深红之心 | Crimson Heart | 每手牌**随机禁用 1 张**小丑牌 |
| 蔚蓝之钟 | Cerulean Bell | 始终有 **1 张随机手牌被强制选中** |

### 5.4 跳过奖励

跳过小盲注或大盲注时，随机获得以下奖励之一：

| 奖励 | 效果 |
|------|------|
| 金钱奖励 | 获得 **$5** |
| 刷新奖励 | 下一个商店**首次刷新免费** |
| 小丑奖励 | 下一个商店出现 **1 张免费小丑牌** |
| 经济奖励 | 当前金钱 **翻倍**（上限 $20） |
| 手牌奖励 | 下一个回合 **+3 手牌上限** |

---

## 6. 商店系统

### 6.1 商店内容

每次进入商店提供以下内容：
- **2 张小丑牌**（按稀有度权重随机：70% 普通、25% 罕见、5% 稀有）
- **1 个丑角补充包**（从多张小丑牌中选取，见 6.3）

```
function run_shop():
  shop_items = [random_joker(), random_joker(), random_booster_pack()]
  reroll_cost = 5

  loop until player chooses LEAVE:
    player can:
      BUY item     → if money >= item.cost and has_slot: money -= item.cost, add to inventory
      SELL joker   → money += floor(joker.cost / 2), remove from slots
      REROLL       → if money >= reroll_cost: money -= reroll_cost, regenerate 2 jokers, reroll_cost += 1
      LEAVE        → exit shop
```

### 6.2 价格体系

| 物品 | 价格 |
|------|------|
| 普通小丑 (Common) | $2-$6 |
| 罕见小丑 (Uncommon) | $5-$8 |
| 稀有小丑 (Rare) | $8-$10 |

**售卖价值**：`floor(购买价格 / 2)`

### 6.3 丑角补充包

| 等级 | 价格 | 内容 |
|------|------|------|
| 普通丑角包 | $4 | 2 张选 1 |
| 大型丑角包 | $6 | 4 张选 1 |
| 超级丑角包 | $8 | 4 张选 2 |

### 6.4 刷新机制

- 初始刷新费用：**$5**
- 同一商店每次刷新：费用 **+$1**
- 进入新商店时：刷新费用**重置为 $5**
- 刷新替换 2 张小丑牌，不替换补充包

---

## 7. 小丑牌列表

共 **50 张**小丑牌，分为 4 个稀有度。

### 7.1 普通小丑 (Common) — 出现权重 70%

| # | 名称 | 费用 | 效果 |
|---|------|------|------|
| 1 | Joker | $2 | **+4** 倍率 |
| 2 | Greedy Joker | $5 | **方块**牌计分时 +3 倍率 |
| 3 | Lusty Joker | $5 | **红心**牌计分时 +3 倍率 |
| 4 | Wrathful Joker | $5 | **黑桃**牌计分时 +3 倍率 |
| 5 | Gluttonous Joker | $5 | **梅花**牌计分时 +3 倍率 |
| 6 | Jolly Joker | $3 | 牌型包含**对子**时 +8 倍率 |
| 7 | Zany Joker | $4 | 牌型包含**三条**时 +12 倍率 |
| 8 | Crazy Joker | $4 | 牌型包含**顺子**时 +12 倍率 |
| 9 | Droll Joker | $4 | 牌型包含**同花**时 +10 倍率 |
| 10 | Sly Joker | $3 | 牌型包含**对子**时 +50 筹码 |
| 11 | Wily Joker | $4 | 牌型包含**三条**时 +100 筹码 |
| 12 | Devious Joker | $4 | 牌型包含**顺子**时 +100 筹码 |
| 13 | Crafty Joker | $4 | 牌型包含**同花**时 +80 筹码 |
| 14 | Half Joker | $5 | 打出 **≤3 张牌**时 +20 倍率 |
| 15 | Banner | $5 | 每剩余 1 次弃牌 **+30 筹码** |
| 16 | Mystic Summit | $5 | 弃牌次数为 **0** 时 +15 倍率 |
| 17 | Even Steven | $4 | **偶数牌** (2,4,6,8,10) 计分时 +4 倍率 |
| 18 | Odd Todd | $4 | **奇数牌** (3,5,7,9,A) 计分时 +31 筹码 |
| 19 | Scholar | $4 | **A** 计分时 +20 筹码和 +4 倍率 |
| 20 | Scary Face | $4 | **人脸牌**计分时 +30 筹码 |
| 21 | Smiley Face | $4 | **人脸牌**计分时 +5 倍率 |
| 22 | Misprint | $4 | +0 到 +23 **随机**倍率 |
| 23 | Supernova | $5 | +倍率 = 本局该牌型累计出牌次数 |
| 24 | Blue Joker | $5 | 牌组中每张**剩余牌** +2 筹码 |
| 25 | Ice Cream | $5 | +100 筹码，每出 1 手牌 **-5** 筹码 (减到 0 消失) |
| 26 | Splash | $3 | 所有打出的牌都**参与计分** |
| 27 | Abstract Joker | $4 | 每拥有 1 张小丑牌 **+3 倍率** |
| 28 | Juggler | $4 | **+1** 手牌上限 |
| 29 | Drunkard | $4 | 每回合 **+1** 弃牌 |
| 30 | Popcorn | $5 | +20 倍率，每回合 **-4** 倍率 (减到 0 消失) |

### 7.2 罕见小丑 (Uncommon) — 出现权重 25%

| # | 名称 | 费用 | 效果 |
|---|------|------|------|
| 1 | Four Fingers | $7 | **4 张牌**即可构成同花和顺子 |
| 2 | Ride the Bus | $6 | 连续不打出人脸牌时每手 **+1 倍率**（累积，打出人脸牌后归零） |
| 3 | Loyalty Card | $5 | 每 **6 手**牌 ×4 倍率 (第6手触发) |
| 4 | Fibonacci | $8 | **A/2/3/5/8** 计分时 +8 倍率 |
| 5 | Blackboard | $6 | 手中全是**黑桃或梅花**时 ×3 倍率 |
| 6 | Card Sharp | $6 | 本回合**重复**打出同种牌型时 ×3 倍率 |
| 7 | Photograph | $5 | 第一张**人脸牌**计分时 ×2 倍率 |
| 8 | Raised Fist | $5 | +倍率 = 手中**最低点数**牌的 2 倍 |
| 9 | Green Joker | $4 | 每出 1 手牌 **+1 倍率**，每弃 1 次牌 **-1 倍率**（累积，最低 0） |
| 10 | Acrobat | $6 | 回合**最后一手**牌 ×3 倍率 |
| 11 | Flash Card | $5 | 每次商店**刷新** +2 倍率（永久累积） |
| 12 | Spare Trousers | $6 | 包含**两对**时 +2 倍率（永久累积） |
| 13 | Bull | $6 | 每持有 **$1** 获得 +2 筹码 |
| 14 | Golden Joker | $6 | 回合结束时获得 **$4** |
| 15 | Flower Pot | $6 | 计分牌包含**全部 4 种花色**时 ×3 倍率 |
| 16 | Smeared Joker | $7 | 红心=方块(红色同算)，黑桃=梅花(黑色同算) |
| 17 | Chaos the Clown | $4 | 每个商店 **1 次免费刷新** |
| 18 | Dusk | $5 | 回合**最后一手**牌重新触发所有计分牌 |

### 7.3 稀有小丑 (Rare) — 出现权重 5%

| # | 名称 | 费用 | 效果 |
|---|------|------|------|
| 1 | Baron | $8 | 手中每张 **K** ×1.5 倍率 |
| 2 | Blueprint | $10 | **复制**右侧相邻小丑的能力 |
| 3 | The Duo | $8 | 牌型包含**对子**时 ×2 倍率 |
| 4 | The Trio | $8 | 牌型包含**三条**时 ×3 倍率 |
| 5 | The Family | $8 | 牌型包含**四条**时 ×4 倍率 |
| 6 | The Order | $8 | 牌型包含**顺子**时 ×3 倍率 |
| 7 | The Tribe | $8 | 牌型包含**同花**时 ×2 倍率 |
| 8 | Campfire | $9 | 每卖出 1 张牌 ×0.25 倍率（累积），击败 Boss 盲注后**重置** |
| 9 | Baseball Card | $8 | 每张**罕见小丑** ×1.5 倍率 |
| 10 | DNA | $8 | 首手打出**单张牌**时，复制该牌到牌组并抽取 |

### 7.4 传奇小丑 (Legendary) — 不在商店出现，仅通过特殊方式获得

> 简化版中，传奇小丑有 **极低概率** (0.5%) 出现在超级丑角补充包中。

| # | 名称 | 效果 |
|---|------|------|
| 1 | Triboulet | **K 和 Q** 计分时 ×2 倍率 |
| 2 | Chicot | **禁用**所有 Boss 盲注效果 |

---

## 8. 数值速查表

### 8.1 默认参数

| 参数 | 默认值 |
|------|--------|
| 手牌上限 | **8** |
| 每回合出牌次数 | **4** |
| 每回合弃牌次数 | **3** |
| 小丑牌槽位 | **5** |
| 每次打出牌数 | **1-5 张** |
| 每次弃牌数 | **最多 5 张** |
| 起始金钱 | **$4** |
| 利息公式 | **floor($/5)**，上限 $5 |

### 8.2 经济参考

- 基础奖励: SmallBlind=$3, BigBlind=$4, BossBlind=$5
- 手数奖励: remaining_hands * $1
- 利息: min(floor(money / 5), 5)
- 典型收入(底注1-3): 省2手+$12存款 = $3+$2+$2 = $7

### 8.3 小丑牌功能分类速查

```
加筹码 (+Chips):  Sly/Wily/Devious/Crafty, Banner, Odd Todd, Scary Face,
                  Scholar, Ice Cream, Blue Joker, Bull

加倍率 (+Mult):   Joker, Greedy/Lusty/Wrathful/Gluttonous,
                  Jolly/Zany/Crazy/Droll, Half Joker, Mystic Summit,
                  Even Steven, Smiley Face, Misprint, Supernova,
                  Abstract, Fibonacci, Raised Fist, Flash Card,
                  Spare Trousers, Popcorn, Ride the Bus, Green Joker

乘倍率 (×Mult):  Loyalty Card, Blackboard, Card Sharp, Photograph,
                  Acrobat, Dusk, Flower Pot, The Duo/Trio/Family/Order/Tribe,
                  Baron, Campfire, Baseball Card, Blueprint, Triboulet

经济类 ($):       Golden Joker

功能类:           Four Fingers, Splash, Smeared Joker, DNA,
                  Juggler, Drunkard, Chaos the Clown, Chicot
```

### 8.4 与完整版的主要差异

| 特性 | 完整版 | 简化版 |
|------|--------|--------|
| 消耗品 (塔罗/星球/幽灵) | 52 张 | ❌ 移除 |
| 券码 | 32 张 | ❌ 移除 |
| 卡牌强化 | 8 种 | ❌ 移除 |
| 卡牌版本 | 5 种 | ❌ 移除 |
| 封印 | 4 种 | ❌ 移除 |
| 牌型升级方式 | 使用星球卡 | 打出 3 次自动升级 |
| 赌注难度 | 8 级 | 仅白色赌注 (基础难度) |
| 贴纸 | 永恒/易腐/租赁 | ❌ 移除 |
| 小丑牌 | 150 张 | 50 张 |
| Boss 盲注 | 23 种 | 15 种 |
| 决战盲注 | 5 种 | 3 种 |
| 标签系统 | 24 种 | 简化为 5 种跳过奖励 |
| 牌组选择 | 15 种 | 标准牌组 (52张) |
| 补充包 | 5 类 | 仅丑角包 |
| 商店内容 | 小丑+消耗品+券码+补充包 | 小丑+丑角包 |

---

> 本文档基于 Balatro v1.0.1 完整规则精简而成。完整规则参见 `game_rules.md`。
