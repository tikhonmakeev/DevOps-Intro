# Cloudflare quick tunnel (bonus)

```bash
docker run -d --name qn -p 127.0.0.1:8080:8080 ghcr.io/tikhonmakeev/devops-intro/quicknotes:v0.1.0
cloudflared tunnel --url http://localhost:8080      # prints https://<random>.trycloudflare.com
hyperfine --runs 50 --warmup 3 'curl -s -o /dev/null https://<random>.trycloudflare.com/health'
```

Teardown: Ctrl+C in `cloudflared`, `docker rm -f qn`. The URL dies with the process.
