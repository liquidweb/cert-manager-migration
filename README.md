# cert-manager-migration
Cert Manager Migration from BoltDB to PostgreSQL

## Build Instructions
1. Clone repository to your `$GOPATH/src/github.com/liquidweb`
2. [Download](https://golang.github.io/dep/docs/installation.html) dep if you haven't already   
3. Download dep dependencies for this project
   `dep ensure`
4. Modify your `conf.yaml` file with appropriate values
5. Build this project
   `go build`
6. Modify your `conf.yaml` file with appropriate values   
7. Run this project
    1. Migrate Data `./cert-manager-migration migrate` 
    2. Print Bolt Data `./cert-manager-migration print-bolt-data`
    3. Create DB Tables `./cert-manager-migration create-tables`
    4. Drop DB Tables `./cert-manager-migration drop-tables`
    5. Migrate Certs and Secrets `./cert-manager-migration kube-migrate`
    
   
