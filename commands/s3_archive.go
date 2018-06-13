// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mitchellh/cli"
	gitlab "github.com/plouc/go-gitlab-client"
	helper "github.com/rande/gitlab-ci-helper"
)

type S3ArchiveCommand struct {
	Ui           cli.Ui
	Verbose      bool
	Ref          string
	RefName      string
	Job          string
	Project      string
	IncludePaths helper.Paths
	IgnorePaths  helper.Paths
	IgnoreCVS    bool
	TagMatcher   string

	// s3 settings
	AwsRegion   string
	AwsEndPoint string
	AwsProfile  string
	AwsBucket   string
}

func (c *S3ArchiveCommand) Run(args []string) int {

	flags := flag.NewFlagSet("s3:uploads", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	flags.BoolVar(&c.Verbose, "verbose", false, "")
	flags.StringVar(&c.Job, "job", helper.GetEnv("CI_JOB_NAME", os.Getenv("CI_BUILD_NAME")), "The job name")
	flags.StringVar(&c.Ref, "ref", helper.GetEnv("CI_COMMIT_SHA", os.Getenv("CI_BUILD_REF")), "The reference (sha1)")
	flags.StringVar(&c.RefName, "ref-name", helper.GetEnv("CI_COMMIT_REF_NAME", os.Getenv("CI_BUILD_REF_NAME")), "The reference name (tag or branch)")
	flags.StringVar(&c.Project, "project", os.Getenv("CI_PROJECT_ID"), "The project reference")
	flags.StringVar(&c.TagMatcher, "tag-matcher", "(v|)[0-9]{1,}\\.[0-9]{1,}\\.[0-9]{1,}(-[A-Za-z]*|)", "Regular expression to match tag (default: semver format)")

	flags.StringVar(&c.AwsRegion, "region", os.Getenv("AWS_REGION"), "The s3 region")
	flags.StringVar(&c.AwsEndPoint, "endpoint", os.Getenv("AWS_ENDPOINT"), "The s3 endpoint")
	flags.StringVar(&c.AwsProfile, "profile", helper.GetEnv("AWS_PROFILE", "default"), "The aws credentials")
	flags.StringVar(&c.AwsBucket, "bucket", os.Getenv("AWS_BUCKET"), "The s3 bucket")

	flags.BoolVar(&c.IgnoreCVS, "ignore-cvs", true, "Ignore CVS files")

	c.IgnorePaths = make(helper.Paths, 0)
	c.IncludePaths = make(helper.Paths, 0)

	config := helper.NewConfig()
	client := gitlab.NewGitlab(config.Gitlab.Host, config.Gitlab.ApiPath, config.Gitlab.Token)

	flags.Var(&c.IgnorePaths, "exclude", "-ignore path/to/ignore")
	flags.Var(&c.IncludePaths, "include", "-include path/to/ignore")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	project, err := helper.GetProject(c.Project, client)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

		return 1
	}

	if len(c.IncludePaths) == 0 {
		c.IncludePaths.Set("./") // add the current path
	}

	if c.IgnoreCVS {
		c.IgnorePaths.Set(".git/")
		c.IgnorePaths.Set(".svn/")
		c.IgnorePaths.Set(".hg/")
		c.IgnorePaths.Set(".bzr/")
	}

	zipTarget := fmt.Sprintf("%s/%s_%s.zip", os.TempDir(), c.Ref, c.Job)

	c.Ui.Output(fmt.Sprintf("Generate zip file: %s", zipTarget))

	if err := helper.Zip(c.IncludePaths, c.IgnorePaths, zipTarget); err != nil {
		c.Ui.Output(fmt.Sprintf("Unable to zip file: %s", err))

		return 1
	}

	credentials, err := helper.GetAwsCredentials(c.AwsProfile)

	if err != nil {
		c.Ui.Output(fmt.Sprintf("Unable to load credentials: %s", err))

		return 1
	}

	awsConfig := &aws.Config{
		Region:           aws.String(c.AwsRegion),
		Endpoint:         aws.String(c.AwsEndPoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials,
	}

	s3client := s3.New(session.New(), awsConfig)

	key := fmt.Sprintf("commits/%s/%s/%s_%s.zip", project.Namespace.Path, project.Path, c.Ref, c.Job)

	if regexp.MustCompile(c.TagMatcher).Match([]byte(c.RefName)) {
		key = fmt.Sprintf("releases/%s/%s/%s_%s.zip", project.Namespace.Path, project.Path, c.RefName, c.Job)
	}

	f, _ := os.Open(zipTarget)
	defer f.Close()

	putObject := &s3.PutObjectInput{
		Bucket:      aws.String(c.AwsBucket),
		Key:         aws.String(key),
		Body:        f,
		ContentType: aws.String("application/zip"),
	}

	count := 0
	for {
		count++

		c.Ui.Output(fmt.Sprintf("Copy zip file: s3://%s/%s (%d)", c.AwsBucket, key, count))

		_, err = s3client.PutObject(putObject)

		if err != nil {
			c.Ui.Output(fmt.Sprintf("Unable to copy zip file: %s, %s", zipTarget, err))

			if count > 5 {
				return 1
			}

			c.Ui.Output(fmt.Sprintf("Retry to copy file: %s", zipTarget))

			time.Sleep(2 * time.Second)

			continue
		}

		break
	}

	os.Remove(zipTarget)

	return 0
}

func (c *S3ArchiveCommand) Synopsis() string {
	return "Send archive to a S3 bucket."
}

func (c *S3ArchiveCommand) Help() string {
	helpText := `
Usage: gitlab-ci-helper s3:archive

  Send archive to a S3 bucket

Options:

  -include            Path to include (one option per path)
  -exclude            Path to exclude (one option per path)
  -ignore-cvs         Exclude CVS files: .git .svn .bzr .hg
  -verbose            Add verbose information to the output
  -job                The job name (default: 9.x: CI_JOB_NAME and 8.x: CI_BUILD_NAME)
  -ref                The reference (sha1) (default: 9.x: CI_COMMIT_SHA and 8.x: CI_BUILD_REF)
  -ref-name           The reference name (default: 9.x: CI_COMMIT_REF_NAME and 8.x: CI_BUILD_REF_NAME)
  -project            The project reference (default: CI_PROJECT_ID)
  -region             The s3 region (default: AWS_REGION)
  -endpoint           The s3 endpoint (default: AWS_ENDPOINT)
  -profile            The aws credentials name (default: AWS_PROFILE, if not set default)
  -bucket             The s3 bucket name (default: AWS_BUCKET)
  -tag-matcher        The regular expression to match a tag (default: semver)

Credentials are retrieved from environment:

  GITLAB_HOST         The gitlab host
  GITLAB_TOKEN        The user's token
  GITLAB_API_PATH     (optional) the api path, default to: "/api/v3"

`
	return strings.TrimSpace(helpText)
}
