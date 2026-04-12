package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yinuoli/balatro-cli/engine"
	"github.com/yinuoli/balatro-cli/i18n"
	"github.com/yinuoli/balatro-cli/ui/components"
)

type ShopScreen struct {
	cursor int // 0=joker0, 1=joker1, 2=booster, 3=reroll, 4=leave
	errMsg string
}

const (
	shopJoker0 = iota
	shopJoker1
	shopBooster
	shopReroll
	shopLeave
	shopMaxCursor = shopLeave
)

func (s *ShopScreen) Update(msg tea.Msg, game *engine.GameState) tea.Cmd {
	s.errMsg = ""
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "left", "h":
			if s.cursor > 0 {
				s.cursor--
			}
		case "right", "l":
			if s.cursor < shopMaxCursor {
				s.cursor++
			}
		case "enter":
			s.handleAction(game)
		case "q", "ctrl+c":
			return tea.Quit
		}
	}
	return nil
}

func (s *ShopScreen) handleAction(game *engine.GameState) {
	shop := game.Shop
	if shop == nil {
		return
	}

	switch s.cursor {
	case shopJoker0, shopJoker1:
		idx := s.cursor
		_, ok := shop.BuyJoker(idx, &game.Money, &game.Jokers, game.MaxJokerSlots)
		if !ok {
			s.errMsg = "无法购买 (金钱不足或槽位已满)"
		}
	case shopBooster:
		pack := shop.GenerateBoosterPack()
		if game.Money < pack.Cost {
			s.errMsg = "金钱不足"
			return
		}
		game.Money -= pack.Cost
		// For simplicity, auto-pick first N jokers from pack
		for i := 0; i < pack.Choose && i < len(pack.Jokers) && len(game.Jokers) < game.MaxJokerSlots; i++ {
			j := engine.NewJokerInstance(pack.Jokers[i])
			game.Jokers = append(game.Jokers, j)
		}
	case shopReroll:
		ok := shop.Reroll(&game.Money)
		if !ok {
			s.errMsg = "金钱不足"
		}
	case shopLeave:
		game.LeaveShop()
		s.cursor = 0
	}
}

func (s *ShopScreen) View(game *engine.GameState, width, height int) string {
	var b strings.Builder

	if game.Shop == nil {
		return "No shop"
	}

	shop := game.Shop

	header := lipgloss.NewStyle().Bold(true).Foreground(components.ColorGold).
		Align(lipgloss.Center).Width(width).
		Render("🏪 " + i18n.T("shop.title"))
	b.WriteString("\n" + header + "\n\n")

	// Money
	moneyLine := lipgloss.NewStyle().Foreground(components.ColorGreen).Bold(true).
		Render(fmt.Sprintf("%s: $%d", i18n.T("ui.money"), game.Money))
	b.WriteString(moneyLine + "\n\n")

	// Joker slots
	var items []string
	for i := 0; i < 2; i++ {
		def := shop.JokerSlots[i]
		focused := s.cursor == i
		if def != nil {
			name := i18n.T(def.NameKey)
			cost := shop.JokerCosts[i]
			costStr := fmt.Sprintf("$%d", cost)
			if cost == 0 {
				costStr = i18n.T("shop.free")
			}
			style := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(0, 1).Width(14)
			if focused {
				style = style.BorderForeground(components.ColorFocused)
			} else {
				style = style.BorderForeground(components.ColorBorder)
			}
			content := fmt.Sprintf("%s\n%s\n%s", name, i18n.T(def.DescKey), costStr)
			items = append(items, style.Render(content))
		} else {
			style := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Padding(0, 1).Width(14).Foreground(components.ColorDim)
			if focused {
				style = style.BorderForeground(components.ColorFocused)
			} else {
				style = style.BorderForeground(components.ColorDim)
			}
			items = append(items, style.Render(i18n.T("shop.sold_out")))
		}
	}

	// Booster pack
	tier := []string{"普通", "大型", "超级"}[shop.BoosterTier]
	packCost := []int{4, 6, 8}[shop.BoosterTier]
	boosterStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).Width(14)
	if s.cursor == shopBooster {
		boosterStyle = boosterStyle.BorderForeground(components.ColorFocused)
	} else {
		boosterStyle = boosterStyle.BorderForeground(components.ColorPurple)
	}
	items = append(items, boosterStyle.Render(
		fmt.Sprintf("%s\n%s\n$%d", tier+i18n.T("shop.booster"), "", packCost)))

	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(width).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, items...)))
	b.WriteString("\n\n")

	// Reroll & Leave buttons
	rerollStyle := buttonStyle(s.cursor == shopReroll)
	rerollCost := shop.RerollCost
	rerollLabel := fmt.Sprintf("%s ($%d)", i18n.T("action.reroll"), rerollCost)
	if shop.FreeReroll {
		rerollLabel = i18n.T("action.reroll") + " (" + i18n.T("shop.free") + ")"
	}
	leaveStyle := buttonStyle(s.cursor == shopLeave)
	buttons := rerollStyle.Render(rerollLabel) + "    " + leaveStyle.Render(i18n.T("action.leave"))
	b.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(width).Render(buttons))

	// Sell hint
	if len(game.Jokers) > 0 {
		sellHint := lipgloss.NewStyle().Foreground(components.ColorDim).Align(lipgloss.Center).Width(width).
			Render("提示: 在出牌界面按 Tab 切到小丑栏后可卖出小丑")
		b.WriteString(sellHint)
	}

	// Error
	if s.errMsg != "" {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(components.ColorRed).Align(lipgloss.Center).Width(width).
			Render("⚠ " + s.errMsg))
	}

	return b.String()
}

func (s *ShopScreen) Focus() components.FocusTarget {
	if s.cursor <= shopJoker1 {
		return components.FocusTarget{Type: components.FocusShopJoker, Index: s.cursor}
	}
	return components.FocusTarget{Type: components.FocusNone}
}
