# run

From WSL/Linux, in the repo root:

```sh
cd app
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../ansible/files/quicknotes .
cd ..
vagrant ssh-config
# update inventory.ini with the actual host, port and IdentityFile
ansible-playbook -i ansible/inventory.ini ansible/playbook.yaml --check --diff
ansible-playbook -i ansible/inventory.ini ansible/playbook.yaml
ansible-playbook -i ansible/inventory.ini ansible/playbook.yaml
curl -s localhost:18080/health
curl -s localhost:18080/notes
ansible-playbook -i ansible/inventory.ini ansible/playbook.yaml -e listen_addr=:9090
ansible-playbook -i ansible/inventory.ini ansible/playbook.yaml -e listen_addr=:9091 --check --diff
# restore the port forwarded by Vagrant
ansible-playbook -i ansible/inventory.ini ansible/playbook.yaml
```

The key path in inventory is a default, not output from a running VM.
On a completely fresh VM check mode may fail after predicting the new user: it does not actually create that user for later tasks.
