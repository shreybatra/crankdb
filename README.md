# CrankDB

CrankDB is an ultra fast and very lightweight Key Value based Document Store.

### Latest version - v0.1.0-beta.2

## Requirements
- Golang 1.16

## Steps to deploy

#### (Option 1) Download binary executable
- Download executable (v0.1.0-beta.1) -
    - [Linux](https://crankdb.blob.core.windows.net/crankdb/crankdb_v0_1-beta_2_linux.tar)
    - [MacOs](https://crankdb.blob.core.windows.net/crankdb/crankdb_v0_1-beta_2_darwin.tar)
- Extract tar and start server -
    ```
    cd Downloads
    tar -xvf <downloaded_tar_file>
    ./crankdb
    ```
- (MacOS) You might need to allow macos to run the file via Settings and Privacy.

#### (Option 2) Using docker image (Recommended)
```
docker run -p 9876:9876 ahsanbarkati/crankdb
```

#### (Option 3) Using go get command
- Download application - `go install github.com/ahsanbarkati/crankdb@latest`
- Run server with command - `crankdb`

## Server configuration
You can provide environment variables HOSTS and PORT to customize your server network binding.

Defaults - 
- HOSTS=localhost (0.0.0.0 for the docker image)
- PORT=9876

## Querying and connecting to CrankDB

As of now, we can query and use the database via the crank-cli or SDKs in the following languages - 

| Language                | SDK/Tool    | Latest Version |
|-------------------------|-------------|----------------|
| CLI (command line tool) | [Crank CLI](https://github.com/shreybatra/crank) | v0.1.0-beta.1  |
| Golang                  | [Gocrank](https://github.com/shreybatra/gocrank)   | v0.1.0-beta.1  |
| Python                  | [Cranky](https://github.com/shreybatra/Cranky)    | 0.1.0b1        |



## Steps to build
- Clone the repo and change directory to project root folder.
- Tidy dependencies using - `go mod tidy`
- Build the application - `go build .`
- Run the server - `./crankdb`
