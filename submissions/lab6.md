# Lab 6

docker 29.5.3, compose v5.1.4

## Task 1

dockerfile: [dockerfile](../app/Dockerfile)

```
$ docker images quicknotes:lab6
IMAGE             ID             DISK USAGE   CONTENT SIZE   EXTRA
quicknotes:lab6   e56d9f060e9e       15.7MB         3.57MB
```

```
$ docker inspect quicknotes:lab6 | jq '.[0].Config'
{
  "User": "nonroot",
  "ExposedPorts": {
    "8080/tcp": {}
  },
  "Env": [
    "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
    "SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt",
    "ADDR=:8080",
    "DATA_PATH=/data/notes.json",
    "SEED_PATH=/seed.json"
  ],
  "Entrypoint": [
    "/quicknotes"
  ],
  "WorkingDir": "/home/nonroot"
}
```

the builder base btw:
```
$ docker images golang:1.24-alpine
IMAGE                ID             DISK USAGE   CONTENT SIZE   EXTRA
golang:1.24-alpine   8bee1901f1e5        395MB         83.5MB
```

docker run from 1.3 works, /health returns {"notes":4,"status":"ok"}, /notes returns the seed notes


a) Layer order matters because: you сhanged one .go file and rebuilt docker. 

if COPY . . first: go mod download reruns every time (0.2-0.3s). 

else if go.mod first: go mod download is CACHED. 

go build runs in both (~6s), source changed. Small gain here only because there are no deps in thi demo-proj, with real deps every code change would download all modules again

b) Without CGO_ENABLED=0 the binary is linked against real libc. distroless/static doesnt have it, so it fails with "no such file or directory"

c) Minimal debian image with ca-certificates, tzdata, /etc/passwd with nonroot (65532). No shell, no package manager, no libc. 

easy logic: less stuff = less CVEs. trivy found 0 in the OS part

d) -s drops the symbol table, -w drops DWARF debug info, -trimpath removes local paths from the binary. Binary from 9.4 Mb to 6.3. Cost: harder to debug with delve/gdb (on prod its fine).

## Task 2

compose: [compose.yaml](../compose.yaml)

```
$ docker compose up --build -d
$ curl -X POST -H 'Content-Type: application/json' -d '{"title":"durable","body":"survive a restart"}' http://localhost:8080/notes
{"id":5,"title":"durable","body":"survive a restart","created_at":"2026-09-24T20:03:36.983292588Z"}
$ curl -s http://localhost:8080/notes | grep -o '"title":"durable"'
"title":"durable"
$ docker compose down
$ docker compose up -d
$ curl -s http://localhost:8080/notes | grep -o '"title":"durable"'
"title":"durable"
$ docker compose down -v
 Volume devops-lab6_quicknotes-data Removed
$ docker compose up -d
$ curl -s http://localhost:8080/notes | grep -o '"title":"durable"'
$
```

e) The binary checks itself: `/quicknotes healthcheck` does GET /health and exits 0 or 1 (small new feature in main.go). Debug image would bring the shell back, sidecar is overkill, and just "process alive" doesnt catch real server problem and gives no healthy state for depends_on

f) Named volume is a separate docker object, down removes only containers and network. down -v, docker volume rm or prune delete it

g) It waits only for the container to start, not for the app to be ready, so the dependent service can get connection refused at startup. condition: service_healthy fixes it

## Bonus

```yaml
    read_only: true
    tmpfs:
      - /tmp
    cap_drop:
      - ALL
    security_opt:
      - no-new-privileges:true
```

```
$ docker inspect quicknotes:lab6 --format '{{ .Config.User }}'
nonroot

$ docker compose exec quicknotes sh
OCI runtime exec failed: exec failed: unable to start container process: exec: "sh": executable file not found in $PATH

$ docker inspect devops-lab6-quicknotes-1 --format '{{ .HostConfig.CapDrop }}'
[ALL]

$ docker cp test.txt devops-lab6-quicknotes-1:/etc/test
Error response from daemon: container rootfs is marked read-only

$ docker inspect devops-lab6-quicknotes-1 --format '{{ .HostConfig.SecurityOpt }}'
[no-new-privileges:true]
```

trivy 0.59.1:
```
quicknotes:lab6 (debian 13.7)
Total: 0 (HIGH: 0, CRITICAL: 0)

quicknotes (gobinary)
Total: 22 (HIGH: 22, CRITICAL: 0)
```
All 22 are Go 1.24 stdlib, 1.24 doesn't get fixes anymore. Best decision is update ver.

Most security per line: read_only: true. Even with RCE nothing can be changed or dropped on disk except /data and /tmp, and with no shell there is nothing to run anyway
