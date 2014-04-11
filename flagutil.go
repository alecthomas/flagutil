package flagutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/doc"
	"io"
	"os"
	"strings"

	"github.com/kr/pty"

	flag "github.com/alecthomas/pflag"
)

// Fatalf prints an error message and exits.
func Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

// UsageErrorf prints an error, the application usage string, then exits.
func UsageErrorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	flag.Usage()
	os.Exit(1)
}

// ParseFlagsFromJSON parses values from a JSON stream into a FlagSet. Keys in
// the JSON file that do not correspond to flags will result in an error.
func ParseFlagsFromJSON(r io.Reader, flags *flag.FlagSet) error {
	decoder := json.NewDecoder(r)
	var data map[string]interface{}
	err := decoder.Decode(&data)
	if err != nil {
		return err
	}
	return ParseFlagsFromMap(data, flags)
}

// ParseFlagsFromMap loads flag values from a map[string]interface{} into a FlagSet. Keys in
// the JSON file that do not correspond to flags will result in an error.
func ParseFlagsFromMap(data map[string]interface{}, flags *flag.FlagSet) error {
	for k, v := range data {
		f := flags.Lookup(k)
		if f == nil {
			return fmt.Errorf("unknown flag '%s'", k)
		}
		sv := fmt.Sprintf("%v", v)
		f.Value.Set(sv)
	}
	return nil
}

func formatFlag(flag *flag.Flag) string {
	flagString := ""
	if flag.Shorthand != "" {
		flagString += fmt.Sprintf("-%s, ", flag.Shorthand)
	}
	flagString += fmt.Sprintf("--%s=%s", flag.Name, flag.DefValue)
	return flagString
}

// PrettyFormatFlags formats standard Go flag FlagSets in a way that doesn't
// make your eyes bleed.
func PrettyFormatFlags(w io.Writer, flags *flag.FlagSet) {
	width := 80
	if t, ok := w.(*os.File); ok {
		if _, cols, err := pty.Getsize(t); err == nil {
			width = cols
		}
	}

	// Find flag column width.
	l := 0
	flags.VisitAll(func(flag *flag.Flag) {
		fl := len(formatFlag(flag))
		if fl > l {
			l = fl
		}
	})

	l += 3

	indent := strings.Repeat(" ", l)

	flags.VisitAll(func(flag *flag.Flag) {
		prefix := fmt.Sprintf("  %-*s", l-2, formatFlag(flag))
		buf := bytes.NewBuffer(nil)
		doc.ToText(buf, flag.Usage, "", "", width-l)
		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		fmt.Fprintf(w, "%s%s\n", prefix, lines[0])
		for _, line := range lines[1:] {
			fmt.Fprintf(w, "%s%s\n", indent, line)
		}
	})
}

// MakeUsage creates a function that generates nicely formatted usage text,
// usable as "flag.Usage".
func MakeUsage(prefix, postfix string) func() {
	return func() {
		fmt.Printf("%s\n\n", prefix)
		PrettyFormatFlags(os.Stdout, flag.CommandLine)
		if postfix != "" {
			fmt.Printf("\n%s\n", postfix)
		}
	}
}
