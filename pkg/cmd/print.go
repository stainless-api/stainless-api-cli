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

func SInfo(indent int, format string, args ...any) string {
	indentStr := strings.Repeat("  ", indent)
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s %s\n", indentStr, aurora.BrightBlue("✱"), msg)
}

func SProperty(indent int, key, msg string) string {
	indentStr := strings.Repeat("  ", indent)
	return fmt.Sprintf("%s%s %s %s\n", indentStr, aurora.Cyan("✱"), aurora.Cyan(key), msg)
}

func SProgress(indent int, format string, args ...any) string {
	indentStr := strings.Repeat("  ", indent)
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s %s\n", indentStr, aurora.Cyan("✱"), msg)
}

func SError(indent int, format string, args ...any) string {
	indentStr := strings.Repeat("  ", indent)
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s %s\n", indentStr, aurora.BrightRed("✗"), msg)
}

func SWarn(indent int, format string, args ...any) string {
	indentStr := strings.Repeat("  ", indent)
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s %s\n", indentStr, aurora.BrightYellow("!"), msg)
}

func SSuccess(indent int, format string, args ...any) string {
	indentStr := strings.Repeat("  ", indent)
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s %s\n", indentStr, aurora.BrightGreen("✓"), msg)
}

func (g Group) Info(format string, args ...any) Group {
	fmt.Fprint(os.Stderr, SInfo(g.indent, format, args...))
	return Group{prefix: "i", indent: g.indent + 1}
}

func (g Group) Property(key, msg string) Group {
	fmt.Fprint(os.Stderr, SProperty(g.indent, key, msg))
	return Group{prefix: "✱", indent: g.indent + 1}
}

func (g Group) Progress(format string, args ...any) Group {
	fmt.Fprint(os.Stderr, SProgress(g.indent, format, args...))
	return Group{prefix: "✱", indent: g.indent + 1}
}

func (g Group) Error(format string, args ...any) Group {
	fmt.Fprint(os.Stderr, SError(g.indent, format, args...))
	return Group{prefix: "✗", indent: g.indent + 1}
}

func (g Group) Warn(format string, args ...any) Group {
	fmt.Fprint(os.Stderr, SWarn(g.indent, format, args...))
	return Group{prefix: "!", indent: g.indent + 1}
}

func (g Group) Success(format string, args ...any) Group {
	fmt.Fprint(os.Stderr, SSuccess(g.indent, format, args...))
	return Group{prefix: "✓", indent: g.indent + 1}
}
