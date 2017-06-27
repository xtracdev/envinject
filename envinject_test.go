package envinject

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

func setNoParamStore() {
	os.Setenv(ParamPrefixEnvVar,"")
}

func setEnv() {
	setNoParamStore()
	os.Setenv("foo", "foo val")
	os.Setenv("bar", "bar val")
}

func TestPassthroughGetenv(t *testing.T) {
	setEnv()

	i,err := NewInjectedEnv()

	if assert.Nil(t,err) {
		assert.Equal(t, "foo val", i.Getenv("foo"))
		assert.Equal(t, "bar val", i.Getenv("bar"))
	}
}

func TestPassthroughLookupEnv(t *testing.T) {
	setEnv()

	i,err := NewInjectedEnv()

	if assert.Nil(t,err) {
		foo,hasFoo := i.LookupEnv("foo")
		assert.Equal(t, "foo val", foo)
		assert.Equal(t, true, hasFoo)

		baz, hasBaz := i.LookupEnv("baz")
		assert.Equal(t, "", baz)
		assert.Equal(t, false, hasBaz)
	}
}

func sliceContains(slice []string, s string) bool {
	for _,v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

func TestPassthroughEnvironment(t *testing.T) {
	setEnv()
	i,err := NewInjectedEnv()
	if assert.Nil(t,err) {
		e := i.Environ()
		assert.True(t, sliceContains(e, "foo=foo val"))
		assert.True(t, sliceContains(e, "bar=bar val"))
	}

}
