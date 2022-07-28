# Genesis and core ideology

## Genesis

[@bodji](https://github.com/bodji) and [I](https://github.com/golgeek) worked for the French Cloud provider 
[OVHcloud](https://github.com/ovh) for about 8 years.

For about these 8 years, our job heavily relied on connecting to various machines, and every connection was made 
through an in-house SSH bastion that was simplifying the process (by a lot!).

When we left OVHcloud early 2020 and started to work at [Inpher](https://inpher.io), we had a similar need for an SSH 
bastion and no real solution, as [OVHcloud's solution](https://github.com/ovh/the-bastion) was only open-sourced 
mid-2020.

This is when we started to code our own version, in Golang, based on a very similar philosophy.

## Core ideology

A few principles drove the development of this project:
- A bastion has to be the last piece of infrastructure standing and has to rely on bulletproof technology
- Its security must be airtight
- Everything must be auditable
- It has to be extremely easy to deploy and setup
- A multi-primary replication was very high on the list of features
- It must bring the same set of core features as the one we were using before

### Resilience

`sb` is a Golang binary and can work without internet access.

Every information is stored as a file or in SQLite databases on the filesystem by a binary that spawns and exits 
at the end of the session.

The only exception to this rule is [the multi-primary replication](#multi-primary-replication).

### Security

Every `sb` user is a system user and every `sb` group is a system group.

This allows the initial authentication to `sb` to be handled via the SSHd process and is secured through PKI exchange.

Once a user is authenticated and the session is established, the `sb` binary acts as the user's shell, and is executed 
by the SSHd process with the system user's permission, granting that the users have unprivileged access to the host, 
and that nobody gets an actual shell.

Every `sb` command has associated [permissions](./permissions.md), and they are handled by `sb`'s core 
before the command is executed. The core also logs every action before even calling the subcommand, ensuring that there 
will always be a trace.

To ensure that nothing bad can happen in case the code is wrong, every personal bit of information (private SSH keys, 
personal accesses, personal sessions) are only accessible by the system user who generated them on a system level 
via Unix permissions. And, of course, every group's bit of information (private SSH keys, group accesses) are only 
accessible by the system users who are in the Unix system group via Unix permissions.

For elevated privilege actions (add a user to a group, create a new user, ...), `sudoers.d` files are created 
with very specific commands available to very specific users and group members.

### Auditability

As every action comes from a SSH connection, the auditability starts from the SSH connection log.

Every time the binary is started, it outputs a log in its SQLite database, this log will be amended every time 
a new piece of information is available:
- is the command valid?
- does the user have access to the host?
- was the remote connection successful?

On top of these, every SSH connection is recorded via TTYRec and can be replayed.

### Simplicity

As the code is written in Golang, it is shipped as a single binary, simplifying the installation process.

The TTYRecs records, playbacks and conversion as GIF are handled in pure Golang, limiting the number 
of external softwares to install to [mostly coreutils](./installation.md#prepare-your-system).

### Multi-primary replication

Resilience is great, but your host can go down, and this is why a multi-primary replication was important to us.

Having multiple instances of `sb` being able to handle any command at any time allows you to have:
- high availability
- geo replication

This part is handled via a long running daemon (`sb -d`) that handles the [replication](./high-availability.md#replication) through a message queue.

### Set of features

This part is covered in [the dedicated documentation](./features.md)
