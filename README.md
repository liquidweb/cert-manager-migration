# cert-manager-migration
Cert Manager Migration from BoltDB to PostgreSQL

## Build Instructions
1. Clone repository to your `$GOPATH/src/github.com/liquidweb`
2. [Download](https://golang.github.io/dep/docs/installation.html) dep if you haven't already   
3. Download dep dependencies for this project
   `dep ensure`
4. Build this project
   `go build migration.go`
5. Run this project 
   `go run migration.go -data-dir={path to bolt db file}` (Code assumes you have a data.db file) (You can leave data-dir file off if data.db is in your current working directory)
