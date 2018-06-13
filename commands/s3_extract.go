// Copyright © 2016-present Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	gitlab "github.com/plouc/go-gitlab-client"
	helper "github.com/rande/gitlab-ci-helper"

	"io"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3ExtractCommand struct {
	Ui          cli.Ui
	Verbose     bool
	Ref         string
	RefName     string
	Job         string
	Project     string
	ExtractPath string
	TagMatcher  string

	// s3 settings
	AwsRegion   string
	AwsEndPoint string
	AwsProfile  string
	AwsBucket   string
}

func (c *S3ExtractCommand) Run(args []string) int {

	flags := flag.NewFlagSet("s3:uploads", flag.ContinueOnError)
	flags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	flags.BoolVar(&c.Verbose, "verbose", false, "")

	flags.StringVar(&c.Job, "job", helper.GetEnv("CI_JOB_NAME", os.Getenv("CI_BUILD_NAME")), "The job name")
	flags.StringVar(&c.Ref, "ref", helper.GetEnv("CI_COMMIT_SHA", os.Getenv("CI_BUILD_REF")), "The reference (sha1)")
	flags.StringVar(&c.RefName, "ref-name", helper.GetEnv("CI_COMMIT_REF_NAME", os.Getenv("CI_BUILD_REF_NAME")), "The reference name (tag or branch)")
	flags.StringVar(&c.Project, "project", os.Getenv("CI_PROJECT_ID"), "The project reference")
	flags.StringVar(&c.ExtractPath, "path", "./", "The project reference")
	flags.StringVar(&c.TagMatcher, "tag-matcher", "(v|)[0-9]{1,}\\.[0-9]{1,}\\.[0-9]{1,}(-[A-Za-z]*|)", "Regular expression to match tag (default: semver format)")

	flags.StringVar(&c.AwsRegion, "region", os.Getenv("AWS_REGION"), "The s3 region")
	flags.StringVar(&c.AwsEndPoint, "endpoint", os.Getenv("AWS_ENDPOINT"), "The s3 endpoint")
	flags.StringVar(&c.AwsProfile, "profile", os.Getenv("AWS_PROFILE"), "The aws credentials")
	flags.StringVar(&c.AwsBucket, "bucket", os.Getenv("AWS_BUCKET"), "The s3 bucket")

	config := helper.NewConfig()
	client := gitlab.NewGitlab(config.Gitlab.Host, config.Gitlab.ApiPath, config.Gitlab.Token)

	if err := flags.Parse(args); err != nil {
		return 1
	}

	project, err := helper.GetProject(c.Project, client)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", err.Error()))

		return 1
	}

	credentials, err := helper.GetAwsCredentials(c.Project)

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

	zipTarget := fmt.Sprintf("%s/%s_%s.zip", os.TempDir(), c.Ref, c.Job)

	f, _ := os.Create(zipTarget)
	defer f.Close()

	putObject := &s3.GetObjectInput{
		Bucket: aws.String(c.AwsBucket),
		Key:    aws.String(key),
	}

	c.Ui.Output(fmt.Sprintf("Retrieve zip file from s3://%s/%s", c.AwsBucket, key))

	objOutput, err := s3client.GetObject(putObject)

	if err != nil {
		c.Ui.Output(fmt.Sprintf("Unable to download archive from s3://%s/%s, %s", c.AwsBucket, key, err))

		return 1
	}

	_, err = io.Copy(f, objOutput.Body)

	if err != nil {
		c.Ui.Output(fmt.Sprintf("Unable to retrieve zip file: %s, %s", zipTarget, err))

		return 1
	}

	c.Ui.Output(fmt.Sprintf("Extract zip to %s", c.ExtractPath))

	err = helper.Unzip(zipTarget, c.ExtractPath)

	if err != nil {
		c.Ui.Output(fmt.Sprintf("Unable to extract zip file: %s, %s", zipTarget, err))

		return 1
	}

	os.Remove(zipTarget)

	return 0
}

func (c *S3ExtractCommand) Synopsis() string {
	return "Send archive to a S3 bucket."
}

func (c *S3ExtractCommand) Help() string {
	helpText := `
Usage: gitlab-ci-helper s3:extract

  Extract archive from a S3 bucket

Options:

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
