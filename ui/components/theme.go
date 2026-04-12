package components

import "github.com/charmbracelet/lipgloss"

// Card dimensions
const (
	CardWidth  = 9
	CardHeight = 7
	JokerWidth  = 10
	JokerHeight = 7
)

// Colors
var (
	ColorRed       = lipgloss.Color("#FF4444")
	ColorWhite     = lipgloss.Color("#FFFFFF")
	ColorDim       = lipgloss.Color("#666666")
	ColorGold      = lipgloss.Color("#FFD700")
	ColorCyan      = lipgloss.Color("#00DDDD")
	ColorGreen     = lipgloss.Color("#44FF44")
	ColorBlue      = lipgloss.Color("#4488FF")
	ColorPurple    = lipgloss.Color("#AA44FF")
	ColorBg        = lipgloss.Color("#1a1a2e")
	ColorSidebarBg = lipgloss.Color("#16213e")
	ColorStatusBg  = lipgloss.Color("#0f3460")
	ColorBorder    = lipgloss.Color("#555555")
	ColorSelected  = lipgloss.Color("#FFD700")
	ColorFocused   = lipgloss.Color("#00DDDD")
)

// Styles
var (
	SidebarStyle = lipgloss.NewStyle().
			Width(26).
			Padding(1, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderRight(true).
			BorderForeground(ColorBorder)

	StatusBarStyle = lipgloss.NewStyle().
			Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(ColorBorder).
			Foreground(ColorWhite)

	CardBorderNormal = lipgloss.RoundedBorder()

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGold).
			Align(lipgloss.Center)

	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorDim)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Bold(true)

	ScoreStyle = lipgloss.NewStyle().
			Foreground(ColorGold).
			Bold(true)

	HandTypeStyle = lipgloss.NewStyle().
			Foreground(ColorCyan).
			Bold(true)
)
