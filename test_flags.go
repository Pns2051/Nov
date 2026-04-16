package main

import (
    "flag"
    "fmt"
)

func main() {
    flag.Parse()
    fmt.Printf("Args: %v\n", flag.Args())
    if len(flag.Args()) > 0 {
        fmt.Printf("Arg0: %s\n", flag.Arg(0))
        fmt.Printf("Prefix match? %v\n", len(flag.Arg(0)) > 15 && flag.Arg(0)[:17] == "chrome-extension:")
    }
}
