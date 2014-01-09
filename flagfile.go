package flagfile

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
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
