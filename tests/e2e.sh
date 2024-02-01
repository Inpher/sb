#!/usr/bin/env bash

set -e
shopt -s expand_aliases

# Start ssh-agent to handle the demo private-key
ssh-agent > /tmp/ssh-agent.sh
source /tmp/ssh-agent.sh

# Fix permissions for demo SSH key
chmod 600 $(pwd)/demo/assets/ssh-keys/id_ed25519

# Add the demo private-key
ssh-add $(pwd)/demo/assets/ssh-keys/id_ed25519

# Setup all aliases
alias sb1="ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -p 22001 t800@127.0.0.1 -A -tt -- "
alias sb2="ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -p 22002 t800@127.0.0.1 -A -tt -- "
alias sbscp="scp -O -S /tmp/.sbdemoscp "

echo "Generate an egress key"

sb1 self egress-key generate --algo ed25519 --size 256

echo "Get the generated public key"

export TRUSTED_KEY=$(sb1 self egress-keys list | sed -e 's/\x1b\[[0-9;]*m//g' | grep -e "^1:" | sed -e 's/^1: //')

echo "Push the key ${TRUSTED_KEY} to the docker examplevm container"

docker exec sbdemo_examplevm /bin/bash -c "echo '$TRUSTED_KEY' > /root/.ssh/authorized_keys"

echo "Grant personal access to examplevm"

sb1 self access add --host examplevm --port 22 --user root

# Wait for replication to be OK
sleep 5

echo "Execute a command on examplevm"

if ! sb2 root@examplevm -- echo "test"; then
   echo "Unable to connect to examplevm via sb2";
   exit 1;
fi

echo "Check that the session was recorded"

if ! sb2 self sessions list | grep -q "Session ID"; then
    echo "Session was not saved in logs on sb2";
    exit 1;
fi

# Wait for replication to be OK
sleep 5

echo "Check that the session was replicated"

if ! sb1 self sessions list | grep -q "Session ID"; then
    echo "Session was not replicated in logs on sb1";
    exit 1;
fi

echo "Check that we cannot connect to an unauthorized host"

if ! sb2 test@examplevm -- echo "test" | grep -q "user can't access the host"; then
    echo "User shouldn't be able to access examplevm as user test";
    exit 1;
fi

echo "Get the SCP program"

if ! sb1 scp --get-script > /tmp/.sbdemoscp; then
    echo "Unable to get the SCP program";
    exit 1;
fi

chmod +x /tmp/.sbdemoscp

# Adapt the script to our unsecured test environment
sed -i 's/ssh /ssh -o UserKnownHostsFile=\/dev\/null -o StrictHostKeyChecking=no /' /tmp/.sbdemoscp

echo "SCP recursively to the VM"

if ! sbscp -r ./docs root@examplevm:/root; then
    echo "Unable to SCP recursively to examplevm";
    exit 1;
fi

echo "Check that the content was indeed copied over to examplevm"

if ! sb1 root@examplevm -- cat docs/demo.md | grep -q "# Demo"; then
    echo "Unable to check if the content was copied over to examplevm"
    exit 1
fi

exit 0;