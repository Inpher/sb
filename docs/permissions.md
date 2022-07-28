# Permissions

## General

Having an account on `sb` doesn't grant much privileges by itself except the permission to interact with `sb`.

A registered account will be able to:
- list all `sb`'s groups
- get partial information on these groups
- manage their own SSH key-pairs
- list their past sessions
- replay their past sessions as ttyrecs or as GIF
- access their own personal hosts configured (see [Self accesses](#self-accesses))
- manage their connection parameters (generate TOTP codes, ...)

## SSH key pairs: ingress vs. egress keys

`sb` distinguishes two types of SSH key-pairs:
1. ingress (from your computer -> `sb`)
2. egress (from `sb` -> distant host)

### Accessing sb with ingress keys

Ingress key pairs are personal to each account.

When a new account is created on `sb`, the `sb` owner asks the users for their first SSH public keys, 
and creates the accounts with them.

Once users have an account on `sb`, they can manage their own ingress keys with the following commands:
- `sb self ingress-keys list`: list the authorized SSH public keys
- `sb self ingress-key add`: add a new authorized SSH public key
- `sb self ingress-key delete`: delete an authorized SSH public key

### Accessing distant hosts with egress keys

It is not supported to import SSH key pairs in `sb` via the CLI. 
You can always generate new key pairs, and `sb` will output the public SSH keys that you can copy 
in the `.ssh/authorized_keys` file of the distant host.

There are two types of egress key pairs:
- the personal ones
- the group ones

The personal egress key pairs are dedicated to one account, and only the account has permissions on them.

The group egress key pairs are shared by all [members](#group-member) of a [group](#groups).

It is not the responsibility of `sb` to propagate the public keys to the distant hosts to allow access. 
These keys will have to be trusted in the `authorized_keys` file of the distant host to allow access via `sb`.

## Self accesses

Once they have generated an [egress key](#accessing-distant-hosts-with-egress-keys), accounts can add personal accesses 
to their own authorized hosts list.

This, of course, doesn't grant access to the distant host, and the public part of the egress key pair will have 
to be trusted on the distant host.

These self accesses can be useful to test things or to share `sb` between multiple usages (company and private).

Accounts can manage their own accesses with the following commands:
- `sb self accesses list`: list the personal accesses
- `sb self access add`: add a personal access to a host
- `sb self access remove`: delete a personal access to a host

## Groups

`sb` implements the concept of groups.

Groups are an easy way to share accesses between multiple accounts (e.g.: all members of the `dev` group can access 
the development hosts).

A group has four level of permissions:
1. `member`: members can access the hosts authorized for the group
2. `acl-keeper`: ACL keepers can manage the group's hosts
3. `gate-keeper`: gate keepers can manage the group's member
4. `owner`: owners can promote accounts to `gate-keeper`, `acl-keeper` and `owner` for the group

Groups have their own SSH key pairs, and only members of the group are able to use the group's private keys 
thanks to system level permissions.

### Group member

A group member has access to the hosts authorized by the group, with the group-shared SSH private key-pairs.

A group member will also be able to list all hosts accessible via the group (`sb group list accesses`).

### Group ACL keepers

A group ACL keeper can manage the hosts accessible by the group with the following commands:
- `sb group access add`: add an access to the group
- `sb group access remove`: remove an access from the group

### Group gate keepers

A group gate keeper can manage the members in the group with the following commands:
- `sb group member add`: add a member to the group
- `sb group member remove`: remove a member from the group

### Group owners

A group owner can manage the owners, gate keepers and ACL keepers of the group with the following commands:
- `sb group owner add`
- `sb group owner remove`
- `sb group gate-keeper add`
- `sb group gate-keeper remove`
- `sb group acl-keeper add`
- `sb group acl-keeper remove`

## Owners group

A special group exists in `sb`: `owners`.

An owner of the `owners` group gets access to the following additional commands:
- `sb account create`: create an account on `sb`
- `sb account delete`: delete an account from `sb`
- `sb group create`: create a group on `sb`
- `sb group delete`: delete a group from `sb`
