# healthchecker

Healthchecker for docker. From security point of view using curl is bad idea.
This tool is limited to request 127.0.0.1 only. By the way you can build it
with any target you want. On start it will check against hostname const.

## Installation
```sh
$ make
$ make install
```

## Usage
```text
healthchecker [URL (default: http://127.0.0.1:8000/ping)]
  -H string
    	Header
  -V	Print version and exit
  -f	Fail silently (default true)
  -s	InsecureSkipVerify (default true)
  -t duration
    	Timeout
```

## Example

Dockerfile:
```text
HEALTHCHECK CMD healthchecker http://127.0.0.1:8000/ping || exit 1
```
