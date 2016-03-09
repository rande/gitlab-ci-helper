// Copyright © 2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gitlab_ci_helper

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	gitlab "github.com/plouc/go-gitlab-client"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func GetProject(p string, client *gitlab.Gitlab) (*gitlab.Project, error) {
	pId, err := strconv.ParseInt(p, 10, 32)

	if err != nil {
		// invalid build id, search from a project list
		paths := strings.Split(p, "/")

		if len(paths) != 2 {
			return nil, errors.New("Error: Invalid project format, must be namespace/project-name")
		}

		projects, _ := client.Projects()

		try := ""
		for _, p := range projects {
			if p.Name == paths[1] {
				try = fmt.Sprintf("%s/%s", p.Namespace.Name, p.Name)
			}
			if p.Name == paths[1] && p.Namespace.Name == paths[0] {
				pId = int64(p.Id)

				break
			}
		}

		if pId == 0 {
			extra := ""
			if len(try) > 0 {
				extra = fmt.Sprintf("\nDid you mean: %s ?", try)
			}

			return nil, fmt.Errorf("Unable to find the project: %s/%s.%s", paths[0], paths[1], extra)
		}
	}

	project, err := client.Project(strconv.FormatInt(pId, 10))

	if err != nil {
		return nil, errors.New("Error: " + err.Error())
	}

	return project, err
}

func GetBuild(project *gitlab.Project, buildId string, client *gitlab.Gitlab) (*gitlab.Build, error) {
	build, err := client.ProjectBuild(strconv.FormatInt(int64(project.Id), 10), buildId)

	if err != nil {
		return nil, fmt.Errorf("Error: %s.\nUnable to find the build (projectId:%d, buildId:%s)", err.Error(), project.Id, buildId)
	}

	return build, err
}

// from http://blog.ralch.com/tutorial/golang-working-with-zip/
func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {

			if fileReader != nil {
				fileReader.Close()
			}

			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()

			if targetFile != nil {
				targetFile.Close()
			}

			return err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			fileReader.Close()
			targetFile.Close()

			return err
		}

		fileReader.Close()
		targetFile.Close()
	}

	return nil
}

func GetEnv(name, deflt string) string {
	if len(os.Getenv(name)) == 0 {
		return deflt
	}

	return os.Getenv(name)
}

type FakeRequest struct {
	Host   string
	Path   string
	Method string
	Called int

	Response *http.Response
}

func WrapperTestCommand(reqs []*FakeRequest, envs map[string]string, t *testing.T, fn func(ts *httptest.Server)) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Handling a request: %s\n", r.URL.Path)
		// dummy matcher
		for _, req := range reqs {
			if req.Path == r.URL.Path && req.Method == r.Method {
				fmt.Printf("Found request: %s\n", req.Path)

				req.Called++

				for name, values := range req.Response.Header {
					for _, v := range values {
						w.Header().Add(name, v)
					}
				}

				buf := bytes.NewBuffer([]byte(""))

				io.Copy(buf, req.Response.Body)

				//bytes.NewBuffer(buf.Bytes()).WriteTo(os.Stdout)
				bytes.NewBuffer(buf.Bytes()).WriteTo(w)

				req.Response.Body.Close()

				return
			}
		}

		t.Error("Unable to find a response to handle the request")
	}))

	envs["GITLAB_HOST"] = ts.URL

	defer func() {
		ts.Close()

		for name := range envs {
			os.Unsetenv(name)
		}
	}()

	for name, value := range envs {
		err := os.Setenv(name, value)

		if err != nil {
			panic(err)
		}
	}

	fn(ts)
}
