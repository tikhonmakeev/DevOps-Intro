# Render setup (Option A)

Source: **Existing Image** `ghcr.io/tikhonmakeev/devops-intro/quicknotes:v0.1.0` (public package, no credentials).
I took the image from CI and not "build from repo", so Render runs exactly the bytes that CI pushed and I scanned.

| Setting | Value |
|---|---|
| Type | Web Service, instance type **Free** |
| Region | Frankfurt |
| Env `PORT` | `8080` |
| Env `ADDR` | `:8080` |
| Health check path | `/health` |

`PORT` and `ADDR` are set to the same port, so Render does not restart the deploy with `New primary port detected`.

Deploy hook: Settings -> Deploy Hook. The URL is stored only as the GitHub Actions secret `RENDER_DEPLOY_HOOK`
(`Settings -> Secrets and variables -> Actions`), the workflow adds `&imgURL=<url-encoded image:tag>`.

Free instance has no persistent disk, so notes created after the start are lost when the service spins down (`/data` is container FS).

## Teardown

Render dashboard -> service -> Settings -> Suspend (or Delete). The ghcr package stays, it costs nothing.
