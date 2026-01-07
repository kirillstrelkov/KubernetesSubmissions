# DummySite custom resource

## Development

Install kubebuilder - <https://book.kubebuilder.io/quick-start>

```bash
go mod init stable.dwk
kubebuilder init --domain stable.dwk

kubebuilder create api --group stable.dwk --version v1 --kind DummySite

# edit next files:
# ./api/v1/dummysite_types.go
# ./internal/controller/dummysite_controller.go

# build yaml file and install it to the cluster:
make manifests
make install

# this will start the controller and attach it to the cluster
make run
```

## Use DummySite

To create 2 DummySites from [./config/samples](./config/samples):

```bash
make apply-dummysite
```

> In order to access them via localhost, ports should be forwarded. Check output of make command or Makefile for specific command.

To remove DummySites:

```bash
make clean-dummysites
```
