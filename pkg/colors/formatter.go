package colors

import (
	"devchat/pkg"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/acarl005/stripansi"
	chromastyles "github.com/alecthomas/chroma/styles"
	"github.com/jwalton/gchalk"
	markdown "github.com/quackduck/go-term-markdown"
)

func NewFormatter() *Formatter {
	f := &Formatter{}

	f.Init()

	return f
}

type Formatter struct {
	chalk  *gchalk.Builder
	Colors struct {
		Green   *gchalk.Builder
		Red     *gchalk.Builder
		Cyan    *gchalk.Builder
		Magenta *gchalk.Builder
		Yellow  *gchalk.Builder
		Orange  *gchalk.Builder
		Blue    *gchalk.Builder
		White   *gchalk.Builder
	}
	Styles struct {
		Normal []*pkg.Style
		Secret []*pkg.Style
	}
}

func (f *Formatter) Init() {
	markdown.CurrentTheme = chromastyles.ParaisoDark

	f.chalk = gchalk.New(gchalk.ForceLevel(gchalk.LevelAnsi256))

	f.Colors.Green = f.Ansi256(1, 5, 1)
	f.Colors.Red = f.Ansi256(5, 1, 1)
	f.Colors.Cyan = f.Ansi256(1, 5, 5)
	f.Colors.Magenta = f.Ansi256(5, 1, 5)
	f.Colors.Yellow = f.Ansi256(5, 5, 1)
	f.Colors.Orange = f.Ansi256(5, 3, 0)
	f.Colors.Blue = f.Ansi256(0, 3, 5)
	f.Colors.White = f.Ansi256(5, 5, 5)

	f.Styles.Normal = []*pkg.Style{
		{White, f.buildStyle(f.Colors.White)},
		{Red, f.buildStyle(f.Colors.Red)},
		{Coral, f.buildStyle(f.Ansi256(5, 2, 2))},
		{Green, f.buildStyle(f.Colors.Green)},
		{Sky, f.buildStyle(f.Ansi256(3, 5, 5))},
		{Cyan, f.buildStyle(f.Colors.Cyan)},
		{Magenta, f.buildStyle(f.Colors.Magenta)},
		{Pink, f.buildStyle(f.Ansi256(5, 3, 4))},
		{Rose, f.buildStyle(f.Ansi256(5, 0, 2))},
		{Lavender, f.buildStyle(f.Ansi256(4, 2, 5))},
		{Fire, f.buildStyle(f.Ansi256(5, 2, 0))},
		{PastelGreen, f.buildStyle(f.Ansi256(0, 5, 3))},
		{Olive, f.buildStyle(f.Ansi256(4, 5, 1))},
		{Yellow, f.buildStyle(f.Colors.Yellow)},
		{Orange, f.buildStyle(f.Colors.Orange)},
		{Blue, f.buildStyle(f.Colors.Blue)},
	}

	f.Styles.Secret = []*pkg.Style{
		{"ukraine", f.buildStyle(f.chalk.WithHex("#005bbb").WithBgHex("#ffd500"))},
		{"easter", f.buildStyle(f.chalk.WithRGB(255, 51, 255).WithBgRGB(255, 255, 0))},
		{"baby", f.buildStyle(f.chalk.WithRGB(255, 51, 255).WithBgRGB(102, 102, 255))},
		{"hacker", f.buildStyle(f.chalk.WithRGB(0, 255, 0).WithBgRGB(0, 0, 0))},
		{"l33t", f.buildStyleNoStrip(f.chalk.WithBgBrightBlack())},
		{"whiten", f.buildStyleNoStrip(f.chalk.WithBgWhite())},
		{"trans", f.makeFlag([]string{"#55CDFC", "#F7A8B8", "#FFFFFF", "#F7A8B8", "#55CDFC"})},
		{"gay", f.makeFlag([]string{"#FF0018", "#FFA52C", "#FFFF41", "#008018", "#0000F9", "#86007D"})},
		{"lesbian", f.makeFlag([]string{"#D62E02", "#FD9855", "#FFFFFF", "#D161A2", "#A20160"})},
		{"bi", f.makeFlag([]string{"#D60270", "#D60270", "#9B4F96", "#0038A8", "#0038A8"})},
		{"ace", f.makeFlag([]string{"#333333", "#A4A4A4", "#FFFFFF", "#810081"})},
		{"pan", f.makeFlag([]string{"#FF1B8D", "#FFDA00", "#1BB3FF"})},
		{"enby", f.makeFlag([]string{"#FFF430", "#FFFFFF", "#9C59D1", "#000000"})},
		{"aro", f.makeFlag([]string{"#3AA63F", "#A8D47A", "#FFFFFF", "#AAAAAA", "#000000"})},
		{"genderfluid", f.makeFlag([]string{"#FE75A1", "#FFFFFF", "#BE18D6", "#333333", "#333EBC"})},
		{"agender", f.makeFlag([]string{"#333333", "#BCC5C6", "#FFFFFF", "#B5F582", "#FFFFFF", "#BCC5C6", "#333333"})},
		{"rainbow", func(a string) string {
			rainbow := []*gchalk.Builder{
				f.Colors.Red,
				f.Colors.Orange,
				f.Colors.Yellow,
				f.Colors.Green,
				f.Colors.Cyan,
				f.Colors.Blue,
				f.Ansi256(2, 2, 5),
				f.Colors.Magenta,
			}
			return f.ApplyRainbow(rainbow, a)
		}}}
}

