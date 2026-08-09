package main

import (
    "fmt"
    "net/http"
    "os"
)

type result struct {
    text string
    err error
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please provide a url")
        return
    }
    results := make(chan result)
    urls := os.Args[1:]
    for _, url := range urls {
        go func(url string) {
            resp, err := http.Get(url)
            if err != nil {
                results <- result{err: err}
                return
            }
            defer resp.Body.Close()
            results <- result{
                text: url + " " + resp.Status,
            }            
        }(url)
    }
    for range urls {
        res := <-results
        if res.err != nil {
            continue
        }
        fmt.Println(res.text)
        return
    }
    fmt.Println("All requests failed")
}
