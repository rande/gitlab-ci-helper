package gitlab_ci_helper

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_S3_Upload_Paths(t *testing.T) {

	paths := make(Paths, 0)

	paths.Set("Hello")
	paths.Set("World")

	assert.Equal(t, 2, len(paths))
}

func Test_Semversion_Matcher(t *testing.T) {

	values := []struct {
		Value  string
		Expect bool
	}{
		{"v1.0.0", true},
		{"v1.0.0-DEV", true},
		{"v1.0.0-dev", true},
		{"1.0.0-dev", true},
		{"1.0.0", true},
		{"v1-a-1", false},
	}

	for _, v := range values {
		assert.Equal(t, v.Expect, SemVersion.Match([]byte(v.Value)))
	}
}

func Test_GetEnv(t *testing.T) {
	assert.Equal(t, "Default Value", GetEnv("FOOBAR", "Default Value"))

	os.Setenv("FOOBAR", "FOOBAR")
	defer os.Unsetenv("FOOBAR")

	assert.Equal(t, "FOOBAR", GetEnv("FOOBAR", "Default Value"))
}

func Test_Path(t *testing.T) {

	p := Paths{"bonjour", "la terre!"}

	assert.Equal(t, "[bonjour la terre!]", p.String())
}
