# FIGlet.go

```txt
___________.___  ________.__          __
\_   _____/|   |/  _____/|  |   _____/  |_      ____   ____
 |    __)  |   /   \  ___|  | _/ __ \   __\    / ___\ /  _ \
 |     \   |   \    \_\  \  |_\  ___/|  |     / /_/  >  <_> )
 \___  /   |___|\______  /____/\___  >__| /\  \___  / \____/
     \/                \/          \/     \/ /_____/
```

This project aims to fully implement the FIGfont spec in Go. It is a port of the JavaScript FIGdriver project, which can be found at [figlet.js](https://github.com/patorjk/figlet.js). The goal is to create a Go library that can generate ASCII art from text using FIGlet fonts while maintaining the same feature set as [figlet.js](https://github.com/patorjk/figlet.js), as well as a command-line tool for easy usage.

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
    text, err := figlet.Text("Hello, World!")
    if err != nil {
        panic(err)
    }

    fmt.Println(text)
}
```

This will print out:

```txt
  _   _      _ _        __        __         _     _ _ _
 | | | | ___| | | ___   \ \      / /__  _ __| | __| | | |
 | |_| |/ _ \ | |/ _ \   \ \ /\ / / _ \| '__| |/ _` | | |
 |  _  |  __/ | | (_) |   \ V  V / (_) | |  | | (_| |_|_|
 |_| |_|\___|_|_|\___/     \_/\_/ \___/|_|  |_|\__,_(_|_)
```

## Getting Started - Command Line

To use figlet.go on the command line, install globally:

```sh
go install github.com/damienbutt/figlet/cmd/figlet@latest
```

And then you should be able run from the command line. Example:

```sh
figlet "Hello, World!"
```

## LICENSE

[MIT](./LICENSE)
