# high error rate

## What this alert means
More than 5% of quicknotes requests are returning 4xx/5xx. above the limit for 5 minutes

## Triage steps
1. Open localhost:9090/alerts and localhost:3000. Check when errors started and if traffic also changed
2. Run `curl -i localhost:8080/health` and `curl -i localhost:8080/notes`. Health alone doesnt prove writes work
3. Query `sum by (code) (rate(quicknotes_http_responses_by_code_total[5m]))` in Prometheus. If mostly 400 - means bad requests and doesnt mean a bug; 500 needs to look into app logs
4. Run `docker compose logs --tail=100 quicknotes` and check recent deploys. Use `docker system df -v` to check docker disk usage

## Mitigations
- If a deploy started it, revert the bad app/config commit, rebuild and run `docker compose up -d --build quicknotes`.
- If one client sends bad JSON, stop and fix that client. Check request body before restarting - app may be ok
- If disk is full free bad files or increase space. Dont delete the notes volume!

## Post-incident
Check errors recover and the alert becomes inactive. Record impact, timeline, root cause and one prevention action. No blaming the person who deployed.
