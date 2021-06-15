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

## Server configuration
You can provide environment variables HOSTS and PORT to customize your server network binding.

Defaults - 
- HOSTS=localhost (0.0.0.0 for the docker image)
- PORT=9876

## Querying and connecting to CrankDB

As this is a very early release you can use 2 ways to interact with the database -

- [Crank CLI](https://github.com/shreybatra/crank)
- [Cranky - Python Driver for CrankDB](https://github.com/shreybatra/Cranky)


## Steps to build
- Clone the repo and change directory to project root folder.
- Tidy dependencies using - `go mod tidy`
- Build the application - `go build .`
- Run the server - `./crankdb`
