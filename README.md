# cert-manager-migration
Cert Manager Migration from BoltDB to PostgreSQL

## Build Instructions
1. Clone repository to your $GOPATH/github.com/liquidweb 
2. Download glide if you haven't already
   `go get github.com/Masterminds/glide`
3. Download glide dependencies for this project
   `glide install`
4. Build this project
   `go build migration.go`
5. Run this project 
   `go run migration.go -data-dir={path to bolt db file}` (Code assumes you have a data.db file) (You can leave data-dir file off if data.db is in your current working directory)
