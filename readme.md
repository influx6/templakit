Templatekit
============
TemplaKit code generates code from provided template files of the user.

## Install

```
go get -u github.com/gokit/sqlkit
```

## Usage

Running the following commands instantly generates all necessary files and packages for giving code gen.

```go
> templakit generate
```

## How It works

TemplaKit works on the idea that you know what you wish to generate, it exposes the structures created through [moz](github.com/influx6/moz), which 
lets you extract specific information from code structures to generate new ones.

The process follows:

1. Create a template source by refering either a giving inlined template or a file within directory using the `@source` annotation with a unique `id` and a optional generation criteria. The `@source` can be used either on a package or type level.

    Optional Generation Criteria:

        - `partial_test.go` - means to create filename from package name but as a test file where cli is called.

        - `partial.go` - means to create filename from package name where cli is called.

        - `go` - means to create a pure go file from template with any lteration.

- Inline Source

```go
@source(id => Mob, gen => partial.Go, {
    func Add(m {{sel TYPE1}}, n {{sel TYPE2}}) {{sel TYPE3}} {

    }
})
```

- Referenced Source

```go
@source(id => Mob, gen => partial.Go, file => "_source.tml")
```

2. Add `@makeFor` annotation with source `id` to be used for generation and provided attributes to pass in data.

```go
 @makeFor(id => Mob, filename => bob_gen.go, TYPE1 => int32, TYPE2 => int32, TYPE3 => int64)
```
