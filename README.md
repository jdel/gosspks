# gosspks: A Go SPK Server for Synology NAS 

[![Travis](https://img.shields.io/travis/jdel/gosspks.svg)](https://travis-ci.org/jdel/gosspks)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/jdel.org/gosspks)
[![GoReport](https://goreportcard.com/badge/github.com/jdel/go-syno)](https://goreportcard.com/report/jdel.org/gosspks)
[![Maintainability](https://img.shields.io/codeclimate/maintainability/jdel/gosspks.svg)](https://codeclimate.com/github/jdel/gosspks/maintainability)
[![Test Coverage](https://img.shields.io/codeclimate/coverage/jdel/gosspks.svg)](https://codeclimate.com/github/jdel/gosspks/test_coverage)
[![Dependencies](https://tidelift.com/badges/github/jdel/gosspks?style=flat)](https://tidelift.com/repo/github/jdel/gosspks)

`go get gopkg.in/jdel/gosspks.v0`

gosspks is the successor of [sspks](https://github.com/jdel/sspks).

It aims at providing full backwards compatibility as well as improvements, better performance and easier deployment.

## Installation

gosspks is provided as a single statically linked binary for darwin (x86, amd64), linux (x86, amd64, arm), windows (x86, amd64).

Only the linux-amd64 is extensively tested by myself.

Tagged releases are in the [Releases Page](https://jdel.org/gosspks/releases) while the latest build from master is avilable on [Bintray](https://bintray.com/jdel/gosspks/master/master#files).

To install, simply run:

```bash
curl -L https://raw.githubusercontent.com/jdel/gosspks/master/install.sh | sh
```

## Technical considerations

First off, I need to warn whoever that wants to dig in the code that I wrote `gosspks` to teach myself Go programming. Instead of thinking this from scratch, I mirrored the behavior of `sspks`, which results in experientations with Go, bad decisions taken early that I now cannot amend without a lot of rework. These will be dealt with in due time.

gosspks is databaseless. This is a deliberate choice to keep the design simple and avoid having a dependency and having to support multiple databases engines. In the world of stateless containers, this is not pretty, but I am aiming at providing support for object storage in the future.

In the meantime, in order to linit the IO overhead and increase response time, I implemented an in-memory cache mechanism. This works well, but the mechanism is complex and could be optimized.

The cache TTL and refresh can be tweaked but gosspks provides safe defaults (5 minutes TTL, 1 minute refresh). If you add or delete packages while gosspks is running, you will either have to wait for the cache to expire (maximum 5 minutes), or restart gosspks.

## Initial run

If it doesn't exist, gosspks will create a config file in `$HOME/gosspks/`.

In this case, the file will contain all default configuration values, overriden by any `--` flag or environment variable passed on the first run.

## Usage

```
Serving your Synology Packages.

Usage:
  gosspks [flags]
  gosspks [command]

Available Commands:
  config      Get the current running config
  help        Help about any command
  version     Get the version of sspks

Flags:
      --cache string                     cache directory (gosspks extracts INFO and images here) (default "cache")
  -C, --config string                    config file (default is $HOME/gosspks/config.yml)
  -d, --debug-package                    generates a debug package visible in Synology Package Center
      --download string                  prefix to serve packages (default "download")
  -h, --help                             help for gosspks
  -H, --home string                      gosspks home (default is $HOME/gosspks/
      --hostname string                  hostname to use when generating urls
  -l, --log-level string                 log level [Error,Warn,Info,Debug] (default "Error")
      --md5                              enable md5 calculation
      --models string                    models file (default "models.yml")
      --models-cache-duration string     models in-memory cache TTL (default "7d")
      --models-cache-refresh string      models in-memory cache automatic refresh rate (default "1d")
      --packages string                  packages directory (default "packages")
      --packages-cache-count int         im-memory cache size (0 to read packages from disk every time) (default 15)
      --packages-cache-duration string   packages in-memory cache TTL (default "5m")
      --packages-cache-refresh string    packages in-memory cache automatic refresh rate (default "1m")
  -p, --port int                         port to listen to (default 8080)
      --scheme string                    scheme to use when generating urls (default "http")
      --static string                    prefix to serve static images (default "static")

Use "gosspks [command] --help" for more information about a command.
```

Place your `.spk` files in `$HOME/gosspks/packages/` and run gosspks:

```bash
gosspks --port 80
```

A `systemd` unit file will be provided later on.

## Run with Docker

Official automated Docker build is available at [Docker Hub](https://hub.docker.com/r/jdel/gosspks/tags/).

Simply run a container with:

```bash
docker run -d --name gosspks \
           -p 80:8080 \
           -v $(pwd)/packages/:/home/user/gosspks/packages/:rw \
           -e GOSSPKS_HOSTNAME=yourdomain.com \
           jdel/gosspks:v0.1 
```

Is strictly equivalent to:

```bash
docker run -d --name gosspks \
           -p 80:8080 \
           -v $(pwd)/packages/:/home/user/gosspks/packages/:rw \
           jdel/gosspks:v0.1 gosspks --hostname yourdomain.com
```

It is not necessary to bind mount `/home/user/gosspks/gosspks.yml` when running with docker as options are enforced by environment variables or flags, but feel free to do so if you prefer to use a config file.

Now on your Docker host, place your `.spk` files in `$(pwd)/packages/`.

## Consume the API

### Examples

Get all packages:

```bash
curl -i localhost:80/v1/packages
```

The above command should return a 200 status code together with a json payload containing the packages list.

### API Routes

| Route                          | Method | Description                                                   |
| ------------------------------ | ------ | ------------------------------------------------------------- |
| /about                         | GET    | Shows gosspks version                                         |
| /v1/models                     | GET    | Returns all models (fetches list from the internet)           |
| /v1/packages/                  | GET    | Returns all packages                                          |
| /v1/packages/{synoPackageName} | GET    | Returns a specific package                                    |
| /                              | GET    | Endpoint for the Synology Package Center                      |
| /                              | POST   | Endpoint for the Synology Package Center (older DSM versions) |

There is also a special route that doesn't fit in the table:

`/getList/v0/{synoMajor}/{synoMinor}/{synoMicro}/{synoBuild}/{synoNano}/{synoArch}/{synoChannel}/{synoUnique}/{synoLanguage}`

This route mimics the behaviour of Synology's official package server. While you could hack your NAS to have gosspks impersonate the official Synology store, this has not been tested and is only here for educational purpose.

### Synology Package Center

Point your NAS package center to `yourdomain.com` to enjoy your packages.

## Secure with TLS

gosspks doesn't provide SSL termination. You may want to use `nginx` or `traefik` for this purpose.

If you chose to use https, you will need to pass `--scheme https` when running gosspks for it to generate the right links.

Some `nginx` and `traefik` configuration examples will be provided later on.

## Configuration

All configurables are listed by running `gosspks config --yaml`.

```
gosspks:
  cache:
    models:
      duration: 7d
      refresh: 1d
    packages:
      count: 15
      duration: 5m
      refresh: 1m
  debug-package: false
  filesystem:
    cache: cache
    models: models.yml
    packages: packages
  hostname: ""
  log-level: Error
  md5: false
  port: 8080
  router:
    download: /download
    static: /static
  scheme: http
```

It is possible to override any configuration options at runtime with flags or environment variables.

For example, `gosspks --filesystem-cache=my_cache` is equivalent to writing the following in the config file:

```
gosspks:
  filesystem:
    cache: my_cache
```

The same way, using environment variables like `GOSSPKS_FILESYSTEM_CACHE=my_cache gosspks` will have the same effect.

:warning: command line flags take precedence over environment variables.

I will not detail all configuration options but the most important 

## More configuration options

While gosspks will run fine with default options, it is possible to override the home directory by passing `--home /path/to/anywhere`.

It is possible to bypass completely the home directory by specifying the config file itself with `--config /etc/gosspks/gosspks.yml`.

## Ok, but, I want a web UI !

Right, there is [gosspks-ui](https://jdel.org/gosspks-ui), a React UI, but there is no real documentation yet, 

A Docker image is available, you can run the followinf `docker-compose.yml` file through `docker-compose up -d` or `docker stack deploy -c docker-compose.yml gosspks`:

```yaml
version: "3.1"
services:
  gosspks:
    image: jdel/gosspks:v0.1
    container_name: gosspks
    networks:
      - mobynet
    volumes:
      - $HOME/packages/:/home/user/gosspks/packages/:rw
    environment:
      - GOSSPKS_HOSTNAME=localhost
    logging:
      options:
        max-size: 50m

  gosspks-ui:
    image: jdel/gosspks-ui:v0.1
    container_name: gosspks-ui
    depends_on:
      - gosspks
    networks:
      - mobynet
    ports:
      - "80:80"
    environment:
      - GOSSPKS_UI_HOSTNAME=localhost
      - GOSSPKS_URL=http://gosspks:8080/
      - GOSSPKS_UI_TITLE=Your Package Server
      - GOSSPKS_UI_HEADER=Your Package Server
      - GOSSPKS_UI_SYNO_URL=http://localhost
    logging:
     options:
       max-size: 50m
       
networks:
  mobynet:
```