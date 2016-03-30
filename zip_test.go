// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gitlab_ci_helper

import (
	"archive/zip"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Zip_No_Error(t *testing.T) {

	targetPath := fmt.Sprintf("%s/test_zip.zip", os.TempDir())

	fmt.Print(targetPath)

	includePath := make(Paths, 0)
	includePath.Set("./README.md")

	excludePath := make(Paths, 0)
	excludePath.Set(".git")

	err := Zip(includePath, excludePath, targetPath)
	assert.NoError(t, err)

	r, err := zip.OpenReader(targetPath)

	defer r.Close()
	assert.NoError(t, err)

	count := 0
	for _, f := range r.File {
		count++
		assert.Equal(t, "README.md", f.Name)
	}

	assert.Equal(t, 1, count)
}
