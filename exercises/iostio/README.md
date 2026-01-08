# Iostio

Stop default cluster

```bash
k3d cluster stop k3s-default
```

Check [./Makefile](./Makefile) for more information:

```bash
# create and install iostio
make k3d-iostio

# create sample apps
make k3d-iostio-apply-sample

# to see app run port forwaring
make port-forward-iostio

# install kiali
make install-kiali

# open kiali dashboard
make kiali

# generate load
make send-100requests

# to enable L4 auth policy
make apply-auth-policy

# to enable L7 and L4 auth policies
make enable-l7

# split review routing
make split-traffic

# test that routing is split between v1 and v2 for reviews
make test-traffic-split
```
