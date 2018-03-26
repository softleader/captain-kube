![](./doc/captainkube-01.svg)

> A command line tool for Kubernetes helps SoftLeader client deploying SoftLeader products

## Install

安裝的過程中, 必須要指定執行的系統 `$GOOS` 及 `$GOARCH`, 支援的參數可 [參考這篇](https://golang.org/doc/install/source#environment)

```shell
$ docker run -it --rm -v $BINARY:/data \
	softleader/captain-kube \
	GOOS=$GOOS GOARCH=$GOARCH
```

- `$BINARY` - 安裝的目錄, 如 `/usr/local/bin`
- `$GOOS` - 預設為 `linux`, 若輸入 `macos` 則自動轉換為 `darwin`
- `$GOARCH` - 預設為 `amd64`

> You need Docker installed

### Tag

- `softleader/captain-kube` - 以最新 release 的程式安裝
- `softleader/captain-kube:remote` - 每次都以 captain-kube 當前的 master branch 去安裝

### Example: 

安裝在 Mac 系統中的當前目錄

```shell
$ docker run -it --rm -v "$(pwd)":/data softleader/captain-kube GOOS=macos
```

安裝在 Ubuntu

```shell
$ docker run -it --rm -v /usr/local/bin:/data softleader/captain-kube
```

## Usage

![](./doc/overview.svg)
