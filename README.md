CLOG
==

Count Lines of GitHub Organization Code: CLOG counts blank lines, comment
lines, and physical lines source code in many programming languages on GitHub
Organization.

[![travis](https://img.shields.io/travis/linyows/clog.svg?style=for-the-badge)][travis]
[![codecov](https://img.shields.io/codecov/c/github/linyows/clog.svg?style=for-the-badge)][codecov]
[![release](http://img.shields.io/github/release/linyows/clog.svg?style=for-the-badge)][release]
[![godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=for-the-badge)][godoc]
[![license](http://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge)][license]

[travis]: https://travis-ci.org/linyows/clog
[release]: https://github.com/linyows/clog/releases
[license]: https://github.com/linyows/clog/blob/master/LICENSE
[godoc]: http://godoc.org/github.com/linyows/clog
[codecov]: https://codecov.io/gh/linyows/clog

Installation
--

```sh
$ go get github.com/linyows/clog
```

Usage
--

```sh
$ clog orgname --name foobar --token xxxxxxxxxx
```

Contribution
------------

1. Fork ([https://github.com/linyows/clog/fork](https://github.com/linyows/clog/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

Author
--

[linyows](https://github.com/linyows)
