# Lab 10

## Task 1

[release workflow](../.github/workflows/release.yml). Triggers on tag `v*`, builds `Dockerfile` (same one as Lab 6), pushes to
`ghcr.io/tikhonmakeev/devops-intro/quicknotes` with tags `v0.1.0` and `latest`. Job has `packages: write` + `contents: read`, nothing else.
All third-party actions are pinned by 40-char SHA (version in comment).

Not done yet (needs my GitHub account, not possible from the checkout):
- [ ] push tag `v0.1.0`, put the green run URL here
- [ ] flip package to public, `docker rmi` + `docker pull` output here

### a) OIDC vs GITHUB_TOKEN
For ghcr in the same repo `GITHUB_TOKEN` is enough. OIDC is for pushing to *another* system (AWS, GCP, Azure...): the job gets a short-lived token signed by GitHub and the cloud trusts it by claims (repo, branch, tag). So there is no long-lived cloud key stored in secrets that can leak.

### b) why :latest
Because humans and simple tooling (`docker run ...quicknotes`, a quick start in README, a new teammate) need a name which always means "the newest release". Production should still deploy the immutable `:v0.1.0`; `latest` is only a convenience pointer that moves with every release.

### c) packages: write only
Least privilege. If a step or a compromised action is running in the job, the token can only do what the job needs. With `write: all` the same token could push code or change a workflow in the repo, create releases or edit issues, so one bad action turns into repo takeover (like tj-actions in Lab 3).

## Task 2

Option A (Render). [cloud/render.md](../cloud/render.md) has the settings. I use *Existing Image* from ghcr, not build from repo, so Render runs exactly the image CI built (and I scanned in Lab 9). The deploy hook step is in the release workflow, the hook URL is in the `RENDER_DEPLOY_HOOK` secret.

Not done yet (needs my Render account, I will not invent numbers):
- [ ] service URL + `curl -v .../health`
- [ ] deploy log line with port
- [ ] warm p50 (5 requests), 3 cold starts, result of the note after spin down

### d) Render spin-down vs Cloud Run
Cloud Run is built around fast scale from zero (small container start, microVM/gVisor, image already cached near the node), so seconds or less. Render free tier optimizes cost: the instance is really stopped, the next request has to schedule it, pull/start the container and pass the health check, so about a minute. It is a free tier, nobody promises latency.

### e) PORT
Render is generic: it can run any image and it does not trust `EXPOSE` (it is only a hint and can be missing or wrong), so it gives the app a port by `PORT` (default 10000) and routes to it. I set `PORT=8080` and `ADDR=:8080` so both are equal on the first boot. If they differ, Render finds the real port, logs `New primary port detected` and restarts the deploy, around 45 s lost for every such deploy.

### f) existing image vs build in Render
Image from ghcr: the same bytes as tested/scanned, fast deploy (no build), and rollback is just another tag. But there is one more moving part (registry, public package, hook with the tag). Build from repo: fewer steps, but Render's build can differ from CI (cache, base image changed), is slower and I cannot say that it is the image from Lab 9. The note from step 5 is gone: free service has no disk, `/data` is in the container FS, which is thrown away on spin down, and the app starts again from `seed.json`.

## Bonus

[cloud/tunnel.md](../cloud/tunnel.md) has the commands. Measurements, the comparison table and checking from cellular are not done yet, so no table here.

### g) really cloud?
Render is the cloud by the usual meaning: someone else's datacenter, I don't own the machine. Tunnel is only a public door to my laptop. For a user it matters in availability: when I close the laptop the tunnel is dead, and the URL changes after restart. If the user only sees HTTPS and JSON, they can't tell the difference until it breaks.

### h) latency
Render: warm latency is mostly network distance to the region + TLS + Render's proxy; QuickNotes itself needs microseconds. Tunnel: the path is user -> Cloudflare edge -> tunnel to my laptop -> app, so my home uplink and the extra hop through the edge dominate.

### i) when tunnel
Right: home lab, on-prem service that must be reachable from outside without opening ports, a temporary URL for a stakeholder review or webhook testing. Never: real production with SLA, anything that depends on a laptop being on, or a thing that needs a stable address (quick tunnel URL is random).
