# gitlab-ci-helper

[![Build Status](https://travis-ci.org/rande/gitlab-ci-helper.png?branch=master)](https://travis-ci.org/rande/gitlab-ci-helper)
[![Coverage Status](https://coveralls.io/repos/github/rande/gitlab-ci-helper/badge.svg?branch=master)](https://coveralls.io/github/rande/gitlab-ci-helper?branch=master)
[![GoDoc](https://godoc.org/github.com/rande/gitlab-ci-helper?status.svg)](https://godoc.org/github.com/rande/gitlab-ci-helper)
[![GitHub license](https://img.shields.io/github/license/rande/gitlab-ci-helper.svg)](https://github.com/rande/gitlab-ci-helper/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rande/gitlab-ci-helper)](https://goreportcard.com/report/github.com/rande/gitlab-ci-helper)
[![GitHub issues](https://img.shields.io/github/issues/rande/gitlab-ci-helper.svg)](https://github.com/rande/gitlab-ci-helper/issues)

This tool provides a binary cli to execute common commands inside a gitlab's job.

## Installation

**gitlab-ci-helper** is a single binary with no external dependencies, released for several platforms.
Go to the [releases page](https://github.com/rande/gitlab-ci-helper/releases),
download the package for your OS, and copy the binary to somewhere on your PATH.
Please make sure to rename the binary to `gitlab-ci-helper` and make it executable.

As it primarily targets GitLab CI environment, a simple curl will also do :).

A release also includes checksums for various build flavors,
as a build has access to sensible information (GitLab tokens, AWS credentials…),
please use them.

## Commands 

A detailed commands list is available in the [commands.md](commands.md) file

## Build Commands
   
- ``ci:revision``: dump a REVISION file
- ``ci:meta``: dump a ci.json file, with build information 
- ``project:builds:artifacts``: download an artifacts file from a previous job


## Integration Commands

- ``hipchat:message``: send a message to hipchat
- ``flowdock:message``: send a message to flowdock
- ``flowdock:status``: create a build status on flowdock 

## Tools commands

- ``project:builds``: list builds
- ``project:list``: list projets