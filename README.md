
# Quick start

```bash
# requirements: installed docker-compose, running docker
make test
make test-canon
```

# Overview

Make a git checkout for the *public* git repo via [go-git](https://github.com/go-git/go-git) via ssh-transport does not work because the public key is not added to the github repository.

case 1: Auth via `ssh.PublicKeysCallback` (SSHAgentAuth)

```
ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain
```

case 2: AuthMethod is nil. It is error is expected - not found known_hosts (e.g. `ssh-keyscan -t rsa github.com >> /root/.ssh/known_hosts`) or add `ssh.InsecureIgnoreHostKey()` in auth method.

```
unable to find any valid known_hosts file, set SSH_KNOWN_HOSTS env variable
```

case 3: Auth via `ssh.Password`

```
ssh: handshake failed: ssh: unable to authenticate, attempted methods [none], no supported methods remain
```

case 4: Auth via `ssh.PublicKeys`

```
ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain
```

For *public* repo git checkout from anywhere only via https-transport

case 5: via https
