// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gitlab_ci_helper

import (
	"errors"
	"github.com/rande/garchive"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Unzip(archive, target string) error {
	return garchive.ExtractZipFile(archive, target)
}

func Zip(includePaths, excludePaths Paths, target string) error {
	var excludes []*regexp.Regexp
	for _, path := range excludePaths {
		excludes = append(excludes, regexp.MustCompile(path))
	}

	files := []string{}

	for _, source := range includePaths {
		info, err := os.Stat(source)
		if err != nil {
			return nil
		}

		var baseDir string
		if info.IsDir() {
			baseDir = source
		}

		filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			for _, exclude := range excludes {
				if exclude.Match([]byte(path)) {
					return nil
				}
			}

			if baseDir != "" {
				path = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			files = append(files, path)

			return err
		})
	}

	if len(files) == 0 {
		return errors.New("No file to zip")
	}

	return garchive.CreateZipFile(target, files)
}
