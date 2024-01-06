# hosts-go

[中文](README.md) | English

## Overview

hosts-go is a command-line tool for fetching and merging hosts files from the internet. It can periodically fetch hosts files from specified URLs and merge them into the local hosts file.

## Installation

### use pre-compiled binary
```bash
curl https://github.com/hunshcn/hosts-go/releases/latest/download/hosts-go_linux_amd64 -L -o /usr/bin/hosts-go && chmod +x /usr/bin/hosts-go
```

### go install
```
go install github.com/hunshcn/hosts-go@latest
```

## Usage

hosts-go provides the following command-line options:

- `--url` or `-u`: Specify the URLs to fetch hosts files from. Multiple URLs can be specified.
- `--test` or `-t`: Only output the merged hosts file content.
- `--content-only`: Only output the fetched hosts file content.
- `--service` or `-s`: Install or uninstall hosts-go as a system service.
- `--duration` or `-d`: Specify the duration between each fetch of hosts files. The default is 1 hour.
- `--reload-command`：Command to execute after successfully updating the hosts file.

### Examples

Fetch and merge hosts files:

```
hosts-go -u https://gitlab.com/ineo6/hosts
```

Install hosts-go as a system service:

```
hosts-go -u https://gitlab.com/ineo6/hosts -s install
```

Uninstall hosts-go service:

```
hosts-go -s uninstall
```

### Notes

- Before running the hosts-go command, make sure you have sufficient permissions to read and write the hosts file.

## License

MIT License.