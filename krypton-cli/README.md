## Overview
`krypton-cli` is a client to aid with integration and load testing of `krypton` services.

### Compatibility
Instructions and documentation below are assuming linux os. All tests below are
done on `ubuntu 22.04`.
Code for `krypton-cli` is written in `go` and should compile and run on windows.
See [How to build for windows](#how-to-build-for-windows)

### How to build for linux
```
make

$ ls -al bin/krypton-cli
-rwxrwxr-x 1 pgp pgp 8859648 Oct 28 03:45 bin/krypton-cli

$ file bin/krypton-cli
bin/krypton-cli: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, Go BuildID=REvo3JRnBlChaxijNA8M/tQSSW2_bmMFKOdIvKkUj/NcX2sIBfQCnkRQOmSnAi/1E54thrHNG3ou1oPTQb6, stripped
```

### How to build for windows
```
TARGET=windows make

$ ls -al bin/krypton-cli.exe
-rwxrwxr-x 1 pgp pgp 9136128 Oct 28 03:43 bin/krypton-cli.exe

$ file bin/krypton-cli.exe
bin/krypton-cli: PE32+ executable (console) x86-64, for MS Windows
```
If you are building on windows, you might be missing `make`, `docker` and other tools so substitute
accordingly to run the above commands.

Once built, you can copy the resultant binary to the desktop and the `config/config.yaml` to a folder named
`.krypton-cli` in your home directory. Please note that autocomplete will not work.

### How to install (linux only)
```
make install
```
The following changes are made to your system:
- folder `$HOME/.krypton-cli` is created
  - if folder already exists, relevant config file(s) are backed up in the same folder
- file `/usr/share/bash-completion/completions/krypton-cli` is created or updated
  - this is to allow auto-complete and requires `bash` shell
  - note: requires sudo if you are non root
- file `/usr/local/bin/krypton-cli` is created or updated
  - note: requires sudo if you are non root

### Override default configuration
You can override some of the default configurations.
- Config file (default is $HOME/.krypton-cli/config.yaml)
  - `CLI_CONFIG_FILE=/tmp/config.yaml krypton-cli`
  - This will load from `/tmp/config.yaml`
- Default profile (default is `aws`)
  - `CLI_PROFILE_NAME=local krypton-cli`
  - This will use `local` as the profile.
  - The profile you override must be defined in the current config.

Note: The above config override examples are linux specific. For windows,
set the environment variable accordingly.

### General command architecture
Commands support modules at the first level or as the first argument.
Each module will support a list of commands. Navigating this structure is as shown below

#### List modules
To list supported modules, invoke `krypton-cli` without arguments.
```
$ bin/krypton-cli
Please specify a module. Supported modules are:
- auth  (Acquire user tokens)
- dsts  (Device STS Service)
- es    (Enroll Service)
- fs    (File Service)
- fds   (File Distribution Service)
- iot   (interact with mqtt server)
- ss    (Scheduler Service)
- util  (Utility commands)
```

#### List commands in a module
To list commands in a module, use the module name as the first argument
```
$ bin/krypton-cli fs
Supported commands are:
- get_download_url
- get_file_details
- get_upload_url
- create_file
```

#### Global command flags
All commands support a base set of arguments that are globally applicable.
Each command can further support flags that are applicable to their functionality.

```
Usage of krypton-cli:
  -batch_size uint
        number of routines to run. parallel mode only. (default 100)
  -count uint
        how many times to run command (default 1)
  -doc string
        generate documentation (none, shell) (default "none")
  -parallel
        should the runs be in parallel
  -retry_count uint
        how many times to retry failures (default 10)
  -retry_wait uint
        number of seconds to wait between retries (default 5)
  -server string
        server (default "https://usdevms.adminx.hppipeline.com")
  -stdin
        use stdin for command input
  -verbose
        verbose logs (default false)
```

#### Getting help on a command
To find more info on any command, type the command name and `-help` as the argument.
```
$ bin/krypton-cli fs get_file_details -help
Usage of this:
  -batch_size int
        number of routines to run in parallel. (default 100)
  -count int
        number of iterations (default 1)
  -id int
        file id (default 1)
  -parallel
        run in parallel
  -server string
        server url (default "http://localhost:1234")
  -verbose
        show verbose logs
```
Notice that `get_file_details` has a `-id` param in addition to global arguments.

#### Get portable docs on a command
Currently supports curl commands with `-doc shell`. Since the commands are generated
from current params, it offers an opportunity to troubleshoot and create doc pages from
actual working commands.
```
$ bin/krypton-cli auth device_code -doc shell
2023/08/15 15:57:49 curl command:
curl -X POST -x http://web-proxy.austin.hpicorp.net:8080 -H "User-Agent: krypton-cli" -H "Content-Type: application/x-www-form-urlencoded"  -d "scope=openid%2Benroll" https://usdevms.adminx.hppipeline.com/services/oauth_handler/device/authorization
2
```
By default `-verbose` will use a `-doc none` which will still log an http request.
`-doc shell` will log as shown above with a well formed `curl` command. This allows
debugging and tracking a header/param option quickly without spending time reading documentation.

#### Chaining commands
Command results are in `json` format. This helps combining multiple commands with minimal external tools
to create useful workflows.

See external tools needed below
- [jq](https://stedolan.github.io/jq/)

Note that our ado runner machines have `jq` installed. `ci` runs can be supported
without extra tool installation.

For specific examples, see modules
- [es](es/README.md)
- [fs](fs/README.md)
- [fds](fds/README.md)

### Packaging
`make docker` will create a docker image. All `ci` tests use `krypton-cli` from a container
so name resolutions work without issues within a test environment. Pre-packaged images from `ci`
runs are available with `docker pull ghcr.io/hpinc/krypton/krypton-cli`

### Making a github release
`make release` and use the resultant `krypton-cli.0.0.1.tar.gz` to upload as release artifact.
To update version, update `SEMVER` in `tools/make_release.sh`