func (f *Formatter) makeFlag(colors []string) func(a string) string {
	flag := make([]*gchalk.Builder, len(colors))
	for i := range colors {
		flag[i] = f.chalk.WithHex(colors[i])
	}
	return func(a string) string {
		return f.ApplyRainbow(flag, a)
	}
}

func (f *Formatter) ApplyRainbow(rainbow []*gchalk.Builder, a string) string {
	a = stripansi.Strip(a)
	buf := ""
	colorOffset := rand.Intn(len(rainbow))
	for i, r := range []rune(a) {
		buf += rainbow[(colorOffset+i)%len(rainbow)].Paint(string(r))
	}
	return buf
}

func (f *Formatter) buildStyle(c *gchalk.Builder) func(string) string {
	return func(s string) string {
		return c.Paint(stripansi.Strip(s))
	}
}

func (f *Formatter) buildStyleNoStrip(c *gchalk.Builder) func(string) string {
	return func(s string) string {
		return c.Paint(s)
	}
}

// with r, g and b values from 0 to 5
func (f *Formatter) Ansi256(r, g, b uint8) *gchalk.Builder {
	return f.chalk.WithRGB(255/5*r, 255/5*g, 255/5*b)
}

func (f *Formatter) BgAnsi256(r, g, b uint8) *gchalk.Builder {
	return f.chalk.WithBgRGB(255/5*r, 255/5*g, 255/5*b)
}

// Applies color from name
func (f *Formatter) ChangeColor(u *pkg.User, colorName string) error {
	style, err := f.GetStyle(colorName)
	if err != nil {
		return err
	}

	if strings.HasPrefix(colorName, "bg-") {
		u.Color.Background = style.Name // update bg color
	} else {
		u.Color.Foreground = style.Name // update fg color
	}

	u.Name, _ = f.ApplyColorToData(u.Name, u.Color.Foreground, u.Color.Background) // error can be discarded as it has already been checked earlier

	u.Term.SetPrompt(fmt.Sprintf("%s: ", u.Name))

	// TODO: having savebans here is wildly incoherent, but this was noticed during a refactor.
	// it stays until i determine something else to do with it.
	if err = u.Room.Server.SaveBans(); err != nil {
		return fmt.Errorf("could not save the bans file: %v", err)
	}

	return nil
}

func (f *Formatter) ApplyColorToData(data string, color string, colorBG string) (string, error) {
	styleFG, err := f.GetStyle(color)
	if err != nil {
		return "", err
	}

	styleBG, err := f.GetStyle(colorBG)
	if err != nil {
		return "", err
	}

	return styleBG.Apply(styleFG.Apply(data)), nil // fg clears the bg color
}

