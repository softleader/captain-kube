![](./docs/captainkube-01.svg)

> A command line tool for Kubernetes helps SoftLeader client deploying SoftLeader products

Documentation and Other Links:

- [Setup Documentation](https://github.com/softleader/captain-kube/wiki/Installation)
- [Usage Documentation](https://github.com/softleader/captain-kube/wiki/Getting-Started)
- [Manual Documentation](./docs/man/ckube.md)


## Set up the installation environment

```
$ docker run --rm -v "$(pwd)":/data hub.softleader.com.tw/captain-kube sh /initial.sh
```

## staging

```
$ 18:28:19 ➜ anible docker run -it --rm -v "$(pwd)":/data ck ansible-playbook -i /data/hosts staging.yml
```