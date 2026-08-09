package main

import (
    "flag"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

type result struct {
    url string
    text string
    err error
}

func main() {
    timeoutValue := 15
    flag.IntVar(
        &timeoutValue,
        "timeout",
        timeoutValue,
        "timeout fot HTTP-requests in sec",
    )
    flag.IntVar(
        &timeoutValue,
        "t",
        timeoutValue,
        "timeout for HTTP-requests in sec",
    )
    flag.Parse()
    if len(os.Args) < 2 {
        fmt.Println("Please provide a url")
        return
    }
    results := make(chan result)
    urls := flag.Args()
    client := http.Client{
        Timeout: time.Duration(timeoutValue) * time.Second,
    }
    for _, url := range urls {
        go func(url string) {
            resp, err := client.Get(url)
            if err != nil {
                results <- result{url:url, err: err,}
                return
            }
            defer resp.Body.Close()
            body, err := io.ReadAll(resp.Body)
            if err != nil {
                results <- result{url: url, err: err,}
                return
            }
            text := resp.Proto + " " + resp.Status + "\n"
            for name, values := range resp.Header {
                for _, value := range values {
                    text += name + ": " + value + "\n"
                }
            }
            text += "\n"
            text += string(body)
            results <- result{url: url, text: text,}            
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
    os.Exit(1)
}
