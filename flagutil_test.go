package flagutil

import (
	"os"
	"testing"
	"time"

	flag "github.com/ogier/pflag"

	"github.com/stretchrcom/testify/assert"
)

func TestParseFlagsFromJSON(t *testing.T) {
	flags := flag.NewFlagSet("test", flag.PanicOnError)
	testBool := flags.Bool("testbool", false, "test bool")
	testInt := flags.Int("testint", 0, "test int")
	testFloat := flags.Float64("testfloat", 0.0, "test float")
	testString := flags.String("teststring", "", "test string")
	testDuration := flags.Duration("testduration", 0, "test duration")
	r, err := os.Open("flagutil_test.json")
	assert.NoError(t, err)
	defer r.Close()
	assert.NoError(t, ParseFlagsFromJSON(r, flags))
	assert.True(t, *testBool)
	assert.Equal(t, *testInt, 99)
	assert.Equal(t, *testFloat, 99.9)
	assert.Equal(t, *testString, "a string")
	assert.Equal(t, *testDuration, time.Second*99)
}

func TestPrettyFormatFlags(t *testing.T) {
}
