# Admin guide: TUMonline status synchronization

This guide enables the backend to update four Firebase Remote Config parameters in
project `tca-backend-0001`. It assumes the production Docker Compose deployment from
this repository and a Linux host where an administrator has Docker/root access.

## What goes where

| Location | What the administrator adds |
|---|---|
| Google Cloud project `tca-backend-0001` | A dedicated service account, Remote Config permissions, and the enabled Remote Config API |
| Production server | The credential JSON file and a private Compose override containing the runtime settings |
| GitHub repository | Merge the status-sync PR after review; no new Actions secrets or variables are required by the current workflow |
| Firebase Remote Config | Check the four managed parameters and remove any conditional values on them before enabling sync |

The repository's [CI workflow](../.github/workflows/main.yml) tests the code and
publishes the Docker image on pushes to `main`. It does not deploy to the production
host. Adding a GitHub Actions secret alone will not make it available to the running
backend. Firebase credentials must never be included in Docker build arguments,
the image, or a committed file.

## 1. Create the Google service account

1. Open [Service Accounts for tca-backend-0001](https://console.cloud.google.com/iam-admin/serviceaccounts?project=tca-backend-0001).
2. Click **Create service account**.
3. Enter `tumonline-status-sync` as the name/ID and click **Create and continue**.
4. Grant **Firebase Remote Config Admin** (`roles/cloudconfig.admin`) on this project.
5. Click **Continue**, leave the optional user-access fields empty, then click **Done**.
6. Open the [Firebase Remote Config API](https://console.cloud.google.com/apis/library/firebaseremoteconfig.googleapis.com?project=tca-backend-0001)
   and click **Enable** if it is not already enabled.

The expected identity is
`tumonline-status-sync@tca-backend-0001.iam.gserviceaccount.com`.
The predefined role grants access to the project's Remote Config resources; it is
not restricted to the four parameter names. The backend limits its edits to those
four names. See [Firebase IAM roles](https://firebase.google.com/docs/projects/iam/roles-predefined-product#firebase_remote_config_roles).

If creating service accounts or granting roles is unavailable, a project IAM
administrator must perform these steps.

## 2. Provide credentials on the production server

For a Google Cloud runtime with an attached service account or an existing Workload
Identity Federation setup, use that identity with the role above. The backend uses
Application Default Credentials and supports those methods without a private key.

For the existing standalone Compose deployment, the following example uses a mounted
service-account key:

1. Open the newly created service account, then **Keys → Add key → Create new key → JSON**.
2. Download the JSON file and transfer it securely to the production administrator.
   Use an approved file-transfer/secret-management channel; do not paste it into a
   GitHub issue, PR, chat, or committed `.env` file.
3. On the production host, install it outside the repository. Replace the source path
   below with the securely transferred file:

   ```sh
   sudo install -d -m 0700 /root/credentials
   sudo install -o root -g 3000 -m 0640 /secure/upload/tumonline-status-sync.json /root/credentials/firebase-remote-config.json
   ```

The production backend runs as UID `1000`, GID `3000`, so the mounted key must be
readable by that identity. These permissions give its group read access. Keep the
host directory private and remove the transfer copy after installation.

If organization policy prohibits key creation, use an approved workload identity
setup rather than changing the policy just for this job.
See [Google's credential/key instructions](https://docs.cloud.google.com/iam/docs/keys-create-delete).

## 3. Configure the backend container

Create `/root/credentials/campus-backend.firebase.yml` on the production server:

```yaml
services:
  backend:
    environment:
      FIREBASE_PROJECT_ID: tca-backend-0001
      TUMONLINE_STATUS_SYNC_ENABLED: "true"
      GOOGLE_APPLICATION_CREDENTIALS: /run/secrets/firebase-remote-config.json
    volumes:
      - type: bind
        source: /root/credentials/firebase-remote-config.json
        target: /run/secrets/firebase-remote-config.json
        read_only: true
        bind:
          create_host_path: false
```

The credential path is the path **inside the container**. Setting only a host
environment variable does not mount the file or pass that variable to the backend.
The additional volume preserves the existing `/Storage` volume.

Include this override in every subsequent Compose deployment, restart/recreation,
and operational command. If a separate deployment tool manages the host, configure
the equivalent environment values and read-only mount there.

## 4. What to add in GitHub

For the workflow currently in this repository: **no Firebase Actions secret or
variable is required**. Review and merge [PR #697](https://github.com/TUM-Dev/Campus-Backend/pull/697),
then let the existing `CI` workflow publish
`ghcr.io/tum-dev/campus-backend/backend-server:main`. The existing `GITHUB_TOKEN`
used for publishing is supplied by GitHub; no Firebase credential is needed to
build or test the image.

If your organization deploys from another repository or external system, configure
the runtime credentials in that deployment system. Its setup cannot be inferred
from this repository.

If a GitHub deployment workflow is added later and genuinely needs the key, the
suggested contract is:

| GitHub setting | Name | Value |
|---|---|---|
| Environment secret | `FIREBASE_REMOTE_CONFIG_CREDENTIALS_JSON` | Complete service-account JSON, including its private key |
| Environment variable | `FIREBASE_PROJECT_ID` | `tca-backend-0001` |
| Environment variable | `TUMONLINE_STATUS_SYNC_ENABLED` | `true` |

Create these under **Repository Settings → Environments → production** (create the
environment if needed). The deployment job must declare `environment: production`,
explicitly read those settings, transfer the credential securely to the host, and
configure the mount above. **No current workflow consumes these suggested names;
creating them alone does nothing.** Prefer workload identity over storing a key
when adding a new deployment workflow. See [GitHub environment secrets](https://docs.github.com/en/actions/how-tos/write-workflows/choose-what-workflows-do/use-secrets#creating-secrets-for-an-environment).

## 5. Deploy and verify

Before enabling sync, review the existing defaults in
[Firebase Remote Config](https://console.firebase.google.com/project/tca-backend-0001/config).
Save any manually maintained message you want to retain: the four keys become owned
by the sync. Remove conditional values from those keys; the publisher rejects them
instead of silently overriding them.

From the production checkout, add the Firebase override to your existing Compose
invocation. The examples below use `.env.local` for the existing production database,
mail, and other settings. If your deployment uses another environment file or flags,
retain those instead; do not replace production configuration with the committed
local-development defaults.

```sh
sudo docker compose --env-file .env.local -f docker-compose.yml -f /root/credentials/campus-backend.firebase.yml pull backend
sudo docker compose --env-file .env.local -f docker-compose.yml -f /root/credentials/campus-backend.firebase.yml up -d --no-deps backend
sudo docker compose --env-file .env.local -f docker-compose.yml -f /root/credentials/campus-backend.firebase.yml exec backend sh -c 'test -r "$GOOGLE_APPLICATION_CREDENTIALS" && echo "Credential file is readable"'
sudo docker compose --env-file .env.local -f docker-compose.yml -f /root/credentials/campus-backend.firebase.yml logs --since 10m backend
```

The backend applies SQL migration `20261106000000_tumonline_status.sql` on startup.
It adds the cron type and schedules a run every 300 seconds. Allow roughly five
minutes plus the scheduler's polling/request time for a run.

Inspect these default parameters in Firebase:

| Parameter | Type | Expected behavior |
|---|---|---|
| `isTUMOnlineMaintenanceMode` | Boolean | True for a current TUMonline maintenance notice |
| `isTUMOnlineDegraded` | Boolean | True for another current TUMonline incident; false if maintenance takes priority |
| `tumOnlineStatusMessage_de` | String | German summary with a usable future end estimate, otherwise empty |
| `tumOnlineStatusMessage_en` | String | English summary with a usable future end estimate, otherwise empty |

Compare them with [TUM's current status notices](https://status.it.tum.de/).
No matching current notice means both flags are false and both messages empty.
All four values are published together. Unchanged values do not create a new
Firebase version, and successful runs currently do not emit a dedicated success
log. No new version alone is not proof of failure.

The app must fetch and activate Remote Config to show changes. Empty custom messages
allow its translated Boolean-based defaults. Debug builds disable these banners;
verify in an appropriate release build. Do not create a fake production incident
just to test delivery; use a separate test project if needed.

## 6. Troubleshooting and manual overrides

| Symptom | Check |
|---|---|
| No activity and no errors | Both runtime settings are present; enabled must be exactly `true`; the cron row exists with `interval=300`; allow a scheduled run |
| `remote config credentials` error | JSON is valid, mounted at the configured container path, and readable by UID/GID `1000:3000` |
| Firebase HTTP 403 | Correct project, enabled API, and service account has `roles/cloudconfig.admin` on that project |
| Conditional-value error | Remove conditional values from the four managed keys, or pause sync before maintaining them manually |
| TUM fetch/parser error | Existing Firebase values are intentionally retained; inspect source availability/markup |
| Firebase values correct, app unchanged | Client fetch/activation, real-time subscription, and build mode |

To make a manual override, change `TUMONLINE_STATUS_SYNC_ENABLED` to `"false"` in the
private override and recreate the backend using the same `up -d --no-deps backend`
command. Once the previous container has stopped, edit the four values in Firebase.
Disabling sync does not clear the last published values. Re-enabling it returns
ownership to the source on the next scheduled run.

To rotate a JSON key, install the new key at the host path and recreate the backend
so the bind mount opens the new file. Verify access and a scheduled run, then revoke
the old key in Google Cloud.
