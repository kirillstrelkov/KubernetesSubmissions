# Notes

## Reset k3d

```bash
make
```

## Working with SOPS and age

### Create key with age

```bash
gen-key
```

`key.txt` will be create in root directory

### Use sops to encrypt file

```bash
export SOPS_AGE_PUB_KEY=$(grep '# public key:' key.txt | cut -d ':' -f 2 | tr -d ' ')
sops --encrypt --age $SOPS_AGE_PUB_KEY --encrypted-regex '^(data)$' secret.yaml > secret.enc.yaml
```

NOTE: choose proper regex `^(data)$` or `'(Data)$'`

### Use sops to decrypt file

```bash
export SOPS_AGE_KEY_FILE=../key.txt
sops --decrypt manifests/enc/secrets.yaml
```
