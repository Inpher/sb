# Setup of your first account

Once your `sb` instance is up and running, you can setup your own account and groups.

This tutorial assumes you have access to one `sb` instance host of the cluster you want to configure, 
and that you already have an SSH key pair generated.

# Create your own account

You can create your first user with one command executed on one instance of the cluster.

Please update the following command with your desired username and your public SSH key:

```console
root@sb-host1:~# /opt/sb/sb account create --username [USERNAME] --public-key "'[YOUR SSH PUBLIC KEY]'"
Adding user
Creating home skeleton
Pushing pk in authorized_keys file
User [USERNAME] was successfully created
```

# Create the owners group

As described in the [permissions documentation](./permissions.md), 
[one group has privileged permissions](./permissions.md#owners-group).

Now that your own user has been created, we're going to create this group and set your user as owner.

```console
root@sb-host1:~# /opt/sb/sb group create --name owners --owner-account [USERNAME] --algo ed25519 --size 256
Group owners was successfully created
The new egress private SSH key was successfully generated
```

# Physically allow the owners group to access the instance

To fully secure your `sb` instance, as of when this setup is over, you should only access it through `sb`. 

To ensure that the owners group will have access to it, just add the key that was generated 
to the `authorized_keys` file of the root user:

```console
root@sb-host1:~# cat /home/bg_owners/.ssh/*.pub >> /root/.ssh/authorized_keys
```

# Setup your local machine

Now that you have your own user, you don't need to setup the `sb` alias that will allow you to use `sb`.

On your local machine, after replacing `[USERMAME]` with the username of the user you just created 
and `[SB_INSTANCE.domain]` with the DNS name or the IP address of one of your `sb` instances, 
add this alias to your shell.

```console
t1000@skynet:~# alias sb="ssh [USERNAME]@[SB_INSTANCE.domain] -A -t --"
```

Then, test that your connection to `sb` works properly:

```console
t1000@skynet:~# sb info
Hi, [USERNAME]!

Here are a few information about the wonderful place you just connected to:
  -> my name is sb-host1
  -> I'm located in us, my address is 127.0.0.1
  -> here is the list of my IP addresses: 127.0.0.1, 192.168.192.5
  -> I currently have 1 guest account
  -> I'm currently hosting 1 group

If you want to interact with me, just one of these aliases:
  -> SSH:  alias sb='ssh [USERNAME]@sb-host1.domain -t -A --'
  -> MOSH: alias msb='mosh [USERNAME]@sb-host1.domain -A --'

But let's talk a bit about yourself, [USERNAME]!
  -> you're a member of 1 groups
  -> TOTP is disabled on your account
```

If everything works properly, you can make your alias permanent.

# Allow the owners group to access the instance

From your local machine, programmatically trust the access to your `sb` instance(s) for the owner group.

For each instance, you will have to get its IP address (reachable by other instances), and optionally define an alias.

Let's say the first instance's IP address is `10.0.0.1` and we want to alias it `sb-host1`:

```console
t1000@skynet:~# sb group access add --host 10.0.0.1 --port 22 --user root --alias sb-host1
The access has been granted: Prefix: 10.0.0.1/32 | Host: 10.0.0.1 | Alias: sb-host1 | User: root | Port: 22
```

You can now connect to your `sb` instance via `sb`:

```console
t1000@skynet:~# sb sb-host1
*------------------------------------------------------------------------------*
|THIS IS A PRIVATE COMPUTER SYSTEM, UNAUTHORIZED ACCESS IS STRICTLY PROHIBITED.|
|ALL CONNECTIONS ARE LOGGED. IF YOU ARE NOT AUTHORIZED, DISCONNECT NOW.        |
*------------------------------------------------------------------------------*
Hey! Welcome to sb-host1, [USERNAME]!
Access to this host is granted by:
- group access from group owners
... connecting you to the distant host (if it's alive :)) ...
---
Linux c9d86d6c9976 5.10.76-linuxkit #1 SMP Mon Nov 8 10:21:19 UTC 2021 x86_64

The programs included with the Debian GNU/Linux system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.

Debian GNU/Linux comes with ABSOLUTELY NO WARRANTY, to the extent
permitted by applicable law.
root@c9d86d6c9976:~# 
```

You can now secure your SSH access on the `sb` instance to make sure that the connections can only come from `sb`.