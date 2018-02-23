# Captain Kube

> A command line tool for Kubernetes hepls SoftLeader client deploying SoftLeader products

## Install

### By Docker

```shell
$ docker run -it --rm -v $PATH:/data softleader/captain-kube GOOS=$GOOS GOARCH=$GOARCH
```

- `$PATH` - 安裝的目錄, 如 `/usr/local/bin`
- `$GOOS` - 預設為 `linux`, 若輸入 `macos` 則自動轉換為 `darwin`
- `$GOARCH` - 預設為 `amd64`

> `$GOOS` 及 `$GOARCH` 支援的參數可 [參考這篇](https://golang.org/doc/install/source#environment)

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

### By Make

在此專案目錄下

```shell
$ make BINARY=$BINARY GOOS=$GOOS GOARCH=$GOARCH
```

- `$BINARY` - 執行檔的目錄, 預設 `./build`
- `$GOOS` - 預設為 `linux`, **(注意: 這邊不支援 `macos` 自動轉換為 `darwin`)**
- `$GOARCH` - 預設為 `amd64`

> `$GOOS` 及 `$GOARCH` 支援的參數可 [參考這篇](https://golang.org/doc/install/source#environment)
