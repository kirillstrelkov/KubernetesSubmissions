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
sops --encrypt \
       --age $(grep '# public key:' key.txt | cut -d ':' -f 2 | tr -d ' ')$ \
       --encrypted-regex '^(data)$' \
       secret.yaml > secret.enc.yaml
```

### Use sops to decrypt file

```bash
export SOPS_AGE_KEY_FILE=key.txt
sops --decrypt secret.enc.yaml | kubectl apply -f -
```
