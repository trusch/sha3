sha3
====

A simple command line tool to work with the sha3 hash family.

## Installation

```bash
go get github.com/trusch/sha3
```

## Usage
```bash
> sha3 --help
Usage of sha3:
  -c, --check         check sum files
  -h, --help          show usage info
  -l, --length int    output length in bytes if using a shake hash (default 32)
  -t, --type string   one of 'sum224', 'sum256', 'sum384', 'sum512', 'shake128', 'shake256' (default "shake256")
```

## Examples

### Get 32 byte of the shake256 hash of a file:

```bash
> sha3 LICENSE
2198cf7b3addc20267107a90795c77845d4db24f4a54d77ba524b7cb87783510  LICENSE
```

### Get 64 byte of the shake256 hash of a file:

```bash
> sha3 --length 64 LICENSE
2198cf7b3addc20267107a90795c77845d4db24f4a54d77ba524b7cb8778351075f57cf6c53c238e286b7c42dd90a3c4adc92b04cc2fc1205483d108db624354  LICENSE
```

### Get the sha3-256 hash of a file

```bash
> sha3 --type sum256 LICENSE
8a6874403b717b0f66552d0ef6c1237de3e8b79b3d08da25bbd3f30dc00e7aa8  LICENSE
```

### Create a sum file of all files in a directory and verify it

```bash
> sha3 . > /tmp/sums.txt
> cat /tmp/sums.txt
c8763dd367aad67181a1d25d3ee663ffd05dbdf1a7f3af02fa82893b9c084373  .git/HEAD
e448e661ae2b1fb6a26af2abdd38004552e9beeb8ef389bd7afc8e6563f34e5f  .git/config
a6bee17be622db8c261d64ba3fcc3ae84772989873beb4b21ee9e53c72b4c119  .git/description
0a232b6c99c2d8211aad32be6d271890a50ee59bd447c97f698b2cadc2e30886  .git/info/exclude
2198cf7b3addc20267107a90795c77845d4db24f4a54d77ba524b7cb87783510  LICENSE
23a97616a5a74ce0078f9633bc195efa8508057fcedcce2b477631f6118acb4b  README.md
12b5cf1cc2ccf4e2c85b233a435275adf1ca02e70f6a01d23e4d31d1ca666e6c  go.mod
0963e30b7af1a3b770e94d78692c0112cc9e711580b0bf0df27325073d37e326  go.sum
6dd7fa48b87f1d1190b1b083f6563e57be7b1fce8fdcd3bd5805413eb798d670  main.go
> sha3 -c /tmp/sums.txt
.git/HEAD: OK
.git/config: OK
.git/description: OK
.git/info/exclude: OK
LICENSE: OK
README.md: OK
go.mod: OK
go.sum: OK
main.go: OK
```
