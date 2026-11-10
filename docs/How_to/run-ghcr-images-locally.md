# Running published GHCR images locally

## Prerequisites

- A GitHub Personal Access Token (classic) with `read:packages`
  - You need the original key string that is shown when it is created
- The full Git commit SHA of the images you want to test

## 1. Authenticate with GHCR

Create a PAT if you don't have a string key for an existing one
- Create a Personal Access Token (classic) with `read:packages` at https://github.com/settings/tokens

Run this command (edit it to have your own GitHub username, and enter the PAT string when asked):

```bash
read -rsp "GitHub PAT: " GHCR_TOKEN
echo
printf '%s' "$GHCR_TOKEN" | docker login ghcr.io \
  -u YOUR_GITHUB_USERNAME \
  --password-stdin
unset GHCR_TOKEN
```

Docker should display `Login Succeeded`.

## 2. Select the image version

Find the image sha in GitHub Actions history, packages or on the deploy log on the Azure virtual machine.

Set `IMAGE_SHA` to the full commit SHA corresponding to the published images:

```bash
export IMAGE_SHA=<FULL_COMMIT_SHA>
```

Both the application and Nginx images use this same SHA

## 3. Pull the images

```bash
docker compose -p monk-ghcr -f compose.ghcr.yaml pull
```

## 4. Start the application

```bash
docker compose -p monk-ghcr -f compose.ghcr.yaml up -d --no-build
```

Open http://localhost:18081
(made like that to not mess up your work if you have something going on at port 8080)

## 5. View logs

```bash
docker compose -p monk-ghcr -f compose.ghcr.yaml logs -f
```

## 6. Stop the application

```bash
docker compose -p monk-ghcr -f compose.ghcr.yaml down
```

## Notes

- The images on GHCR are private, so each developer needs authorized GitHub credentials.
- The local database directory must be writable by the container user configured in `compose.ghcr.yaml`
