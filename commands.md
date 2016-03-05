## Commands

### ci:meta

    Usage: gitlab-helper ci:meta [options]
    
      Dump meta information about ci into a ci.json file
    
    Options:
    
      -verbose            Add verbose information to the output

### ci:revision

    Usage: gitlab-helper ci:revision [options]
    
      Dump a REVISION file with the current sha1
    
    Options:
    
      -verbose            Add verbose information to the output
    
    Env Variables:
    
      CI_BUILD_REF        Get the revision from this variable

### flowdock:message

    Usage: gitlab-helper flowdock:message [options] organisation flow message
    
      Build a flowdock thread from the current build. Information are retrieved from
      environment variables.
    
    Arguments:
      organisation        The organisation name
      flow                The flow reference
      message             The message to send
    
    Options:
      -token              The flow's token (default: env var FLOWDOCK_SOURCE_TOKEN)
      -verbose            Add verbose information to the output
    
    
    Configuration:
    
      Please note, the command use the new Flowdock API.
    
      1. Go to https://www.flowdock.com/oauth/applications
      2. Create a new Application
      3. Enter name, description and make sure "Short application" is selected.
      4. Once validated, Go to "Tools for testing" and create a new source.
      5. Press "Generate Source" and store the generated token for later use
         as the FLOWDOCK_SOURCE_TOKEN

### flowdock:status

    Usage: gitlab-helper flowdock:message [options] organisation flow
    
      Build a flowdock thread from the current build. Information are retrieved from
      environment variables.
    
      The external thread id is: gitlab:sha1
    
      You can use the -last option to indicate that the current job is the last one.
    
    Arguments:
      organisation        The organisation name
      flow                The flow reference
    
    Options:
      -ref                The commit related to the build (default: env var CI_BUILD_REF)
      -project            The project related to the build (default: env var CI_PROJECT_ID)
      -name               The build's name (default: env var CI_BUILD_NAME)
      -ref-name           The reference name (default: env var CI_BUILD_REF_NAME)
      -last               Indicate if the current build is the last one
      -token              The flow's token (default: env var FLOWDOCK_SOURCE_TOKEN)
      -verbose            Add verbose information to the output
    
    
    Configuration:
    
      Please note, the command use the new Flowdock API.
    
      1. Go to https://www.flowdock.com/oauth/applications
      2. Create a new Application
      3. Enter name, description and make sure "Short application" is selected.
      4. Once validated, Go to "Tools for testing" and create a new source.
      5. Press "Generate Source" and store the generated token for later use
         as the FLOWDOCK_SOURCE_TOKEN
    
    
    Gitlab's credentials are retrieved from environment:
    
      GITLAB_HOST         The gitlab host
      GITLAB_TOKEN        The user's token
      GITLAB_API_PATH     (optional) the api path, default to: "/api/v3"

### hipchat:message

    Usage: gitlab-helper ci:notify:hipchat [options] room message
    
      Send a message to one HipChat's room
    
    Arguments:
      room                The room reference
      message             The message to send
    
    Options:
      -color              The message color (default: gray, values: yellow, green, red, purple, gray, random)
      -notify             Whether this message should trigger a user notification (default: false)
      -token              The room's token (default: env var HIPCHAT_TOKEN)
      -server             The hipchat server, default to env var HIPCHAT_SERVER, then https://api.hipchat.com
      -verbose            Add verbose information to the output

### project:builds

    Usage: gitlab-helper project:builds:list [options] project
    
      List all builds available for the provide project
    
    Arguments:
    
      project             Can be an id or a string: namespace/name
    
    Options:
    
      -verbose            Add verbose information to the output
    
    Credentials are retrieved from environment:
    
      GITLAB_HOST         The gitlab host
      GITLAB_TOKEN        The user's token
      GITLAB_API_PATH     (optional) the api path, default to: "/api/v3"

### project:builds:artifacts

    Usage: gitlab-helper project:builds:artifacts [options] project build
    
      Download an artifacts and extract it if the 'path' option is provided
    
    Options:
    
      -build=XX           The build number used to retrieved the related artifact
      -stage=XX           The stage to search the build (must be used with -ref, default: package)
      -ref=XX             The sha1 linked to the build (must be used with -stage)
      -file=artifacts.zip The path to the artifact file (default: artifacts.zip)
      -path=./package     The path to extract the command. If not set, the artifact will not
                          be extracted.
      -verbose            Add verbose information to the output
    
    Credentials are retrieved from environment:
    
      GITLAB_HOST         The gitlab host
      GITLAB_TOKEN        The user's token
      GITLAB_API_PATH     (optional) the api path, default to: "/api/v3"

### project:list

    Usage: gitlab-helper project:list [options] project
    
      List all projects available
    
    Options:
    
      -verbose            Add verbose information to the output
    
    Credentials are retrieved from environment:
    
      GITLAB_HOST         The gitlab host
      GITLAB_TOKEN        The user's token
      GITLAB_API_PATH     (optional) the api path, default to: "/api/v3"
