# LocalSend CLI

## Installation

You can download executable files from the Releases.

### Prerequisites

- [Go](https://golang.org/dl/) 1.22 or higher

### Clone the Repository

```sh
git clone https://github.com/ilius/localsend-go.git
cd localsend-go
```

### Build

Use the `Makefile` to build the program.

```sh
make build
```

This will generate binary files for all supported platforms and save them in the `bin` directory.

## Usage

### Run the Program

#### Receive Mode

```sh
./localsend-go -mode receive
```

Choose the appropriate binary file for your operating system and architecture.
On Linux, you need to execute this command to enable its ping functionality:
`sudo setcap cap_net_raw=+ep localsend-go`

#### Send Mode

```sh
./localsend-go -mode send -file FILE_PATH -to your_ip
```

Example:

```sh
./localsend-go -mode send -file ./hello.tar.gz -to 192.168.3.199
```

## Contribution

You are welcome to submit issues and pull requests to help improve this project.

## License

[MIT](LICENSE)

# Todo

- \[ \] Improve send functionality to display sent text directly on the device
