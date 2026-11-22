# Testing news imports locally

Sentry is optional. The backend also writes logrus messages to the console;
with a running Compose stack, use `dory compose -f docker-compose.local.yml logs -f backend`
from the repository root (`docker compose` works with other Docker engines).

## Supported news and event feeds

The news importer (`newsCron` in `server/backend/cron/news.go`) supports both
standard RSS/Atom/JSON feeds (via `gofeed`) and TUM HTML event/news pages (TUM
Campus Heilbronn SolR pages, TUM Community events, and TYPO3 news lists):

| Source | Language | Type | URL |
| --- | --- | --- | --- |
| TUM News | DE | Atom | `https://www.tum.de/news.rss` |
| TUM News (EN) | EN | Atom | `https://www.tum.de/en/news.rss` |
| TUM Events | DE | RSS 2.0 | `https://www.tum.de/events.rss` |
| TUM Events (EN) | EN | RSS 2.0 | `https://www.tum.de/en/events.rss` |
| TUM Community | EN | HTML (`li.se-event`) | `https://www.community.tum.de/en/events/` |
| Studentische Vertretung (SV) | DE | Atom (`rss.xml`) | `https://www.sv.tum.de/sv/nachrichtenarchiv/rss.xml` |
| Student Representation (SV EN) | EN | Atom (`rss.xml`) | `https://www.sv.tum.de/en/sv/news-archive/rss.xml` |
| TUM Campus Heilbronn Events | DE | HTML (`c-events-solr-list__item`) | `https://chn.tum.de/de/news-events/events` |
| TUM Campus Heilbronn Events (EN) | EN | HTML (`c-events-solr-list__item`) | `https://chn.tum.de/news-events/events` |
| TUM Campus Heilbronn News | DE | HTML (`c-news-solr-list__item`) | `https://chn.tum.de/de/news-events/newsroom` |
| TUM Campus Heilbronn News (EN) | EN | HTML (`c-news-solr-list__item`) | `https://chn.tum.de/news-events/newsroom` |
| TUM School of CIT | DE | Atom (`rss.xml`) | `https://www.cit.tum.de/cit/aktuelles/rss.xml` |
| TUM School of CIT (EN) | EN | Atom (`rss.xml`) | `https://www.cit.tum.de/en/cit/news/rss.xml` |
| TUM School of Natural Sciences | DE | Atom (`rss.xml`) | `https://www.nat.tum.de/nat/aktuelles/rss.xml` |
| TUM School of Life Sciences | DE | Atom (`rss.xml`) | `https://www.ls.tum.de/ls/aktuelles/rss.xml` |

## Live feed and MySQL import

Run from `server/`. For Dory:

```bash
DOCKER_HOST="unix://$HOME/.dory/dory.sock" \
TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE=/var/run/docker.sock \
NEWS_FEED_URL=https://www.tum.de/news.rss \
NEWS_TEST_MYSQL=1 \
go test ./backend/cron -run '^TestNewsLiveFeed$' -count=1 -v
```

The test fetches the live feed (or supported HTML page), reports item counts and
the newest dated article, and checks whether its fields fit the production news
columns. It then replays that response through `newsCron` into a fresh database
with the production schema. Two imports must preserve all importable articles
without duplicating them or changing the source URL. The test database is
removed afterward. Image metadata is stored, but the separate image downloader
and the cron scheduler are not run.

For other Docker engines, omit the Dory-specific variables. Alternatively, set
`TEST_MYSQL_DSN='root:password@tcp(localhost:3306)/'` to use an existing MySQL server;
the test creates and drops its own randomly named database on that server.

Omit `NEWS_TEST_MYSQL=1` to inspect only the feed without a database. Change
`NEWS_FEED_URL` to check another source. Normal tests skip this test unless a
URL is supplied, so CI does not depend on live upstream services.

To check the database-to-REST path separately, run:

```bash
DOCKER_HOST="unix://$HOME/.dory/dory.sock" \
TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE=/var/run/docker.sock \
NEWS_TEST_MYSQL=1 \
go test ./backend -run '^TestNewsMySQLGateway$' -count=1 -v
```

This seeds a current article in a disposable MySQL database and verifies that
the generated REST gateway returns it, including source, date, and pagination
filters. It uses an in-process handler rather than production's network routing.

## Investigating production configuration

To inspect production configuration, these read-only queries show the source
URLs, cron intervals, and last attempted runs:

```sql
SELECT source, title, url, hook FROM news_sources ORDER BY source;
SELECT cron, id, type, `interval`, FROM_UNIXTIME(lastRun) AS last_attempt
FROM crontab WHERE type = 'news' ORDER BY id, cron;
SELECT src, COUNT(*) AS articles, MAX(date) AS newest_article,
       MAX(created) AS latest_import
FROM news GROUP BY src ORDER BY src;
```

`lastRun` records an attempt before execution, so a recent timestamp alone does
not prove a successful import. Console errors named `Cron job failed`,
`parsing rss`, and `Inserting new news failed` help locate failures.
