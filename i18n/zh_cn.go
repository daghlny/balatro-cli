package i18n

func init() {
	Register("zh-CN", map[string]string{
		// Hand types
		"hand.high_card":      "高牌",
		"hand.pair":           "对子",
		"hand.two_pair":       "两对",
		"hand.three_of_kind":  "三条",
		"hand.straight":       "顺子",
		"hand.flush":          "同花",
		"hand.full_house":     "葫芦",
		"hand.four_of_kind":   "四条",
		"hand.straight_flush": "同花顺",
		"hand.royal_flush":    "皇家同花顺",

		// Suits
		"suit.spades":   "黑桃",
		"suit.hearts":   "红心",
		"suit.clubs":    "梅花",
		"suit.diamonds": "方块",

		// Ranks
		"rank.J": "J", "rank.Q": "Q", "rank.K": "K", "rank.A": "A",

		// Blinds
		"blind.small": "小盲注",
		"blind.big":   "大盲注",
		"blind.boss":  "Boss盲注",

		// Boss names
		"boss.the_club.name":       "棍棒",
		"boss.the_club.desc":       "所有梅花牌被削弱",
		"boss.the_goad.name":       "刺棒",
		"boss.the_goad.desc":       "所有黑桃牌被削弱",
		"boss.the_head.name":       "头颅",
		"boss.the_head.desc":       "所有红心牌被削弱",
		"boss.the_window.name":     "窗户",
		"boss.the_window.desc":     "所有方块牌被削弱",
		"boss.the_plant.name":      "植物",
		"boss.the_plant.desc":      "所有人脸牌被削弱",
		"boss.the_psychic.name":    "灵媒",
		"boss.the_psychic.desc":    "必须打出5张牌",
		"boss.the_manacle.name":    "镣铐",
		"boss.the_manacle.desc":    "手牌上限-1",
		"boss.the_flint.name":      "燧石",
		"boss.the_flint.desc":      "基础筹码和倍率减半",
		"boss.the_wall.name":       "墙壁",
		"boss.the_wall.desc":       "目标分×4",
		"boss.the_needle.name":     "针",
		"boss.the_needle.desc":     "只能出1手牌",
		"boss.the_water.name":      "水",
		"boss.the_water.desc":      "没有弃牌机会",
		"boss.the_mouth.name":      "嘴巴",
		"boss.the_mouth.desc":      "只能打同一种牌型",
		"boss.the_eye.name":        "眼睛",
		"boss.the_eye.desc":        "每手必须不同牌型",
		"boss.the_arm.name":        "手臂",
		"boss.the_arm.desc":        "打出的牌型降1级",
		"boss.the_tooth.name":      "牙齿",
		"boss.the_tooth.desc":      "每打出1张牌-$1",
		"boss.violet_vessel.name":  "紫罗兰之器",
		"boss.violet_vessel.desc":  "目标分×6",
		"boss.crimson_heart.name":  "深红之心",
		"boss.crimson_heart.desc":  "每手随机禁用1张小丑",
		"boss.cerulean_bell.name":  "蔚蓝之钟",
		"boss.cerulean_bell.desc":  "1张手牌被强制选中",

		// Joker names & descriptions
		"joker.joker.name":         "小丑",
		"joker.joker.desc":         "+4 倍率",
		"joker.greedy_joker.name":  "贪婪小丑",
		"joker.greedy_joker.desc":  "方块牌计分时 +3 倍率",
		"joker.jolly_joker.name":   "快乐小丑",
		"joker.jolly_joker.desc":   "包含对子时 +8 倍率",
		"joker.sly_joker.name":     "狡猾小丑",
		"joker.sly_joker.desc":     "包含对子时 +50 筹码",
		"joker.half_joker.name":    "半个小丑",
		"joker.half_joker.desc":    "打出≤3张牌时 +20 倍率",
		"joker.banner.name":        "旗帜",
		"joker.banner.desc":        "每剩余1次弃牌 +30 筹码",
		"joker.four_fingers.name":  "四指",
		"joker.four_fingers.desc":  "4张牌即可构成顺子和同花",
		"joker.blackboard.name":    "黑板",
		"joker.blackboard.desc":    "手中全黑色时 ×3 倍率",
		"joker.acrobat.name":       "杂技师",
		"joker.acrobat.desc":       "最后一手牌 ×3 倍率",
		"joker.the_duo.name":       "二重奏",
		"joker.the_duo.desc":       "包含对子时 ×2 倍率",

		// Rarity
		"rarity.common":    "普通",
		"rarity.uncommon":  "罕见",
		"rarity.rare":      "稀有",
		"rarity.legendary": "传奇",

		// UI labels
		"ui.ante":             "底注",
		"ui.score":            "得分",
		"ui.target":           "目标",
		"ui.money":            "金钱",
		"ui.hands":            "出牌",
		"ui.discards":         "弃牌",
		"ui.interest":         "利息",
		"ui.level":            "等级",
		"ui.chips":            "筹码",
		"ui.mult":             "倍率",
		"ui.round":            "回合",

		// Actions
		"action.play":         "出牌",
		"action.discard":      "弃牌",
		"action.sort":         "排序",
		"action.select":       "选中",
		"action.buy":          "购买",
		"action.sell":         "卖出",
		"action.reroll":       "刷新",
		"action.leave":        "离开",
		"action.skip":         "跳过",
		"action.challenge":    "挑战",
		"action.continue":     "继续",
		"action.start":        "开始游戏",
		"action.quit":         "退出",

		// Screens
		"title.name":          "小丑牌 BALATRO",
		"title.subtitle":      "终端版",
		"title.press_start":   "按 Enter 开始游戏",
		"gameover.win":        "恭喜通关！",
		"gameover.lose":       "游戏结束",
		"gameover.play_again": "按 Enter 重新开始，按 Q 退出",

		// Shop
		"shop.title":          "商店",
		"shop.sold_out":       "已售罄",
		"shop.booster":        "丑角补充包",
		"shop.reroll_cost":    "刷新费用",
		"shop.free":           "免费",

		// Status bar
		"status.empty_slot":   "空槽位 — 可在商店购买小丑牌",
		"status.cost":         "费用",
		"status.sell_value":   "售价",
		"status.rarity":       "稀有度",
		"status.chips_value":  "筹码值",
		"status.debuffed":     "已削弱 — 不提供筹码",

		// Skip rewards
		"skip.money":          "获得 $5",
		"skip.reroll":         "下次商店首次刷新免费",
		"skip.free_joker":     "下次商店1张免费小丑",
		"skip.double":         "金钱翻倍(上限$20)",
		"skip.hand_size":      "下回合+3手牌上限",

		// Settlement
		"settlement.base":     "基础奖励",
		"settlement.hand_bonus": "剩余手数奖励",
		"settlement.interest": "利息",
		"settlement.total":    "总计",
	})
}
