# high error rate

## What this alert means
More than 5% of QuickNotes requests are returning 4xx/5xx, with the ratio above the limit for 5 minutes.

## Triage steps
1. Open localhost:9090/alerts and localhost:3000. Check when errors started and if traffic also changed.
2. Run `curl -i localhost:8080/health` and `curl -i localhost:8080/notes`. Health alone doesnt prove writes work.
3. Query `sum by (code) (rate(quicknotes_http_responses_by_code_total[5m]))` in Prometheus. Mostly 400 means bad requests; 500 needs app investigation.
4. Run `docker compose logs --tail=100 quicknotes`, check recent deploys and free disk with `docker compose exec quicknotes df -h /data`.

## Mitigations
- If a deploy started it, revert the bad app/config commit, rebuild and run `docker compose up -d --build quicknotes`. Keep the notes volume.
- If one client sends malformed JSON, stop/fix that client (including the lab traffic script). Check the request body before restarting a healthy app.
- If disk is full, free unrelated files or increase capacity. Dont delete the notes volume to make the alert green.

## Post-incident
Check errors recover and the alert becomes inactive. Record impact, timeline, root cause and one prevention action using [Lecture 1 postmortem guidance](../../lectures/lec1.md). No blaming the person who deployed.
