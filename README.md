# FIGlet.go

```txt
___________.___  ________.__          __
\_   _____/|   |/  _____/|  |   _____/  |_      ____   ____
 |    __)  |   /   \  ___|  | _/ __ \   __\    / ___\ /  _ \
 |     \   |   \    \_\  \  |_\  ___/|  |     / /_/  >  <_> )
 \___  /   |___|\______  /____/\___  >__| /\  \___  / \____/
     \/                \/          \/     \/ /_____/
```

[![CI](https://github.com/damienbutt/figlet/actions/workflows/ci.yml/badge.svg)](https://github.com/damienbutt/figlet/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/damienbutt/figlet)](https://goreportcard.com/report/github.com/damienbutt/figlet)
[![Go Reference](https://pkg.go.dev/badge/github.com/damienbutt/figlet.svg)](https://pkg.go.dev/github.com/damienbutt/figlet)
[![GitHub release](https://img.shields.io/github/v/release/damienbutt/figlet)](https://github.com/damienbutt/figlet/releases/latest)
[![License](https://img.shields.io/github/license/damienbutt/figlet)](./LICENSE)

A Go port of [figlet.js](https://github.com/patorjk/figlet.js) that implements the FIGfont spec to generate ASCII art from text. Provides both an importable Go library and a command-line tool.

## Quick Start

### Install

```sh
go get github.com/damienbutt/figlet@latest
```

### Simple Usage

```go
package main

import (
    "fmt"
    "github.com/damienbutt/figlet"
)

func main() {
    text, err := figlet.Text("Hello World!")
    if err != nil {
        panic(err)
    }

    fmt.Println(text)
}
```

This will print out:

```txt
  _   _      _ _        __        __         _     _ _
 | | | | ___| | | ___   \ \      / /__  _ __| | __| | |
 | |_| |/ _ \ | |/ _ \   \ \ /\ / / _ \| '__| |/ _` | |
 |  _  |  __/ | | (_) |   \ V  V / (_) | |  | | (_| |_|
 |_| |_|\___|_|_|\___/     \_/\_/ \___/|_|  |_|\__,_(_)
```

## Basic Usage

### Text

Generates ASCII art from the given text. Takes two parameters:

- **text** — the string to render as ASCII art.
- **opts** — an optional `*FigletOptions` to control rendering (font, layout, width, etc.). Omit or pass `nil` to use the package defaults.

Returns the generated ASCII art as a `string` and an `error`.

```go
package main

import (
    "fmt"
    "github.com/damienbutt/figlet"
)

func main() {
    text, err := figlet.Text("Boo!", &figlet.FigletOptions{
        Font:             "Ghost",
        HorizontalLayout: "default",
        VerticalLayout:   "default",
        Width:            80,
        WhitespaceBreak:  true,
    })
    if err != nil {
        fmt.Println("Something went wrong:", err)
        return
    }

    fmt.Println(text)
}
```

That will print out:

```txt
.-. .-')                            ,---.
\  ( OO )                           |   |
 ;-----.\  .-'),-----.  .-'),-----. |   |
 | .-.  | ( OO'  .-.  '( OO'  .-.  '|   |
 | '-' /_)/   |  | |  |/   |  | |  ||   |
 | .-. `. \_) |  |\|  |\_) |  |\|  ||  .'
 | |  \  |  \ |  | |  |  \ |  | |  |`--'
 | '--'  /   `'  '-'  '   `'  '-'  '.--.
 `------'      `-----'      `-----' '--'
```

### Options

The `FigletOptions` struct has the following fields:

#### Font

Type: `FontName` (string alias)
Default value: `"Standard"`

The FIGlet font to use for rendering.

#### HorizontalLayout

Type: `KerningMethods` (string alias)
Default value: `"default"`

Controls horizontal kerning/smushing. Accepted values:

| Value                   | Description                                                      |
| ----------------------- | ---------------------------------------------------------------- |
| `"default"`             | Kerning as intended by the font designer                         |
| `"full"`                | Full letter spacing, no overlap                                  |
| `"fitted"`              | Letters moved together until they almost touch                   |
| `"controlled smushing"` | Standard FIGlet controlled smushing rules                        |
| `"universal smushing"`  | Universal smushing — overlapping characters override one another |

#### VerticalLayout

Type: `KerningMethods` (string alias)
Default value: `"default"`

Controls vertical kerning/smushing. Accepts the same values as `HorizontalLayout`.

#### Width

Type: `int`
Default value: `0` (no limit)

Limits the output to this many characters wide. For example, if you want your output to be a max of 80 characters wide, you would set this option to `80`. Set to `0` to disable.

#### WhitespaceBreak

Type: `bool`
Default value: `false`

When used with `Width`, attempts to break text on whitespace boundaries rather than mid-word.

#### PrintDirection

Type: `PrintDirection`
Default value: `DefaultDirection` (uses the font's own setting)

Controls the rendering direction. Use `LeftToRight` or `RightToLeft` to override the font default.

#### ShowHardBlanks

Type: `bool`
Default value: `false`

When true, hardblank characters (used internally for spacing) are preserved in the output rather than replaced with spaces.

### Understanding Kerning

The 2 layout options allow you to override a font's default "kerning". Below you can see how this effects the text. The string "Kerning" was printed using the "Standard" font with horizontal layouts of "default", "fitted" and then "full".

```txt
  _  __               _
 | |/ /___ _ __ _ __ (_)_ __   __ _
 | ' // _ \ '__| '_ \| | '_ \ / _` |
 | . \  __/ |  | | | | | | | | (_| |
 |_|\_\___|_|  |_| |_|_|_| |_|\__, |
                              |___/
  _  __                   _
 | |/ / ___  _ __  _ __  (_) _ __    __ _
 | ' / / _ \| '__|| '_ \ | || '_ \  / _` |
 | . \|  __/| |   | | | || || | | || (_| |
 |_|\_\\___||_|   |_| |_||_||_| |_| \__, |
                                    |___/
  _  __                        _
 | |/ /   ___   _ __   _ __   (_)  _ __     __ _
 | ' /   / _ \ | '__| | '_ \  | | | '_ \   / _` |
 | . \  |  __/ | |    | | | | | | | | | | | (_| |
 |_|\_\  \___| |_|    |_| |_| |_| |_| |_|  \__, |
                                           |___/
```

In most cases you'll either use the default setting or the "fitted" setting. Most fonts don't support vertical kerning, but a hand full of them do (like the "Standard" font).

### Metadata

`Metadata` retrieves a font's parsed header options and comment string.

```go
package main

import (
    "fmt"
    "github.com/damienbutt/figlet"
)

func main() {
    meta, comment, err := figlet.Metadata("Standard")
    if err != nil {
        fmt.Println("something went wrong:", err)
        return
    }

    fmt.Printf("%+v\n", meta)
    fmt.Println(comment)
}
```

### Fonts

`Fonts` returns a sorted list of all available font names from the embedded font library.

```go
package main

import (
    "fmt"
    "github.com/damienbutt/figlet"
)

func main() {
    fonts, err := figlet.Fonts()
    if err != nil {
        fmt.Println("something went wrong:", err)
        return
    }

    fmt.Println(fonts)
}
```

### LoadedFonts

`LoadedFonts` returns the names of all fonts currently loaded in the in-memory cache (i.e. already parsed and ready to use).

```go
package main

import (
    "fmt"
    "github.com/damienbutt/figlet"
)

func main() {
    loaded := figlet.LoadedFonts()
    fmt.Println(loaded)
}
```

### ClearLoadedFonts

`ClearLoadedFonts` resets the in-memory font cache so that no fonts are loaded.

```go
package main

import (
    "fmt"
    "github.com/damienbutt/figlet"
)

func main() {
    figlet.ClearLoadedFonts()
}
```

### ParseFont

`ParseFont` allows you to load a font from any source — for example, from a file on disk — and register it under a given name for use with `Text`.

```go
package main

import (
    "fmt"
    "github.com/damienbutt/figlet"
)

func main() {
    data, err := os.ReadFile("myfont.flf")
    if err != nil {
        fmt.Println("something went wrong:", err)
        return
    }

    if _, err := figlet.ParseFont("myfont", string(data)); err != nil {
        fmt.Println("something went wrong:", err)
        return
    }

    text, err := figlet.Text("myfont!", &figlet.FigletOptions{Font: "myfont"})
    if err != nil {
        fmt.Println("something went wrong:", err)
        return
    }

    fmt.Println(text)
}
```

### Defaults

`Defaults` allows you to read or update the package-level defaults. Pass `nil` to read without modifying.

```go
package main

import (
    "fmt"
    "github.com/damienbutt/figlet"
)

func main() {
    // Update defaults
    figlet.Defaults(&figlet.FigletDefaults{
        Font:     "Standard",
        FontPath: "some-random-place/fonts",
    })

    // Read current defaults
    d := figlet.Defaults(nil)
    fmt.Println(d.Font)
}
```

## Getting Started - Command Line

To use figlet.go on the command line, install globally:

```sh
go install github.com/damienbutt/figlet/cmd/figlet@latest
```

And then you should be able run from the command line. Example:

```sh
figlet "Hello World!"
```

## LICENSE

[MIT](./LICENSE)
