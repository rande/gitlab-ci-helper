// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gitlab_ci_helper

import (
	"archive/zip"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Zip_No_Error(t *testing.T) {
	targetPath := fmt.Sprintf("%s/gitlab_ci_helper_zip.zip", os.TempDir())

	includePath := make(Paths, 0)
	includePath.Set("README.md")

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

	os.Remove(targetPath)
}

func Test_Zip_File_Mode(t *testing.T) {

	targetPath := fmt.Sprintf("%s/gitlab_ci_helper_zip.zip", os.TempDir())
	binPath := "gitlab_ci_helper_zip.bin"

	os.Remove(binPath)
	os.Remove(fmt.Sprintf("%s/%s", os.TempDir(), binPath))

	f, err := os.OpenFile(binPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0755))

	assert.NoError(t, err)

	f.Write([]byte("some content"))
	f.Close()

	includePath := make(Paths, 0)
	includePath.Set(binPath)

	excludePath := make(Paths, 0)

	err = Zip(includePath, excludePath, targetPath)
	assert.NoError(t, err)

	r, err := zip.OpenReader(targetPath)
	assert.NoError(t, err)

	defer r.Close()
	assert.NoError(t, err)

	count := 0
	for _, f := range r.File {
		count++
		assert.Equal(t, "gitlab_ci_helper_zip.bin", f.Name)
		assert.Equal(t, os.FileMode(0755), f.Mode().Perm())
	}

	os.Remove(binPath)

	err = Unzip(targetPath, os.TempDir())
	assert.NoError(t, err)

	if f, err := os.Stat(fmt.Sprintf("%s/%s", os.TempDir(), binPath)); err == nil {
		assert.Equal(t, os.FileMode(0755), f.Mode().Perm())
	} else {
		assert.NoError(t, err, fmt.Sprintf("Error, no file extracted: %s", err))
	}

	os.Remove(fmt.Sprintf("%s/%s", os.TempDir(), binPath))
	os.Remove(binPath)
	os.Remove(targetPath)
}
