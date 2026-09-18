## Task 1
### TCP three-way handshake
```
10:43:40.603727 IP6 ::1.51058 > ::1.8080: Flags [S], seq 3844353658, ...   # SYN
10:43:40.606465 IP6 ::1.8080 > ::1.51058: Flags [S.], seq 3910441164, ack 3844353659, ...   # SYN/ACK
10:43:40.606500 IP6 ::1.51058 > ::1.8080: Flags [.], ack 1, ...   # ACK
```
### Http req
```
10:43:40.607730 IP6 ::1.51058 > ::1.8080: Flags [P.], seq 1:176, ...: HTTP: POST /notes HTTP/1.1
POST /notes HTTP/1.1
Host: localhost:8080
User-Agent: curl/8.22.0
Accept: */*
Content-Type: application/json
Content-Length: 39

{"title":"trace me","body":"in flight"}
```
### Http resp
```
10:43:40.616120 IP6 ::1.8080 > ::1.51058: Flags [P.], seq 1:207, ...: HTTP: HTTP/1.1 201 Created
HTTP/1.1 201 Created
Content-Type: application/json
Date: Fri, 18 Sep 2026 07:43:40 GMT
Content-Length: 93

{"id":5,"title":"trace me","body":"in flight","created_at":"2026-09-18T07:43:40.608316546Z"}
```
### Connection close
```
10:43:40.616484 IP6 ::1.51058 > ::1.8080: Flags [F.], seq 176, ack 207, ...   # FIN from client
10:43:40.616592 IP6 ::1.8080 > ::1.51058: Flags [F.], seq 207, ack 177, ...   # FIN from server
10:43:40.616694 IP6 ::1.51058 > ::1.8080: Flags [.], ack 208, ...  # one more final ACK
```


### All five command outputs from 1.3
```
 ss -tlnp   | grep :8080         # 1. what's listening?
ip route show                    # 2. routes from your host
mtr -rwc 5 localhost             # 3. reachability (loop on lo)
dig +short example.com @1.1.1.1  # 4. DNS works
journalctl --user -u quicknotes -n 20 || true  # 5. logs (if installed as service)

LISTEN 0      4096                *:8080            *:*    users:(("quicknotes",pid=32168,fd=3))
default via 172.27.64.1 dev eth0 proto kernel
172.27.64.0/20 dev eth0 proto kernel scope link src 172.27.67.40
Start: 2026-09-18T10:58:04+0300
HOST: tikhon        Loss%   Snt   Last   Avg  Best  Wrst StDev
  1.|-- ip6-localhost  0.0%     5    0.1   0.6   0.0   2.6   1.1
8.6.112.0
8.47.69.0
-- No entries --
```

### what would you check first if QuickNotes returned 502?
502 means the proxy or load balancer is up but the thing behind is not answering normally. so I would start from the backend, not the proxy. First `ss -tlnp | grep 8080` - is anything listening where proxy goes to? Then `curl -s localhost:8080/health` directly, without proxy, to find if it is a "proxy problem" or an "app problem". If the app is fine, check proxy upstream address + its error log (connection refused vs timeout tells a lot), and only after that dns/firewall. Last is app logs (`journalctl -u quicknotes`) for panics or restarts.

## Task 2

### 2.1
```
[1] 35649
2026/09/18 23:40:36 quicknotes listening on :8080 (notes loaded: 5)
[2] 35713
2026/09/18 23:40:36 quicknotes listening on :8080 (notes loaded: 5)
2026/09/18 23:40:36 listen: listen tcp :8080: bind: address already in use
exit status 1
[2]+  Done                       ADDR=:8080 go run . 2>&1 | tee /tmp/qn-broken.log
root       35649   28895 49 23:40 pts/2    00:00:01 go run .
```

### 2.2


`ps -ef \| grep quicknotes`
`root       35699   35649  0 23:40 pts/2    00:00:00 /root/.cache/go-build/42/42f644342d2ec31f4853932285c82d413b14012d331293f70c366a3bd031564e-d/quicknotes`
only one process, second one already died


`ss -tlnp \| grep 8080`
`LISTEN 0      4096                *:8080            *:*    users:(("quicknotes",pid=35699,fd=3))`
port is in use by pid 35699 (the first instance), thats what broke the second one


`curl -s -o /dev/null -w "%{http_code}\n" localhost:8080/health`
`200`
service is reachable and healthy - first instance answers. Do not look on the fact that green healthcheck! It doesnt prove the new deploy is running


`iptables -L -n -v`
```Chain INPUT (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination

Chain FORWARD (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination

Chain OUTPUT (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination
 ```
firewall is not blocking


`dig +short localhost`
`127.0.0.1` / `::1`
dns fine, nothing to fix here

To sum up: not network, not firewall, not dns - the second instance just couldnt bind on the same port


### 2.3
```
 kill  35699
2026/09/18 23:55:20 shutting down
[root@tikhon app]# ADDR=:8080 go run . &
sleep 1
curl -s http://localhost:8080/health

[2] 35897
[1]   Done                       ADDR=:8080 go run .
2026/09/18 23:55:31 quicknotes listening on :8080 (notes loaded: 5)
{"notes":5,"status":"ok"}
```

### root cause + postmortem
Root cause: `listen tcp :8080: bind: address already in use` - a previous instance still was holding the port when the new one started.

Mini-postmortem:
Nobody did anything "wrong" here: starting a new version while the old one is alive is a normal thing to do by hand. The systemic issue is that nothing is saying "only one user of the port is allowed". The new instance fails, but the old one keeps serving, so if you check `/health` returns 200 - the deploy looks successful while running old code. That is the dangerous part - failure is silent. The startup log also prints "listening" before the bind succeeds, which makes the logs lying about such important thing as successfull deploy.

What would prevent it: run the app under a supervisor (systemd unit / container such as docker) so start = stop old + start new instead of manual `go run &`; make the healthcheck return the build version (maybe hash)/pid so a deploy can check if we are on new ver; log "listening" only after `net.Listen` succeeds (such as all logs about success in operation). Ideally -- we nead an alert on a service exiting with non-zero return code right after start (restart loop) instead of thinking somebody reads the log (the same at CI).
