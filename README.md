Basic URL crawler. Takes a CSV file containing URLs to run HTTP GETs for and returns response codes and content length

In conjunction with domain databases, [such as this one](https://www.domcop.com/top-10-million-domains), it can be used for stress testing and troubleshooting of HTTP/HTTPS-based content filtering systems on network firewalls.

Example usage:
- Crawl every URL in the CSV file

`go run gowler.go --csv <CSV_file>`

- Only crawl the first n URLs in the CSV file

`go run gowler.go --csv <CSV_file> --num <n>`

- Run crawler using n goroutines in parallel

`go run gowler.go --csv <CSV_file> --goroutines <n>`
