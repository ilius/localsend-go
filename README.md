<div align="center">
<h1>LocalSend CLI</h1>
<img src="doc/images/image.png" alt="LocalSend CLI logo" width="150" height="150">
<p>✨LocalSend CLI✨</p>
</div>

## Documentation

[Chinese](doc/README_zh.md) | [English](doc/README_en.md)

## Installation

> 😊You can download the executable file in Release

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or higher

### Clone the repository

```sh
git clone https://github.com/ilius/localsend_cli.git
cd localsend_cli
```

### Compile

Use `Makefile` to compile the program.

```sh
make build
```

This will generate binaries for all supported platforms and save them in the `bin` directory.

## Usage

### Run the program

#### Receive mode

```sh
.\localsend_cli-windows-amd64.exe -mode receive
```

Select the appropriate binary to run based on your operating system and architecture.

In Linux, you need to execute this command to enable its ping function
`sudo setcap cap_net_raw=+ep localsend_cli`

#### Send mode

```
.\localsend_cli-windows-amd64.exe -mode send -file ./xxxx.xx -to your_ip
```

example:

```
.\localsend_cli-windows-amd64.exe -mode send -file ./hello.tar.gz -to 192.168.3.199
```

## Contribution

Welcome to submit issues and pull requests to help improve this project.

## License

<!-- [MIT](LICENSE) -->

# Todo

- \[ \] Improve the sending function. The sent text can be displayed directly on the device
