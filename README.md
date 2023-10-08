# KVDB

`kvdb` is a lightweight, filesystem-backed, cli-driven key-value store written
in Go. This project started as a submission to a 4-hour coding challenge.

## Get Started & Usage

To get started, clone this repository and install the `kvdb` CLI using
`go install`

```console
$ git clone github.com/adowair/kvdb
...
$ cd kvdb
$ go install .
$ kvdb --help
Kvdb is a lightweight, filesystem-backed key-value store written in Go.
By default, kvdb uses the current directory to store data. To get started:

        kvdb set <key> <value>
        kvdb get <key>

Usage:
  kvdb [command]

Available Commands:
  del         Delete a key
  get         Get the value for a key
  help        Help about any command
  set         Set the value for a key
  ts          Get the created and last-modified timestamps for a key

Flags:
  -h, --help   help for kvdb

Use "kvdb [command] --help" for more information about a command.
```

## Features
- `kvdb` supports valid utf8 strings not containing path delimiters
  ('/', '\') as keys and values.
- By design, `kvdb` does not support empty values for keys. This was done to
  reduce confusion. To unset a key, simple delete it.
- Kvdb operations are atomic. Multiple processes can safely invoke `kvdb` to
  safely read and write the same keys.

## Possible Improvements
- Use shared (not exclusive) locks for reads.
- Support nested keys (i.e. "key/subkey").
- Implement configurable database location, backed by a config file.
- Lazily read data files, saving some work when only timestamps are requested.
- Retain key's metadata when it is deleted, "remembering" first-set timestamps.