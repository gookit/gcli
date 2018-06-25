package main

import (
    "runtime"
    "flag"
    "fmt"
    "os"
)

//var ErrHelp = errors.New("flag: help requested")
var dir string
var q bool
var age int

type options struct {
    dir string
}

var op options

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    config()
}

// go build && ./consumer -d dfdf -d1=ddd -age 34 ddd
func config() {
    //flag.Usage = usage

    flag.StringVar(&dir, "d1", "", "search `directory` for include files")

    op = options{}

    flag.StringVar(&op.dir,"d", "", "search `directory` for include files")

    flag.BoolVar(&q, "q", false, "quit")
    flag.IntVar(&age, "age", 0, "your age")

    flag.Parse()

    fmt.Printf("os args %+v\n", os.Args)
    fmt.Printf("args=%+v, num=%d\n", flag.Args(), flag.NArg())

    fmt.Printf("%+v %s %d\n", op, dir, age)

    //flag.ErrHelp = ErrHelp
    //flag.ErrorHandling = flag.ExitOnError
}

