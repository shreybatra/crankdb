# CrankDB

CrankDB is an ultra fast and very lightweight Key Value based Document Store.

## Requirements
- Golang 1.16

## Steps to deploy

#### Using docker image (Recommended)
```
docker run -p 9876:9876 shreybatra/crankdb
```

#### Download binary executable
- Download executable -
    - [Linux](https://crankdb.blob.core.windows.net/crankdb/crankdb_v0_1-beta_0_linux.tar)
    - [MacOs](https://crankdb.blob.core.windows.net/crankdb/crankdb_v0_1-beta_0_macos_darwin.tar)
- Extract tar and start server -
    ```
    cd Downloads
    tar -xvf <downloaded_tar_file>
    ./crankdb
    ```
- (MacOS) You might need to allow macos to run the file via Settings and Privacy.

#### Using go get command
- Download application - `go get github.com/shreybatra/crankdb`
- Run server with command - `crankdb`

## Querying and connecting to CrankDB
- Download [crank-cli](https://github.com/shreybatra/crank)
- Check [docs](https://github.com/shreybatra/crank#query-language) for query language.

## Steps to build
- Clone the repo and change directory to project root folder.
- Tidy dependencies using - `go mod tidy`
- Build the application - `go build .`
- Run the server - `./crankdb`
