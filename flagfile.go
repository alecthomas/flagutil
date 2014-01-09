package flagfile

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type FlagFile map[string]interface{}

func ParseFlagsFromJSON(filename string, flags *flag.FlagSet) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()
	decoder := json.NewDecoder(r)
	var data FlagFile
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}
	for k, v := range data {
		f := flags.Lookup(k)
		if f == nil {
			continue
		}
		sv := fmt.Sprintf("%v", v)
		f.Value.Set(sv)
	}
	return nil
}
