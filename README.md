# Utilities for making life with Go's flags easier

Example:

```go
package main

import (
    "flag"
    "github.com/alecthomas/flagutil"
)

var (
    debugFlag    = flag.Bool("debug", false, "enable debug mode")
    logLevelFlag = flag.Int("log_level", 0, "set the log level for stdout logging")
)

func main() {
    flag.Usage = flagutil.MakeUsage("usage: test <flags>\n\nA test application.", "")
    flag.Parse()
}
```

Then calling `--help` results in the following:

```
$ test --help
usage: test <flags>

A test application.

  --debug=false   enable debug mode
  --log_level=0   set the log level for stdout logging
```


# API Documentation

```go
import "github.com/alecthomas/flagutil"
```


## Usage

#### func  MakeUsage

```go
func MakeUsage(prefix, postfix string) func()
```
MakeUsage creates a function that generates nicely formatted usage text, usable
as "flag.Usage".

#### func  ParseFlagsFromJSON

```go
func ParseFlagsFromJSON(r io.Reader, flags *flag.FlagSet) error
```
ParseFlagsFromJSON parses values from a JSON stream into a FlagSet. Keys in the
JSON file that do not correspond to flags will result in an error.

#### func  ParseFlagsFromMap

```go
func ParseFlagsFromMap(data map[string]interface{}, flags *flag.FlagSet) error
```
ParseFlagsFromMap loads flag values from a map[string]interface{} into a
FlagSet. Keys in the JSON file that do not correspond to flags will result in an
error.

#### func  PrettyFormatFlags

```go
func PrettyFormatFlags(w io.Writer, flags *flag.FlagSet)
```
PrettyFormatFlags formats standard Go flag FlagSets in a way that doesn't make
your eyes bleed.
