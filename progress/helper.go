package progress

import (
	"fmt"
	"math/rand"
	"time"
)

var timeFormats = [][]int{
	{0},
	{1},
	{2, 1},
	{60},
	{120, 60},
	{3600},
	{7200, 3600},
	{86400},
	{172800, 86400},
}

var timeMessages = []string{
	"< 1 sec", "1 sec", "secs", "1 min", "mins", "1 hr", "hrs", "1 day", "days",
}

// HowLongAgo format a seconds, get how lang ago
func HowLongAgo(sec int64) string {
	intVal := int(sec)
	length := len(timeFormats)

	for i, item := range timeFormats {
		if intVal >= item[0] {
			ni := i + 1
			match := false

			if ni < length { // next exists
				next := timeFormats[ni]
				if intVal < next[0] { // current <= intVal < next
					match = true
				}
			} else if ni == length { // current is last
				match = true
			}

			if match { // match success
				if len(item) == 1 {
					return timeMessages[i]
				}

				// len is 2
				return fmt.Sprintf("%d %s", intVal/item[1], timeMessages[i])
			}
		}
	}

	return "unknown" // He should never happen
}

// format bytes number friendly
func formatMemoryVal(bytes uint64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%dB", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", float64(bytes)/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", float64(bytes)/1024/1024/1024)
	}
}

func repeatRune(char rune, length int) (chars []rune) {
	for i := 0; i < length; i++ {
		chars = append(chars, char)
	}

	return
}

// CharThemes collection. can use for Progress bar, RoundTripSpinner
var CharThemes = []rune{
	CharEqual,
	CharCenter,
	CharSquare,
	CharSquare1,
	CharSquare2,
}

// GetCharTheme by index number
func GetCharTheme(index int) rune {
	if len(CharThemes) > index {
		return CharThemes[index]
	}

	return RandomCharTheme()
}

// RandomCharTheme get
func RandomCharTheme() rune {
	rand.Seed(time.Now().UnixNano())
	return CharThemes[rand.Intn(len(CharsThemes)-1)]
}

// CharsThemes collection. can use for LoadingBar, LoadingSpinner
var CharsThemes = [][]rune{
	{'卍', '卐'},
	{'☺', '☻'},
	{'░', '▒', '▓'},
	{'⊘', '⊖', '⊕', '⊗'},
	{'◐', '◒', '◓', '◑'},
	{'✣', '✤', '✥', '❉'},
	{'-', '\\', '|', '/'},
	{'▢', '■', '▢', '■'},
	[]rune("▖▘▝▗"),
	[]rune("◢◣◤◥"),
	[]rune("⌞⌟⌝⌜"),
	[]rune("◎●◯◌○⊙"),
	[]rune("◡◡⊙⊙◠◠"),
	[]rune("⇦⇧⇨⇩"),
	[]rune("✳✴✵✶✷✸✹"),
	[]rune("←↖↑↗→↘↓↙"),
	[]rune("➩➪➫➬➭➮➯➱"),
	[]rune("①②③④"),
	[]rune("㊎㊍㊌㊋㊏"),
	[]rune("⣾⣽⣻⢿⡿⣟⣯⣷"),
	[]rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"),
	[]rune("▉▊▋▌▍▎▏▎▍▌▋▊▉"),
	[]rune("🌍🌎🌏"),
	[]rune("☰☱☲☳☴☵☶☷"),
	[]rune("⠋⠙⠚⠒⠂⠂⠒⠲⠴⠦⠖⠒⠐⠐⠒⠓⠋"),
	[]rune("🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚🕛"),
}

// GetCharsTheme by index number
func GetCharsTheme(index int) []rune {
	if len(CharsThemes) > index {
		return CharsThemes[index]
	}

	return RandomCharsTheme()
}

// RandomCharsTheme get
func RandomCharsTheme() []rune {
	rand.Seed(time.Now().UnixNano())
	return CharsThemes[rand.Intn(len(CharsThemes)-1)]
}
