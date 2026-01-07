# DummySite custom resource

Install kubebuilder - <https://book.kubebuilder.io/quick-start>

```bash
go mod init stable.dwk
kubebuilder init --domain stable.dwk

kubebuilder create api --group stable.dwk --version v1 --kind DummySite

# edit next files:
# ./api/v1/dummysite_types.go
# ./internal/controller/dummysite_controller.go

make manifests

```
