# kache
A simple and a flexible in memory cache

[![Build Status](https://travis-ci.org/kasvith/kache.svg?branch=master)](https://travis-ci.org/kasvith/kache)
[![Build Status](https://cloud.drone.io/api/badges/kasvith/kache/status.svg)](https://cloud.drone.io/kasvith/kache)
[![Build status](https://ci.appveyor.com/api/projects/status/40cr0460vgqyyor8/branch/master?svg=true)](https://ci.appveyor.com/project/kasvith/kache/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kasvith/kache)](https://goreportcard.com/report/github.com/kasvith/kache)
[![HitCount](http://hits.dwyl.io/kasvith/kache.svg)](http://hits.dwyl.io/kasvith/kache)
[![codecov](https://codecov.io/gh/kasvith/kache/branch/master/graph/badge.svg)](https://codecov.io/gh/kasvith/kache)
[![GitHub](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/kasvith/kache/blob/master/LICENSE)

![gopher is looking at kache](https://user-images.githubusercontent.com/13379595/44355952-a3e7e480-a4cb-11e8-901f-aed77cfd63db.png)

# What is kache
**kache** aims to develop a *redis compatible in memory db* with [golang](https://golang.org/ "go"). Currently kache is powered up with **[RESP Protocol](https://redis.io/topics/protocol "RESP")**.
**kache** also supports simple text protocol so you can issue commands to kache using netcat or telnet as you please. kache has powered with many features managing a simple codebase with golang.

# Roadmap
- [x] Kache Server
- [x] Basic Commands as a POC
- [ ] Cluster Mode
- [ ] Pub/Sub Pattern
- [ ] Snapshots of data
- [ ] Kache CLI
- [ ] Client Libraries for popular languages
- [ ] Documentation
- [ ] Security
- [ ] Improved data Structures
- [ ] Website

# [Running kache](#command-line-opts)

kache is a compiled program, download the one for your platform and extract the package to a directory you wish.

Go to that directory, open a command prompt and run the kache executable like

- `./kache` if you are on **linux** or **mac**
- `.\kache` if you are on **windows**

This will start the application and port **7088** will be open by default.

Try to open **telnet** or **netcat** then
```
$: nc localhost 7088
ping
```

If you get the `+PONG` kache is working as expected.

Default configuration file can be found in `config/kache-default.toml`

kache can produce logs as you wish, in addition to default format it supports
 - json
 - logfmt

To run with a custom config file do

`./kache --config=path/to/config/file.toml`

### Synopsis

A fast and a flexible in memory database built with go

```
kache [flags]
```

### Options

```
      --config string    configuration file
  -d, --debug            output debug information
  -h, --help             help for kache
      --host string      host for running application (default "127.0.0.1")
      --logfile string   application log file
      --logging          set application logs (default true)
      --logtype string   kache can output logs in different formats like json or logfmt. The default one is custom to kache. (default "default")
      --maxClients int   max connections can be handled (default 10000)
      --maxTimeout int   max timeout for clients(in seconds) (default 120)
  -p, --port int         port for running application (default 7088)
  -v, --verbose          verbose output
```

# Development

## Prerequisites
 - Go 1.10.+

### Installing `mage`
`mage` is the build tool we use for build kache. To install `mage` 
 - Run `go get -u github.com/magefile/mage`
 - For a proper installation refer [official documentation](https://github.com/magefile/mage "official documentation")

## Setting up workspace
 - Fork the repo
 - Go to your **GOPATH** if you don't know about it learn from [here](https://github.com/golang/go/wiki/SettingGOPATH "here")
 - Create a directory github.com/kasvith
 - Clone the repo into that directory and cd to it

> Make sure you have an active internet connection as for the first time it will download some depedencies.

## Build the kache
 - `mage vendor` will install all the dependencies of the project(will take some time)
 - `mage kache` will produce the binary of the kache in `bin` directory
 - `mage kachecli` will produce the binary of the kache-cli in `bin` directory
 
### Other options
 - `mage check` will run `gofmt`, `goimports`, `go vet` and all tests with 32 bit platform including
 - `mage fmt` will run only `gofmt` on the code, will warn you when code has format errors
 - `mage vet` will reports suspicious constructs
 - `mage imports` will check import errors
 - `mage test` will run a unit test with defaults
 - `mage test386` will run a test in 32-bit mode
 - `mage testrace` will run a test with `race` conditions enabled
 - `mage -l` for list all commands

Special note : According to your environment executable will be built, for windows users it will need to add `.exe` to the end of `-o` flag like `go build -o bin/kache.exe ./cmd/kache`

# Contributions
**kache** is an **opensource** project. Contributions are welcome

- Fork the repo and star it :star:
- Open issues :boom:
- Raise PRs for issues :raised_hand:
- Help on documentation :page_facing_up:
- [Slack](https://join.slack.com/t/kache-db/shared_invite/enQtNDQ4NzYyNzI2NjQwLTMzNjNiMGIxYTQ1MDRiZjMxOTMwYzRiOTdkOTgyMThlYjM1MDlkZTVkN2Y5MmJjZmQyNGU2MDZlZWE2OTc3OWU "Slack")
