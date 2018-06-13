package gitlab_ci_helper

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

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
