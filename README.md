# Go Common

A Common Library used to abstract models, packages, utils and basic implementation out of the main codebase to
maintain a standard codebase.

___

**This project was written in go@1.22.0**

## Usage

- To use the library the following needs to be set:
### NetRC
A .netrc file is a mechanism used to store 
authentication-related information for remote servers. 
It typically resides in the user’s home directory. 
This file is widely used in Unix systems and by tools like curl.

To set this up run:
```shell
echo "machine github.com login <github_username> password <github_access_token>" > ~/.netrc
```

## GitConfig
The developers ~/.gitconfig file should contain the following:
```text
[url "ssh://git@github.com/"]
        insteadOf = https://github.com/
```

## Go ENV
Typically after setting up the network the developer should set up their GOPRIVATE variable,
This tells go-toolchain to fetch data from a specific private source.

### Project Level
```shell
export GOPRIVATE=github.com/shiroyaavish/*
```

### Global
```shell
go env -w "GOPRIVATE=github.com/shiroyaavish/*"
```

## Implementations

- HTTP Server
- GRPC Server
- Postgres Database
- ProtoBuffs
- Redis
- AWS
- Config Loader
    - .env
    - .json
- Error Abstractions
- Data Structures
    - LinkedList

## Note

All Abstraction pertaining to individual code to be added to Common and used throughout.

## Maintainers

@shuklasaharsh

