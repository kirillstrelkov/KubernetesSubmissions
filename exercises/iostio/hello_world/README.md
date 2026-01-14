# Hello world with mesh

As docker images from <https://github.com/mluukkai/kube-hello> are not for linux amd64, we need to build them manually and import into cluster:

```bash
make docker
```

After that check [./Makefile](./Makefile) and you can use `deploy*` targets.

## Part 3

For this part istio ingressgateway should be install, check [../Makefile](../Makefile):

```bash
cd .. # go to parent folder
make install-ingressgateway
```
