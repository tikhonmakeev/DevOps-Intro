### TCP three-way handshake
10:43:40.603727 IP6 ::1.51058 > ::1.8080: Flags [S], seq 3844353658, ...   # SYN
10:43:40.606465 IP6 ::1.8080 > ::1.51058: Flags [S.], seq 3910441164, ack 3844353659, ...   # SYN/ACK
10:43:40.606500 IP6 ::1.51058 > ::1.8080: Flags [.], ack 1, ...   # ACK

### Http req
10:43:40.607730 IP6 ::1.51058 > ::1.8080: Flags [P.], seq 1:176, ...: HTTP: POST /notes HTTP/1.1
POST /notes HTTP/1.1
Host: localhost:8080
User-Agent: curl/8.22.0
Accept: */*
Content-Type: application/json
Content-Length: 39

{"title":"trace me","body":"in flight"}

### Http resp
10:43:40.616120 IP6 ::1.8080 > ::1.51058: Flags [P.], seq 1:207, ...: HTTP: HTTP/1.1 201 Created
HTTP/1.1 201 Created
Content-Type: application/json
Date: Fri, 18 Sep 2026 07:43:40 GMT
Content-Length: 93

{"id":5,"title":"trace me","body":"in flight","created_at":"2026-09-18T07:43:40.608316546Z"}

### Connection close
10:43:40.616484 IP6 ::1.51058 > ::1.8080: Flags [F.], seq 176, ack 207, ...   # FIN from client
10:43:40.616592 IP6 ::1.8080 > ::1.51058: Flags [F.], seq 207, ack 177, ...   # FIN from server
10:43:40.616694 IP6 ::1.51058 > ::1.8080: Flags [.], ack 208, ...  # one more final ACK



### All five command outputs from 1.3
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