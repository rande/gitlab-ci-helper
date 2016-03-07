# gitlab-ci-helper

This tool provides a binary cli to execute common commands inside a gitlab's job.


## Commands 

A detailled commands list is available in the [commands.md](commands.md) file

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