# hosts-go

中文 | [English](README.EN.md)

## 概述

hosts-go 是一个用于从互联网上获取和合并 hosts 文件的命令行工具。它可以定期从指定的 URL 获取 hosts 文件，并将其合并到本地的 hosts 文件中。

## 安装

### 使用预编译的二进制文件
```bash
curl https://github.com/hunshcn/hosts-go/releases/latest/download/hosts-go_linux_amd64 -L -o /usr/bin/hosts-go && chmod +x /usr/bin/hosts-go
```

### go install
```
go install github.com/hunshcn/hosts-go@latest
```

## 使用
hosts-go 提供了以下命令行选项：

- `--url` 或 `-u`：指定要获取 hosts 文件的 URL。可以指定多个 URL。
- `--test` 或 `-t`：仅输出合并后的 hosts 文件内容。
- `--content-only`：仅输出获取的 hosts 文件内容。
- `--service` 或 `-s`：安装或卸载 hosts-go 作为系统服务。
- `--duration` 或 `-d`：指定更新 hosts 文件的时间间隔，默认为 1 小时。
- `--reload-command`：在更新成功 hosts 文件后执行的命令。

### 示例

获取并合并 hosts 文件：

```
hosts-go -u https://gitlab.com/ineo6/hosts/-/raw/master/next-hosts
```

安装 hosts-go 作为系统服务（使用输入的参数）：

```
hosts-go -u https://gitlab.com/ineo6/hosts/-/raw/master/next-hosts -s install
```

卸载 hosts-go 服务：

```
hosts-go -s uninstall
```

### 注意事项

- 在运行 hosts-go 之前，请确保您具有足够的权限来读取和写入 hosts 文件。

## License

MIT License.

## 其他

本项目不提供 hosts，比较类似 [SwitchHosts](https://github.com/oldj/SwitchHosts) 的功能。不过什么都写好了才发现好像 hosts 效果还是太差。