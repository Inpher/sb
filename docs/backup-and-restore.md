# Backup and restore

As everything in `sb` is stored as file (users, groups, databases, ...) and uses system permissions, 
it's pretty easy to backup and restore a cluster.

Because backups are extremely important for the main door to your infrastructure, this project comes 
with built-in commands to make backups and restores even easier.

## Backup

Connected as `root` on one instance of `sb`, you can use the following command:

```console
root@sb-host1:~# /opt/sb/sb backup --backup-directory /tmp
New backup available: /tmp/sb-backup_sb-host1_c9d86d6c9976_20220614T222746Z.bin
```

This `.bin` backup is actually a `.tar.gz` file encrypted with the [configuration's](./configuration.md#general) 
encryption-key.

It contains the following files and directories:
- /etc/shadow
- /etc/group
- /etc/passwd
- /etc/sudoers.d
- /home/sb/logs.db
- the home folders of every `sb` users
- the home folers of every `sb` groups

## Restore

To restore a backup, it is recommended to stop all [replication daemons](./installation.md#setup-the-daemon) 
in the cluster.

Then, you have to import the backup file to each instance you want to restore it on and know the encryption key 
that was used to create the backup.

With everything setup, you can just execute the following command:

```console
root@sb-host1:~# /opt/sb/sb restore --file [FILE] --decryption-key [KEY]
Backup successfully restored
```

Don't forget to restart the replication daemons after the restore operations.