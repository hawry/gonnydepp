# GOnnyDep(p)

Yeah, I know - a really farfetched wordplay here.

Anyway, GOnnyDep(p) (here on GD for convenience) is a very simple local system binary dependency handler. Just specify different requirement groups and levels and GD takes care of the rest. If a system binary isn't available in PATH, an error will be thrown depending on the level of requirement.

**Disclaimer**: *This is not a dependency package manager or something like that - this is only because I'm lazy and doesn't have the energy to perform any exec.LookPath all the time in my own small projects*

## Levels

There are currently three different levels of requirements: `must`,`should` and `might`. Default is `should`.

**MUST** All binaries/executables in this group **must** exist in the PATH.

**SHOULD** At least one of the binaries/executables in this group **must** exist in the PATH.

**MIGHT** Only for convenience for the developer as a shorthand if a binary/executable can provide a "nice"-feature, and not a need-feature. This isn't really a required part, but it'll at least tell you the full path to these binaries if they exist.

## Install
`go get github.com/hawry/gonnydepp`

and
```golang
import (
  "github.com/hawry/gonnydepp"
)
```

## Usage
```golang

mustGroup:=gonnydepp.NewGroup()
mustGroup.Executable("rails")
mustGroup.Executable("vim")
mustGroup.Must()

shouldGroup:=gonnydepp.NewGroup()
shouldGroup.Executable("nano")
shouldGroup.Executable("git")
shouldGroup.Executable("thisprobablydoesntexist")

if err:=gonnydepp.Parse();err!=nil {
  //Dependencies not met for some reason. See err for which binary that is missing
} else {
  log.Printf("Existing commands: %s",gonnydepp.Existing()) //Returns an array of what binaries that exist in its shorthand form
  log.Printf("Missing commands: %s",gonnydepp.Missing()) //Returns an array of the missing binaries, in shorthand form
}

railsPath:=gonnydepp.Path("rails") //Returns the full path to the given binary

```
