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



## Deploy

### 1. Deploy on Railway

[![railway](https://railway.app/button.svg)](https://railway.app/template/HLP0Ub?referralCode=jch2ME)

### 2. Deploy to Render

[![render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/repo-scm/sandbox)

### 3. Deploy with Vercel

[![vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/repo-scm/sandbox&repository-name=sandbox)



## Screenshot

### Create

![create](create.png)

### Open

![open](open.png)



## License

Project License can be found [here](LICENSE).



## Reference

- [docker-webtop](https://github.com/linuxserver/docker-webtop)
