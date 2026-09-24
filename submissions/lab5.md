# Lab 5

virtualbox 7.2.18, vagrant 2.4.9, windows host

## Task 1

vagrantfile: [Vagrantfile](../Vagrantfile)


first 10 lines of vagrant up:
```
Bringing machine 'default' up with 'virtualbox' provider...
==> default: Box 'bento/ubuntu-24.04' could not be found. Attempting to find and install...
    default: Box Provider: virtualbox
    default: Box Version: 202510.26.0
==> default: Loading metadata for box 'bento/ubuntu-24.04'
    default: URL: https://vagrantcloud.com/api/v2/vagrant/bento/ubuntu-24.04
==> default: Adding box 'bento/ubuntu-24.04' (v202510.26.0) for provider: virtualbox (amd64)
    default: Downloading: https://vagrantcloud.com/bento/boxes/ubuntu-24.04/versions/202510.26.0/providers/virtualbox/amd64/vagrant.box
    default:
==> default: Successfully added box 'bento/ubuntu-24.04' (v202510.26.0) for 'virtualbox (amd64)'!
```
provisioning part:
```
==> default: Running provisioner: shell...
    default: Running: inline script
    default: /tmp/go.tgz: OK
    default: go version go1.24.13 linux/amd64
```

```
$ vagrant ssh -c 'go version'
go version go1.24.13 linux/amd64
quicknotes-vm

$ vagrant ssh -c 'cd ~/app && go build -o /tmp/qn && /tmp/qn &' &

$ vagrant ssh -c 'curl -s localhost:8080/health'
{"notes":4,"status":"ok"}

$ curl -s http://localhost:18080/health     # this fronm host
{"notes":4,"status":"ok"}
```

second `vagrant provision` just says `go 1.24.13 already installed`, so its idempotent

a) virtualbox shared folder (default for this provider). works out of the box with the bento box (guest additions are there), changes are visible both ways right away. rsync would be faster for the guest but one-way and needs rsync on windows. nfs doesnt work on windows host, smb needs credentials

b) nat, the default one. forwarded port is bound to 127.0.0.1 so only my machine can reach it. bridged puts the vm right into the lan (uni wifi etc) with its own ip, anyone there could hit it

c) shell, its just download + tar for go, no need for ansible for that

d) 1.24 means "whatever latest patch today", so two people get different versions on different days. exact patch + sha256 = same bytes for everyone

## Task 2

```
$ vagrant snapshot save clean-go1.24
==> default: Snapshotting the machine as 'clean-go1.24'...
==> default: Snapshot saved!

$ vagrant ssh -c 'sudo rm -rf /usr/local/go'

$ vagrant ssh -c 'go version'
bash: line 1: go: command not found

$ time vagrant snapshot restore clean-go1.24
==> default: Forcing shutdown of VM...
==> default: Restoring the snapshot 'clean-go1.24'...
==> default: Resuming suspended VM...
==> default: Machine booted and ready!

real    0m23.350s

$ vagrant ssh -c 'go version'
go version go1.24.13 linux/amd64

$ curl -s http://localhost:18080/health
{"notes":4,"status":"ok"}
```

deleted the snapshot after (`vagrant snapshot delete clean-go1.24`)

e) my snapshot is on the same disk as the vm. if that disk dies or the vm folder is deleted, I lose both. It helps undo changes inside the vm, like deleting Go here, but I still need a separate backup for losing the host disk.

f) copy-on-write = snapshot doesnt copy the disk, after it only changed blocks go to a new diff file. 10 snapshots dont take 10x disk, but every diff grows with changes and they pile up

g) with a long chain, reads may need to look through several diff files, and merging them on delete takes time. Later snapshots also depend on the earlier disk files. I would keep only a few restore points and fix the vagrantfile if provisioning breaks, otherwise I cant recreate the vm from scratch.

## Bonus

same laptop, same session, quicknotes running in both

| | vagrant vm | docker container |
|---|---:|---:|
| cold start | 1m17.9s | 0.67s |
| idle ram | 309 MiB used | 3.6 MiB |
| on-disk size | 2.4 GB | 15.7 MB |
| process count | 112 | 1 |

Here cold start means starting an existing stopped vm/container, as in the lab commands below. It does not include downloading the box/image, provisioning or building, and I did not time a separate wait for /health.

```
vm:
time vagrant halt && time vagrant up = 17.1s / 1m17.9s
free -h = Mem: 961Mi total, 309Mi used
ps -A --no-headers | wc -l = 112
du -sh ~/VirtualBox\ VMs/quicknotes-lab5 = 2.4G

docker:
docker stop qn5 && time docker start qn5 = 0.666s
docker stats --no-stream = 3.582MiB
docker top qn5 = 1 process (/quicknotes)
docker images --format '{{.Size}}' quicknotes:lab6 = 15.7MB
```

- I was surprised that the vm uses 309 MiB and runs 112 processes for the same small app that uses about 3.6 mb in docker. The vm also starts much slower
- I would use a vm when I need a separate kernel or a full OS. For this Go service the container is enough, and I dont need to run another Ubuntu around it
- These easily explains why containers became common for stateless microservices: less memory per service and faster starts when replacing instances
