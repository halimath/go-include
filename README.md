# include

![CI Status][ci-img-url] 
[![Go Report Card][go-report-card-img-url]][go-report-card-url] 
[![Package Doc][package-doc-img-url]][package-doc-url] 
[![Releases][release-img-url]][release-url]

`include` is a generator which includes file content into `string` or `[]byte` variables of go modules.
This generator is inspired by the Rust macro 
[`include_str`](https://doc.rust-lang.org/std/macro.include_str.html) which is used to include a file as a
rust variable.

Use this generator if you want to build a single binary that includes static assets such as html, scripts,
css or any other data as part of the produced binary.

## Usage

Install the generator:

```shell
$ go install github.com/halimath/include/cmd/include@latest
```

Add the library as a dependency:

```shell
$ go get github.com/halimath/include
```

Now, create one or more source files that contain package variable declarations which are initialized with a 
call to either `include.String` or `include.Bytes`. 

```go
import "github.com/halimath/include"

var html = include.String("./html/index.html")
```

**Only package level variable declarations are processed.**

Add a go build tag comment for `include` so that this file will not be included in a regular build. The 
generated go file will cary the inverse build tag. 

```go
//go:build include
```

You can customize the build tag with the `--build-tag` cli option.

Add a `go:generate` comment to instruct `go generate` how to generate a version of this file with all calls
to the above functions replaced with actual file content:

```go
//go:generate include --out file_gen.go $GOFILE
```

Run 

```shell
$ go generate <file>
```

to run the generator.

## Generator CLI options

The following options are supported by the generator:

Option | Default Value | Description
-- | -- | --
`out` | - | Name of the output file to write to. If not specified the generated source will be written to `stdout`.
`buildtag` | `include` | Build tag to deactivate in generate source.

## License

Copyright 2021 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[ci-img-url]: https://github.com/halimath/include/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/include
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/include
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/include
[release-img-url]: https://img.shields.io/github/v/release/halimath/include.svg
[release-url]: https://github.com/halimath/include/releases