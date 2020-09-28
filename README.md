### Notice Me

Notice Me is a static analysis tool for Go programs that can identify unused
values, typically function call results. This is similar to the `errcheck`
tool, but instead of checking for unused results of type `error`, it checks
for unused results of whatever types you specify. There's also some overlap
with the `unusedresult` example analysis pass.

### Building

To build, run `go build` in `cmd/noticeme`. The `noticeme` tool should work
with any reasonably-recent version of Go. It relies on the `golang.org/x/tools`
package, but shouldn't care about the specific version.

### Usage

To check the package in the current directory for any and all unused
values of type `error`:

        noticeme -types error .

Note that this will be pickier than other utilities -- for instance,
`fmt.Printf` returns an error which is almost never checked. There's no
hooks for disabling these warnings case-by-case right now. To specify
additional types, use a comma-separated list. So for instance, the motive
for developing this involved the usage of `*Container` objects in our
Roaring bitmap implementation:

        noticeme -types "*roaring.Container" .

Type names are checked with full import path qualifiers, package qualifiers
only, and unqualified. Note that these qualifiers are applied to all the
types in a line at once, so a type like this won't work:

        map[package.Foo]github.com/user/repo.Bar

### Limitations

The design here is to assume that the type of an expression which exists
as a standalone "expression statement" has not been "used". If you assign
a value, or use it in another expression, then it's being used. You may be
able to bypass this analysis. However, it does handle fairly arbitrary
expressions, not just direct function or method calls.

There is no guessing about pointers and non-pointers; if you specify a
type of `foo` as important, functions returing `*foo` are not special, and
vise versa.

Nested types are not recognized, but tuples (such as functions with multiple
return values) are handled.

### References

* [errcheck](https://github.com/kisielk/errcheck)
* [unusedresult](https://godoc.org/golang.org/x/tools/go/analysis/passes/unusedresult)
