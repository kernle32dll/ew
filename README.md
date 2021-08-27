[![Build Status](https://travis-ci.com/kernle32dll/ew.svg?branch=master)](https://travis-ci.com/kernle32dll/ew)
[![GoDoc](https://godoc.org/github.com/kernle32dll/ew?status.svg)](http://godoc.org/github.com/kernle32dll/ew)
[![Go Report Card](https://goreportcard.com/badge/github.com/kernle32dll/ew)](https://goreportcard.com/report/github.com/kernle32dll/ew)
[![codecov](https://codecov.io/gh/kernle32dll/ew/branch/master/graph/badge.svg)](https://codecov.io/gh/kernle32dll/ew)

# ew

ew - short for `(run things) e(very)w(here)` is a tool for grouping folders by tags,
and executing tasks in all folders via these tags.

## Getting ew

Grab the latest [release](https://github.com/kernle32dll/ew/releases), and put it in a executable path.
If you want to get hands on, install Go, and check out the latest `master` branch.

**Linux**

The following script downloads the latest release to /usr/local/bin/`, and makes it executable.

```shell script
sudo wget -c https://github.com/kernle32dll/ew/releases/download/v1.0.1/ew-linux-amd64 -O /usr/local/bin/ew
sudo chmod +x /usr/local/bin/ew
```

## Quickstart

Verify ew is working by executing `ew version`.

Now, you can either jump to any folder and get started tagging by calling `ew tags add @mytag`, or edit
the ew config file by hand (see #Config below)

Now you can run arbitrary commands for this tag (and so in every configured folder):

```shell script
ew @mytag cat go.mod
```

## How to use

COMMAND                           | DESCRIPTION
-------                           | ----
ew                                |   list all paths, grouped by their tags
ew help                           |   displays this help
ew --help                         |   displays this help (alias for ew help)
ew version                        |   displays version and build time (if applicable)
ew --version                      |   displays version and build time (if applicable) (alias for ew version)
ew migrate                        |   migrate from mixu/gr config, and keep json format
ew migrate --yaml                 |   migrate from mixu/gr config, and use new yaml format
ew paths                          |   list all paths (alias for ew paths list)
ew paths list                     |   list all paths
ew tags                           |   list all tags (alias for ew tags list)
ew tags list                      |   list all tags
ew tags add @some-tag             |   add current directory to tag "some-tag"
ew tags add \some\path @some-tag  |   add \some\path to tag "some-tag"
ew tags rm @some-tag              |   add current directory to tag "some-tag"
ew tags rm \some\path @some-tag   |   add \some\path to tag "some-tag"
ew status                         |   show quick git status for all paths
ew @tag1 status                   |   show quick git status for all paths of tag1 (supports multiple tags)
ew @tag1 some-cmd                 |   executes some-cmd in all paths of tag1 (supports multiple tags)

## Examples

```shell script
[bgerda@voidlight ~]$ ew @github status 
/home/bgerda/github-git/kernle32dll/emissione-go               master     Clean          @github     
/home/bgerda/github-git/kernle32dll/httpbulk-go                master     Clean          @github     
/home/bgerda/github-git/kernle32dll/jwtcache-go                master     Clean          @github     
/home/bgerda/github-git/kernle32dll/nullable                   master     Clean          @github     
/home/bgerda/github-git/kernle32dll/planlagt                   master     1 modified     @github     
/home/bgerda/github-git/kernle32dll/pooler-go                  master     Clean          @github     
/home/bgerda/github-git/kernle32dll/synchronized-cron-task     master     Clean          @github     
/home/bgerda/github-git/kernle32dll/turtleware                 master     Clean          @github
```

```shell script
[bgerda@voidlight ~]$ ew @github git pull

in /home/bgerda/github-git/kernle32dll/emissione-go:

Already up to date.

in /home/bgerda/github-git/kernle32dll/httpbulk-go:

Already up to date.

in /home/bgerda/github-git/kernle32dll/jwtcache-go:

Already up to date.

in /home/bgerda/github-git/kernle32dll/nullable:

Already up to date.

in /home/bgerda/github-git/kernle32dll/planlagt:

Already up to date.

in /home/bgerda/github-git/kernle32dll/pooler-go:

Already up to date.

in /home/bgerda/github-git/kernle32dll/synchronized-cron-task:

Already up to date.

in /home/bgerda/github-git/kernle32dll/turtleware:

Already up to date.

```

## Config

Per default, ew tries to loads its config from either `~/.ewconfig.yml` or `~/.ewconfig.json`. The default
is  `~/.ewconfig.yml`, if you are migrating from `gr` its `~/.ewconfig.json`.

The config file looks as follows:

`~/.ewconfig.yml`
```yaml
tags:
  server:
  - /home/bgerda/tc-git/turtlecoding/common/server-common
  github:
  - /home/bgerda/github-git/kernle32dll/emissione-go
  - /home/bgerda/github-git/kernle32dll/httpbulk-go
  - /home/bgerda/github-git/kernle32dll/jwtcache-go
  - /home/bgerda/github-git/kernle32dll/nullable
  - /home/bgerda/github-git/kernle32dll/planlagt
  - /home/bgerda/github-git/kernle32dll/pooler-go
  - /home/bgerda/github-git/kernle32dll/synchronized-cron-task
  - /home/bgerda/github-git/kernle32dll/turtleware
```

`~/.ewconfig.json`
```json
{
  "tags": {
    "server": [
      "/home/bgerda/tc-git/turtlecoding/common/server-common"
    ],
    "github": [
      "/home/bgerda/github-git/kernle32dll/emissione-go",
      "/home/bgerda/github-git/kernle32dll/httpbulk-go",
      "/home/bgerda/github-git/kernle32dll/jwtcache-go",
      "/home/bgerda/github-git/kernle32dll/nullable",
      "/home/bgerda/github-git/kernle32dll/planlagt",
      "/home/bgerda/github-git/kernle32dll/pooler-go",
      "/home/bgerda/github-git/kernle32dll/synchronized-cron-task",
      "/home/bgerda/github-git/kernle32dll/turtleware"
    ]
  }
}
```
