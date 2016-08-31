# izip-client

- [Installation](#installation)
- [Development](#development)

## Installation

Follow the instructions from the [releases page](https://github.com/clns/izip-client/releases).

## Development

You'll need a [Go dev environment](https://golang.org/doc/install) and [Docker](https://www.docker.com).

### Release

From the [releases page](https://github.com/clns/izip-client/releases) find out the next version number. Then, add the new version number to the [cmd/version.go](cmd/version.go) file and [build the binaries](#build-binaries).

On GitHub, edit the [latest release](https://github.com/clns/izip-client/releases) and change the tag name, the release name, the previous version for both download URLs and upload the new binaries.

### Build binaries

The following commands will build the binaries into the [build/](build) directory.

The first thing is to build the docker image used for cross compiling Windows and Linux binaries from Linux:

```sh
docker build -t izip .
```

#### Windows

:exclamation: Make sure the docker image is up to date.

```powershell
docker run --rm -v $env:GOPATH/src/github.com/clns/izip-client:/go/src/github.com/clns/izip-client izip bash -c 'CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o build/izip-Windows-x86_64.exe main.go'
```

#### Linux

:exclamation: Make sure the docker image is up to date.

```powershell
docker run --rm -v $env:GOPATH/src/github.com/clns/izip-client:/go/src/github.com/clns/izip-client izip bash -c 'CGO_ENABLED=1 go build -o build/izip-Linux-x86_64 main.go'
```

#### OS X

:exclamation: Make sure the docker image is up to date.

Since OS X binary can't be built from a different platform, you'll have to be on an OS X machine to run the following command:

```sh
CGO_ENABLED=1 go build -o build/izip-Darwin-x86_64 main.go
```