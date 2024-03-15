# Tachiyomi-To-FMD2 Converter

Converter to convert a Tachiyomi Backup to a `favorites.db` to use with [FMD2](https://github.com/dazedcat19/FMD2).

## Installation

```bash
$ git clone https://github.com/dix0nym/Tachiyomi-To-FMD2-Converter
$ ./gen-proto.sh
# prepare mapping before running
$ go run main.go
```

## Usage

```
# prepare Tachiyomi Backup
$ gunzip tachiyomi_2022-03-15_20-42.proto.gz

# use provided mapping or update from source
$ git clone https://github.com/dazedcat19/FMD2
$ python gen-module-mapping.py
```