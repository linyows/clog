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
analyzed:  83/83
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
PHP                          54816        1665809        4324032        7188246
JavaScript                    5485         311133         518803        1508930
XML                           3465           1522            540        1477807
CSS                           2812         153656          65972         858077
Perl                          1947         132547          70862         600158
HTML                          2328          31458          22704         272195
Sass                           494          41512           9607         220136
C Header                       314          10661          10426          48123
SQL                            806           5299           7227          33666
Ruby                           801           4459           1507          25303
XSD                             50           1250            896          11454
C                              118           1774           1012          10059
Python                          35           1220           1379           9289
Terra                          242           2005             25           8184
LESS                           108           1250           1604           6286
Ant                             12            181            123            993
Objective-C++                    1            194              0            963
HCL                             19            159             76            898
XSLT                             6            176             46            800
Bourne Shell                   300           3175          21571            284
Go                               6             71              1            283
Markdown                      3473         140576         473944            175
BASH                           173           1719           8159            173
Mustache                        10              4              0             82
CoffeeScript                     3             13              0             75
Lua                              4             20              8             60
TypeScript                       6             10             72             44
PowerShell                       1              5              0             10
C Shell                          2             10             64              2
Plain Text                    1394          31040        1421353              2
M4                               8            909           8696              1
Awk                              2             15             92              0
Batch                           19            114            615              0
Gherkin                         23             66            416              0
JSON                         13972            346        2294682              0
Makefile                        30            468           1891              0
ReStructuredText               528          28650          75157              0
VimL                             2             11             97              0
YAML                           825           2573          41920              0
-------------------------------------------------------------------------------
TOTAL                        94640        2576060        9385579       12282758
-------------------------------------------------------------------------------
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
