# cert-manager-migration
Cert Manager Migration from BoltDB to PostgreSQL

## Build Instructions
1. Clone repository to your `$GOPATH/src/github.com/liquidweb`
2. [Download](https://golang.github.io/dep/docs/installation.html) dep if you haven't already   
3. Download dep dependencies for this project
   `dep ensure`
4. Modify your `conf.yaml` file with appropriate values
5. Build this project
   `go build migration.go`
6. Modify your `conf.yaml` file with appropriate values   
7. Run this project
    1. Migrate Data `./migration.go migrate` 
    2. Print Bolt Data `./migration.go print-bolt-data`
    3. Create DB Tables `./migration.go create-tables`
    4. Drop DB Tables `./migration.go drop-tables`
    5. Migrate Certs and Secrets `./migration.go kube-migrate`
    
   
