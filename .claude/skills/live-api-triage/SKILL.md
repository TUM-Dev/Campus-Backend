---
name: live-api-triage
description: Run the live API harness (client/livetest) against production or a local server, match its findings to issues and advisories already filed, and file bugs for the new ones. Use when asked to check the live API, triage harness results, or file bugs for API problems.
---

# Live API triage

The goal is to end with every real problem the harness finds tracked exactly once. Findings already filed get no duplicate, new findings get one bug each, and security problems are never disclosed publicly.

## 1. Run the harness

From `client/`:

```bash
go run ./livetest -json > "$SCRATCH/livetest.json"; echo "exit=$?"
```

- `$SCRATCH` is your session scratchpad directory. Keep the output out of the repo.
- Exit 1 only means that some check failed. Read the JSON either way.
- If the user names a target, use it: `-local` for a local server, or `-grpc host:port -rest https://host`. Only file issues for findings from production (`api.tum.app`). Local findings are for the user's own debugging.
- If `client/livetest` doesn't exist, the harness isn't on this branch. Say so and stop.

Flatten the result to one line per problem:

```bash
jq -r '.[] | .check as $c | .notes[]? | select(.status=="FAIL" or .status=="WARN") | "\(.status)\t\($c)\t\(.msg)"' "$SCRATCH/livetest.json"
```

Group the lines into **findings**, one per root cause, not one per note. For example:
- several `sweep/*` notes saying "falls back to Unimplemented" are one finding;
- the same 500 showing up in `cafeteria/*` and in `rest/error-mapping` is one finding.

## 2. Check each finding before filing

The harness is a heuristic. For each finding:

1. **Reproduce it once by hand** with `curl` (REST) or a targeted `go run ./livetest -run '<check>' -v`. Throwaway scripts go in the scratchpad.
2. **Look for the cause in the code** (`server/backend/…`) when that's cheap. A bug that points at a file and line is worth far more than a symptom. Clearly separate a confirmed cause from a hypothesis ("likely cause, not confirmed against prod").
3. **Drop false positives.** Examples:
   - an external link that is temporarily down;
   - a news source that is legitimately quiet;
   - no upcoming movies during the semester break (TUM lectures run roughly mid-Oct to mid-Feb and mid-Apr to mid-Jul).

   If a check is wrong rather than the API, fix the harness instead of filing a bug.
4. **Only state numbers you measured.** Every count, size, age or date in an issue must come from the harness output or your own reproduction. Never estimate.

## 3. Find what is already filed

Collect the candidates:

```bash
gh issue list --state open --limit 200 --json number,title,body,labels
gh issue list --state closed --limit 100 --search "sort:updated-desc" --json number,title,body,closedAt,stateReason
gh pr list --state open --limit 100 --json number,title,body,headRefName
gh api "repos/TUM-Dev/Campus-Backend/security-advisories?per_page=100" -q '.[] | {ghsa_id, state, summary}'
```

Issues filed by this skill carry a hidden marker in the body: `<!-- livetest: <check ids> -->`. Search for markers first:

```bash
gh issue list --state all --search '"livetest:" in:body' --json number,title,state,body
```

Then match the rest by meaning: endpoint, RPC name, error message, and the file the cause is in. A marker narrows the candidates but doesn't decide a match, because one check (e.g. `cafeteria/dishes-validation`) can report several unrelated bugs.

For each finding, decide:

| Situation | Action |
|---|---|
| An open issue or advisory covers it | Don't file. Add a comment only if you have **new** information, such as a root cause, a changed error, or a new affected endpoint. Never post "still happening" comments. If the issue has no marker, add one to its body (`gh issue edit --body-file`). |
| An open PR fixes it but isn't deployed yet | Don't file. Mention the PR in the report. |
| A closed issue covers it | Regression: open a new issue that references the closed one. If it was closed as "not planned", ask the user instead. |
| Nothing covers it | New finding (step 4). |

Also report filed issues whose checks now pass, as candidates for closing. Don't close them yourself.

## 4. Confirm, then file

Before creating anything, show the user one table: each finding, whether it is new or already filed (with the existing issue/PR/advisory number), and the proposed title. Wait for confirmation. The only exception is when the user has already said to file without asking in this conversation.

### Security findings: never in a public issue

`TUM-Dev/Campus-Backend` is public. File anything that exposes user data, credentials or uploads, bypasses access control, or gives an attacker leverage as a **draft security advisory**:

```bash
gh api -X POST repos/TUM-Dev/Campus-Backend/security-advisories --input advisory.json
```

`advisory.json` needs `summary`, `description`, `severity`, `cwe_ids`, and `vulnerabilities` with `{"package":{"ecosystem":"go","name":"github.com/TUM-Dev/Campus-Backend/server"}}`.

- Never download user content to prove the issue. Counting index entries or reading status codes is enough.
- Fixes for an advisory go through its temporary private fork (`POST …/security-advisories/<ghsa>/forks`), never through a public branch or PR.
- Low-impact hardening, like the public `/metrics` endpoint, can be a normal issue.
- If you're unsure whether something is security-relevant, treat it as security-relevant and ask.

### Bug issues

`gh issue create --label bug --title "[Bug] <what is broken>" --body-file -` with this body:

```markdown
<What happens, with the exact request and the response or status as measured.>

**Cause:** <file:line and why, or "Likely cause (not confirmed): …", or omit if unknown>

**Fix:** <concrete suggestion, if there is one>

Found by the live API harness (`client/livetest`, check `<check id>`) on <YYYY-MM-DD>.
<!-- livetest: <check id>[, <check id>…] -->
```

- One issue per root cause.
- Titles name the symptom users see, not the check.
- Write roles and code that start with `@` in backticks, because a bare `@name` notifies that user.

## 5. Report

End with a short summary:
- new issues and advisories, with links;
- findings that were already tracked, with their numbers;
- dropped false positives, with the reason;
- harness changes you made;
- filed issues whose checks now pass.
