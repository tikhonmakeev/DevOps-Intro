# Lab 8

## Task 1

[compose](../compose.yaml), [prometheus](../monitoring/prometheus/prometheus.yml), [datasource](../monitoring/grafana/provisioning/datasources/datasource.yml), [provider](../monitoring/grafana/provisioning/dashboards/dashboard.yml), [dashboard](../monitoring/grafana/dashboards/golden-signals.json).

### a) pull vs push
Prometheus calls QuickNotes, so QuickNotes must be reachable from Prometheus. If connection fails, `up` becomes 0 and fresh app samples stop. Old graph points dont prove the app is alive now.

### b) interval
5s means more samples, storage and scrape overhead; a rate window should still contain several samples. 5m is too sparse for a 5m rate window, short spikes disappear and alerts react slowly. 15s with a 5m window is enough for this app.

### c) rate / irate / delta
`rate` for traffic: requests are a counter, rate handles resets and averages over the window. `irate` uses only the last two points and jumps around. `delta` is for gauges, it doesnt handle counter resets correctly.

### d) provisioning
Fresh stack loads the same datasource/dashboard. No remembering what I clicked last time, and config changes are visible in git.

No duration histogram in this app. Used the request-rate proxy explicitly allowed by lab, but labelled it as a proxy: requests/s is not actual latency. Notes gauge is also just a storage growth signal, not a capacity percentage.

## Task 2

[alert rule](../monitoring/prometheus/alerts.yml), [runbook](../docs/runbook/high-error-rate.md).
The rule uses 4xx + 5xx, `for: 5m`, `severity: page` and a repo runbook link.

### e) why 5 minutes
One bad request is normal, waking somebody for it isnt. Sustained ratio means an ongoing problem. The 5m rate window plus 5m pending time also means firing can take longer than 5m after starting bad traffic.

### f) cause alert
CPU >80% would be a cause/resource alert. Busy CPU can still serve users fine; errors are directly visible to them. Also errors can happen with CPU almost idle (bad config, wrong request body).

### g) noise
As a practical target I would review the rule if >10% of pages had no user impact over the last month. Thats my threshold, not a universal constant. Track false pages / all pages, plus missed incidents so tuning doesnt just hide everything.

## Run

```sh
export GRAFANA_PASSWORD=$(openssl rand -hex 18)
docker compose up -d --build
python monitoring/traffic.py --seconds 40
curl -s localhost:9090/api/v1/targets
python monitoring/traffic.py --errors --seconds 660
curl -s localhost:9090/api/v1/alerts
```

Grafana is on localhost:3000, user `tikhon`, password from the environment. No password committed.
Traffic script sends 4 good requests and 1 bad per second in error mode. Stop it after the demo; dont leave intentional errors running.
## Checks

`docker compose up -d --build` starts all 3 services, QuickNotes is healthy and `/health` returns `{"notes":4,"status":"ok"}`.
Prometheus target is `up`; Grafana API confirms the provisioned dashboard has 4 panels.

```text
promtool check config /etc/prometheus/prometheus.yml
SUCCESS: 1 rule files found
SUCCESS: /etc/prometheus/prometheus.yml is valid prometheus config file syntax
Checking /etc/prometheus/alerts.yml
SUCCESS: 1 rules found

promtool test rules alerts.test.yml
SUCCESS
```

The rule test checks sustained errors fire after the pending period and a single burst does not page.
Real stack outputs are in [lab8-evidence.json](lab8-evidence.json), captured by `python monitoring/verify.py`.
At 01:55:31 MSK on October 2 the real alert was `firing`, with error ratio about 17%. Pending started at 01:50:27 MSK. Stopped the error generator after capturing it.
Dashboard screenshot still needed: open localhost:3000 → Dashboards → QuickNotes while traffic runs.

Bonus not attempted (needs a public URL and two external regions for 30 minutes).
