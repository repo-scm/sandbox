# sandbox

[![Build Status](https://github.com/repo-scm/sandbox/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/repo-scm/sandbox/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/repo-scm/sandbox)](https://goreportcard.com/report/github.com/repo-scm/sandbox)
[![License](https://img.shields.io/github/license/repo-scm/sandbox.svg)](https://github.com/repo-scm/sandbox/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/repo-scm/sandbox.svg)](https://github.com/repo-scm/sandbox/tags)



## Introduction

git workspace sandbox



## Usage

```bash
sandbox serve [--address string]
```



## Settings

[sandbox](https://github.com/repo-scm/sandbox) parameters can be set in the file `sandbox/.env`.

An example of settings can be found in [sandbox/.env.example](https://github.com/repo-scm/sandbox/blob/main/sandbox/.env.example).

```
CUSTOM_USER=admin
PASSWORD=admin
PUID=1000
PGID=1000
TZ=UTC
SUBFOLDER=/
KEYBOARD=en-us-qwerty
```



## Screenshot

### Create

![create](create.png)

### Open

![open](open.png)



## License

Project License can be found [here](LICENSE).



## Reference

- [docker-webtop](https://github.com/linuxserver/docker-webtop)
