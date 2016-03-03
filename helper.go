package gitlab_ci_helper

import (
	"errors"
	gitlab "github.com/plouc/go-gitlab-client"
	"strconv"
	"strings"
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

		for _, p := range projects {
			if p.Name == paths[1] && p.Namespace.Name == paths[0] {
				pId = int64(p.Id)

				break
			}
		}

		if pId == 0 {
			return nil, errors.New("Unable to find the project")
		}
	}

	project, err := client.Project(strconv.FormatInt(pId, 10))

	if err != nil {
		return nil, errors.New("Error: " + err.Error())
	}

	return project, err
}
