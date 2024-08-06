package log

import (
	"fmt"
	"github.com/tech4works/checker"
)

func Print(lvl level, tag, prefix string, msg ...any) {
	tagText := BuildTagText(tag)
	levelText := BuildLevelText(lvl)

	if checker.IsNotEmpty(prefix) {
		fmt.Printf("[%s] %s %s %s", tagText, levelText, prefix, fmt.Sprintln(msg...))
	} else {
		fmt.Printf("[%s] %s %s", tagText, levelText, fmt.Sprintln(msg...))
	}
}

func Printf(lvl level, tag, prefix, format string, msg ...any) {
	tagText := BuildTagText(tag)
	levelText := BuildLevelText(lvl)

	if checker.IsNotEmpty(prefix) {
		fmt.Printf("[%s] %s %s %s\n", tagText, levelText, prefix, fmt.Sprintf(format, msg...))
	} else {
		fmt.Printf("[%s] %s %s\n", tagText, levelText, fmt.Sprintf(format, msg...))
	}
}