// Sets either the foreground or the background with a random color if the
// given name is correct.
func (f *Formatter) GetRandomColor(name string) *pkg.Style {
	var foreground bool

	if name == "random" {
		foreground = true
	} else if name == "bg-random" {
		foreground = false
	} else {
		return nil
	}

	r := rand.Intn(6)
	g := rand.Intn(6)
	b := rand.Intn(6)
	if foreground {
		return &pkg.Style{fmt.Sprintf("%03d", r*100+g*10+b), f.buildStyle(f.Ansi256(uint8(r), uint8(g), uint8(b)))}
	}

	return &pkg.Style{fmt.Sprintf("bg-%03d", r*100+g*10+b), f.buildStyleNoStrip(f.BgAnsi256(uint8(r), uint8(g), uint8(b)))}
}

// If the input is a named style, returns it. Otherwise, returns nil.
func (f *Formatter) GetNamedColor(name string) *pkg.Style {
	for _, s := range f.Styles.Normal {
		if s.Name == name {
			return s
		}
	}

	for _, s := range f.Styles.Secret {
		if s.Name == name {
			return s
		}
	}

	return nil
}

func (f *Formatter) GetCustomColor(name string) (*pkg.Style, error) {
	if strings.HasPrefix(name, "#") {
		return &pkg.Style{name, f.buildStyle(f.chalk.WithHex(name))}, nil
	}

	if strings.HasPrefix(name, "bg-#") {
		return &pkg.Style{name, f.buildStyleNoStrip(f.chalk.WithBgHex(strings.TrimPrefix(name, "bg-")))}, nil
	}

	if len(name) != 3 && len(name) != 6 {
		return nil, nil
	}

	rgbCode := name

	if strings.HasPrefix(name, "bg-") {
		rgbCode = strings.TrimPrefix(rgbCode, "bg-")
	}

	a, err := strconv.Atoi(rgbCode)
	if err != nil {
		return nil, err
	}

	r := (a / 100) % 10
	g := (a / 10) % 10
	b := a % 10

	if r > 5 || g > 5 || b > 5 || r < 0 || g < 0 || b < 0 {
		return nil, errors.New("custom colors have values from 0 to 5 smh")
	}

	if strings.HasPrefix(name, "bg-") {
		return &pkg.Style{name, f.buildStyleNoStrip(f.BgAnsi256(uint8(r), uint8(g), uint8(b)))}, nil
	}

	return &pkg.Style{name, f.buildStyle(f.Ansi256(uint8(r), uint8(g), uint8(b)))}, nil
}

// Turns name into a style (defaults to nil)
func (f *Formatter) GetStyle(name string) (*pkg.Style, error) {
	randomColor := f.GetRandomColor(name)
	if randomColor != nil {
		return randomColor, nil
	}
	if name == "bg-off" {
		return &pkg.Style{"bg-off", func(a string) string { return a }}, nil // Used to remove one's background
	}

	namedColor := f.GetNamedColor(name)
	if namedColor != nil {
		return namedColor, nil
	}
	if strings.HasPrefix(name, "#") {
		return &pkg.Style{name, f.buildStyle(f.chalk.WithHex(name))}, nil
	}

	customColor, err := f.GetCustomColor(name)
	if err != nil {
		return nil, err
	}
	if customColor != nil {
		return customColor, nil
	}
	return nil, errors.New("Which color? Choose from random, " + strings.Join(func() []string {
		colors := make([]string, 0, len(f.Styles.Normal))
		for i := range f.Styles.Normal {
			colors = append(colors, f.Styles.Normal[i].Name)
		}
		return colors
	}(), ", ") + "  \nMake your own colors using hex (#A0FFFF, etc) or RGB values from 0 to 5 (for example, `color 530`, a pretty nice Orange). Set bg color like this: color bg-530; remove bg color with color bg-off.\nThere's also a few secret colors :)")
}
