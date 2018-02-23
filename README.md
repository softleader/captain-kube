# Captain Kube

> A command line tool for Kubernetes hepls SoftLeader client deploying SoftLeader products

## Install

安裝的過程中, 必須要指定執行的系統 `$GOOS` 及 `$GOARCH`, 支援的參數可 [參考這篇](https://golang.org/doc/install/source#environment)

### By Docker


```shell
$ docker run -it --rm -v $BINARY:/data softleader/captain-kube GOOS=$GOOS GOARCH=$GOARCH
```

- `$BINARY` - 安裝的目錄, 如 `/usr/local/bin`
- `$GOOS` - 預設為 `linux`, 若輸入 `macos` 則自動轉換為 `darwin`
- `$GOARCH` - 預設為 `amd64`

> You need Docker installed

#### Example: 

安裝在 Mac 系統中的當前目錄

```shell
$ docker run -it --rm -v "$(pwd)":/data softleader/captain-kube GOOS=macos
```

安裝在 Ubuntu

```shell
$ docker run -it --rm -v /usr/local/bin:/data softleader/captain-kube
```

#### Tag

- `softleader/captain-kube:latest` - 每次都以最新的 captain-kube 的 master branch 去安裝

### By Makefile

```shell
$ go get github.com/softleader/captain-kube
$ make -C $GOPATH/src/github.com/softleader/captain-kube BINARY=$BINARY GOOS=$GOOS GOARCH=$GOARCH
```

- `$GOPATH` - 如果系統沒設定的話, 可以用 `$(go env GOPATH)` 代替, 或通常會裝在 `~/go` 這個目錄下
- `$BINARY` - 安裝的目錄, 預設 `./build`
- `$GOOS` - 預設為 `linux`, **(注意: 這邊不支援 `macos` 自動轉換為 `darwin`)**
- `$GOARCH` - 預設為 `amd64`

> You need Go, Git, Make installed
