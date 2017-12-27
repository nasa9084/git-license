git-license
===========

Get LICENSE from [GitHub Licenses API](https://developer.github.com/v3/licenses/).

## Installation

### go get

`go get github.com/nasa9084/git-license`

### download binary

Download from [releases page](https://github.com/nasa9084/git-license/releases) and put the binary into your PATH.

### build your own

requirements: [dep](https://github.com/golang/dep)

``` shell
$ git clone https://github.com/nasa9084/git-license.git
$ cd git-license
$ dep ensure
$ go build .
```

and copy or move `git-license` binary into your PATH.

## How to use

Basically:

``` shell
$ git license --username YOUR_USER_NAME LICENSE_NAME
```

such as: `git license --username nasa9084 mit`.

You can see the list of license names with:

``` shell
$ git license -l
```

### Show HELP

``` shell
$ git license --help
```
