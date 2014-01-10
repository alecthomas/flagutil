package flagutil

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/doc"
	"io"
	"os"
	"strings"
)

// ParseFlagsFromJSON parses values from a JSON file into a FlagSet. Keys in
// the JSON file that do not correspond to flags will result in an error.
func ParseFlagsFromJSON(filename string, flags *flag.FlagSet) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()
	decoder := json.NewDecoder(r)
	var data map[string]interface{}
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}
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

// PrettyFormatFlags formats standard Go flag FlagSets in a way that doesn't
// make your eyes bleed.
func PrettyFormatFlags(w io.Writer, flags *flag.FlagSet) {
	l := 0
	flags.VisitAll(func(flag *flag.Flag) {
		if len(flag.Name) > l {
			l = len(flag.Name)
		}
	})

	l += 9

	indent := strings.Repeat(" ", l)

	flags.VisitAll(func(flag *flag.Flag) {
		prefix := fmt.Sprintf("  %-*s", l-2, fmt.Sprintf("--%s=%s", flag.Name, flag.DefValue))
		buf := bytes.NewBuffer(nil)
		doc.ToText(buf, flag.Usage, "", "", 80-l)
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
