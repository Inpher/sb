# Demo

This project comes with a docker-compose configuration that presents a demo of how `sb` works

It will spawn 5 services:
1. A simple Debian Bullseye container with an SSH service that we want to log on to
2. An `sb` container
3. Another `sb` container (both replicated)
4. A `Google PubSub` emulator, used for replication between the two `sb` instances
5. A `Google Cloud Storage` emulator, used to offload SSH sessions recordings to an external storage

## Go to the demo directory and start the docker-compose configuration

```console
t1000@skynet:~# cd demo
t1000@skynet:~# docker-compose up -d
```

## Ensure the permissions for the demo private SSH key are correct

```console
t1000@skynet:~# chmod 600 $(PWD)/assets/ssh-keys/id_ed25519
```

## Configure two aliases to simplify your interactions with sb

```console
t1000@skynet:~# alias sb1="ssh -i $(PWD)/assets/ssh-keys/id_ed25519 -p 22001 t800@127.0.0.1 -A -t -- "
t1000@skynet:~# alias sb2="ssh -i $(PWD)/assets/ssh-keys/id_ed25519 -p 22002 t800@127.0.0.1 -A -t -- "
```

## Generate an SSH key for the default account

We are now going to generate a personal egress key (`sb` -> distant host) for your account.

```console
t1000@skynet:~# sb1 self egress-key generate --algo ed25519 --size 256
```

## Push the public SSH key to the example container we will want to access

The docker-compose configuration created a Debian Bullseye container acting as a remote machine.

To be able to access through sb, we need to add the newly generated SSH public egress key to the authorized_keys of this machine:

```console
t1000@skynet:~# export TRUSTED_KEY=$(sb1 self egress-keys list | sed -e 's/\x1b\[[0-9;]*m//g' | grep -e "^1:" | sed -e 's/^1: //')
t1000@skynet:~# docker exec -it sbdemo_examplevm /bin/bash -c "echo '$TRUSTED_KEY' > /root/.ssh/authorized_keys"
```

## Authorize your default account to access the example container

Now, let's authorize our account to access the distant host on port 22 with the user root:

```console
t1000@skynet:~# sb1 self access add --host examplevm --port 22 --user root
```

## Access the example container

Now, just connect to the example container via the other `sb` instance:

```console
t1000@skynet:~# sb2 root@examplevm
```

Congrats, you just made your first connection to a distant host through `sb`!