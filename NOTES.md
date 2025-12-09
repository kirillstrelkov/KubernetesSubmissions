# Notes

## Increase file watcher limits

```bash
# temporary
sudo sysctl -w fs.inotify.max_user_watches=524288
sudo sysctl -w fs.inotify.max_user_instances=8192

# make it permanent
echo "fs.inotify.max_user_watches = 524288" | sudo tee -a /etc/sysctl.conf
echo "fs.inotify.max_user_instances = 8192" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

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
