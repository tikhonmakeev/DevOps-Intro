# Lab 7

## Task 1

[playbook](../ansible/playbook.yaml), [inventory](../ansible/inventory.ini), [unit](../ansible/templates/quicknotes.service.j2).
Static binary and seed are in `ansible/files/`. Commands are in [ansible/README.md](../ansible/README.md).

### a) command vs modules
`command` runs a process, usually every time. It doesnt know the desired state. `file`, `copy`, `user`, `systemd_service` compare current state with what I asked for and only change the difference. Otherwise every deploy would restart a working service for no reason.

### b) handlers
Notify queues a handler only when the task changed. Multiple notifications still run it once. No change means no restart. I flush it before the start task so the new unit is loaded on the first deploy too.

### c) variables
For this small lab: play vars for paths/port, inventory group vars for VM-specific values, `-e` for a temporary experiment. Extra vars override play vars. If this became a role, put the usual values in role defaults so inventory can override them easily.

### d) facts
Not needed here. User, paths and unit are already known, no OS-dependent branches. `gather_facts: false` saves the setup module and transferring facts on every run; exact time depends on VM/SSH.

## Task 2

### e) changed=0
Directory already has the right owner/mode. Copy checks contents and metadata; template renders locally and compares the result and metadata with the remote file. Same result = nothing to write, nothing to notify.

### f) echo instead of template
Shell overwrites it every time and normally reports changed even for the same text, so a handler would restart every time. Also `ADDR=...` alone is not a valid unit. Quoting mistakes can corrupt it and shell doesnt fix owner/mode for me.

### g) check + diff
Check says a unit would change, diff shows *what*. I could notice `SEED_PATH` pointing to a wrong file or a port typo before applying it. Check mode alone doesnt show that detail. Neither is proof the service will start.

## Verification

Linux amd64 binary built with `CGO_ENABLED=0` (`file`: statically linked ELF). Ansible 10.7.0 installed in a temporary WSL venv.

```text
LANG=C.UTF-8 LC_ALL=C.UTF-8 ansible-playbook -i ansible/inventory.ini ansible/playbook.yaml --syntax-check

playbook: ansible/playbook.yaml
```

Deployment evidence is still pending.
No Lab 5 Vagrant VM is available in this checkout, so inventory must be filled from real `vagrant ssh-config` before deploying.
Need to capture first/second PLAY RECAP, variable-change handler output, check diff, health and seeded notes on that VM. No made-up recap here.

Bonus not attempted: proving git pull convergence needs publishing the solution branch first.
