package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

// Group represents a nested logging group
type Group struct {
	prefix string
	indent int
}

func Info(format string, args ...any) Group {
	return Group{}.Info(format, args...)
}

func Property(key, msg string) Group {
	return Group{}.Property(key, msg)
}

func Progress(format string, args ...any) Group {
	return Group{}.Progress(format, args...)
}

func Error(format string, args ...any) Group {
	return Group{}.Error(format, args...)
}

func Warn(format string, args ...any) Group {
	return Group{}.Warn(format, args...)
}

func Success(format string, args ...any) Group {
	return Group{}.Success(format, args...)
}

func Spacer() {
	fmt.Fprintf(os.Stderr, "\n")
}

func (g Group) Info(format string, args ...any) Group {
	indentStr := strings.Repeat("  ", g.indent)
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s%s %s\n", indentStr, aurora.BrightBlue("✱"), msg)
	return Group{prefix: "i", indent: g.indent + 1}
}

func (g Group) Property(key, msg string) Group {
	indentStr := strings.Repeat("  ", g.indent)
	fmt.Fprintf(os.Stderr, "%s%s %s %s\n", indentStr, aurora.Cyan("✱"), aurora.Cyan(key), msg)
	return Group{prefix: "✱", indent: g.indent + 1}
}

func (g Group) Progress(format string, args ...any) Group {
	indentStr := strings.Repeat("  ", g.indent)
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s%s %s\n", indentStr, aurora.Cyan("✱"), msg)
	return Group{prefix: "✱", indent: g.indent + 1}
}

func (g Group) Error(format string, args ...any) Group {
	indentStr := strings.Repeat("  ", g.indent)
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s%s %s\n", indentStr, aurora.BrightRed("✗"), msg)
	return Group{prefix: "✗", indent: g.indent + 1}
}

func (g Group) Warn(format string, args ...any) Group {
	indentStr := strings.Repeat("  ", g.indent)
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s%s %s\n", indentStr, aurora.BrightYellow("!"), msg)
	return Group{prefix: "!", indent: g.indent + 1}
}

func (g Group) Success(format string, args ...any) Group {
	indentStr := strings.Repeat("  ", g.indent)
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s%s %s\n", indentStr, aurora.BrightGreen("✓"), msg)
	return Group{prefix: "✓", indent: g.indent + 1}
}
