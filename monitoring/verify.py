import base64
import json
import subprocess
import urllib.request
import urllib.parse
from datetime import datetime, timezone
from pathlib import Path


def get(url, auth=None):
    request = urllib.request.Request(url)
    if auth:
        request.add_header('Authorization', 'Basic ' + base64.b64encode(auth.encode()).decode())
    with urllib.request.urlopen(request) as response:
        return json.load(response)


env = json.loads(subprocess.check_output([
    'docker', 'inspect', '--format', '{{json .Config.Env}}', 'devops-lab8-grafana-1']))
env = dict(item.split('=', 1) for item in env)
auth = env['GF_SECURITY_ADMIN_USER'] + ':' + env['GF_SECURITY_ADMIN_PASSWORD']
dashboard = get('http://localhost:3000/api/dashboards/uid/quicknotes', auth)
evidence = {
    'captured_at': datetime.now(timezone.utc).isoformat(),
    'health': get('http://localhost:8080/health'),
    'targets': get('http://localhost:9090/api/v1/targets'),
    'dashboard_panels': [panel['title'] for panel in dashboard['dashboard']['panels']],
    'alerts': get('http://localhost:9090/api/v1/alerts'),
    'panel_values': {
        panel['title']: get('http://localhost:9090/api/v1/query?' + urllib.parse.urlencode({
            'query': panel['targets'][0]['expr']}))
        for panel in dashboard['dashboard']['panels']
    },
}
Path('submissions/lab8-evidence.json').write_text(json.dumps(evidence, indent=2) + '\n')
print(json.dumps(evidence, indent=2))
