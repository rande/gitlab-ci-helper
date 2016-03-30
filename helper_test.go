package gitlab_ci_helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_S3_Upload_Paths(t *testing.T) {

	paths := make(Paths, 0)

	paths.Set("Hello")
	paths.Set("World")

	assert.Equal(t, 2, len(paths))
}
