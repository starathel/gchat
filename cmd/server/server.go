package main

import (
    "fmt"
    "google.golang.org/grpc"
)

func main() {
    fmt.Println("hello world,", grpc.Version)
}
