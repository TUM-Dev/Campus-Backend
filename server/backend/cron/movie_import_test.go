package cron

import (
	_ "embed"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron/movie_parsers"
	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type movieTestTransport func(*http.Request) (*http.Response, error)

func (f movieTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func movieTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mock.ExpectationsWereMet()); _ = conn.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)
	return db, mock
}

//go:embed movie_parsers/test_data/current_movie.html
var moviePage string

func TestMovieImportWithoutWorkingOMDb(t *testing.T) {
	for _, key := range []string{"", "invalid-key"} {
		t.Run("key="+key, func(t *testing.T) {
			t.Setenv("OMDB_API_KEY", key)
			db, mock := movieTestDB(t)
			old := webfetch.Client
			t.Cleanup(func() { webfetch.Client = old })
			omdbCalls := 0
			webfetch.Client = &http.Client{Transport: movieTestTransport(func(r *http.Request) (*http.Response, error) {
				body := ""
				switch {
				case r.URL.Host == "www.omdbapi.com":
					omdbCalls++
					body = `{"Response":"False","Error":"Invalid API key"}`
				case r.URL.Host == "www.unifilm.de":
					body = `<div class="kinoprogramm"></div>`
				case r.URL.Path == "/":
					body = `<nav><a href="/program/ws2026/">Programm</a></nav>`
				case r.URL.Path == "/program/ws2026/":
					body = `<div class="plan"><table><tr><td class="title"><a href="/program/ws2026/1506/">The Sheep Detectives</a></td></tr></table></div>`
				case r.URL.Path == "/program/ws2026/1506/":
					body = moviePage
				default:
					t.Fatalf("unexpected request: %s", r.URL)
				}
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
			})}
			mock.ExpectQuery("SELECT DISTINCT").WillReturnRows(sqlmock.NewRows([]string{"Link"}))
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT .* FROM `files`").WillReturnRows(sqlmock.NewRows([]string{"file"}))
			mock.ExpectExec("INSERT INTO `files`").WillReturnResult(sqlmock.NewResult(3, 1))
			mock.ExpectExec("INSERT INTO `movies`").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			require.NoError(t, New(db).movieCron())
			if key == "" {
				require.Zero(t, omdbCalls)
			} else {
				require.Equal(t, 1, omdbCalls)
			}
		})
	}
}

func TestMovieImageNameIgnoresQueryAndUsesDefaultExtension(t *testing.T) {
	require.True(t, strings.HasSuffix(movieImageName("https://example.com/poster.jpg?size=large"), ".jpg"))
	require.True(t, strings.HasSuffix(movieImageName("https://example.com/poster"), ".jpg"))
}

func TestPartialOMDbDataPreservesSourceMetadata(t *testing.T) {
	movie := model.Movie{Year: null.StringFrom("2026"), Director: null.StringFrom("Kyle Balda"), Description: "Source plot"}
	enrichMovie(&movie, &movie_parsers.OmdbResults{ReleaseYear: "N/A", Plot: "N/A", Runtime: "109 min"})
	require.Equal(t, "2026", movie.Year.String)
	require.Equal(t, "Kyle Balda", movie.Director.String)
	require.Equal(t, "Source plot", movie.Description)
	require.Equal(t, "109 min", movie.Runtime.String)
}

func TestCampuskinoHeilbronnMovieImport(t *testing.T) {
	t.Setenv("OMDB_API_KEY", "")
	db, mock := movieTestDB(t)
	old := webfetch.Client
	t.Cleanup(func() { webfetch.Client = old })
	webfetch.Client = &http.Client{Transport: movieTestTransport(func(r *http.Request) (*http.Response, error) {
		body := ""
		switch {
		case r.URL.Host == "www.unifilm.de":
			body = `<div class="anzeigebereich-film">
				<div class="film-showcase" data-id="5315" data-sid="30797" data-cid="194">
					<div class="film-showcase-container">
						<figure class="image_container_placeholder">
							<img src="assets/images/d/the-invite-DE-b16cac5d.jpg" width="300" height="450" alt="">
						</figure>
						<div class="text_container">
							<h1 class="headline-h3"><span>The Invite</span>&nbsp;&nbsp;<span class="headline-normalcase">[OmU]</span></h1>
							<h2 class="headline-h4 film-subtitel"></h2>
							<ul class="film-info-termin">
								<li class="vertical-divider"><span class="film-info-text datum">Donnerstag 26.11.2026</span></li>
								<li class="vertical-divider"><span class="film-info-text uhrzeit">18:00</span></li>
								<li><span class="film-info-text raum">NXTSpace - Bildungscampus</span></li>
							</ul>
							<p>Das Ehepaar Angela und Joe gerät in eine tiefe Krise.</p>
							<ul class="film-info-filmdaten">
								<li class="vertical-divider">R: Olivia Wilde</li>
								<li class="vertical-divider">USA</li>
								<li class="vertical-divider">FSK 12</li>
								<li class="vertical-divider">Drama</li>
								<li>107 Min.</li>
							</ul>
						</div>
					</div>
				</div>
			</div>`
		case r.URL.Path == "/":
			body = `<nav><a href="/program/ws2026/">Programm</a></nav>`
		case r.URL.Path == "/program/ws2026/":
			body = `<div class="plan"><table><tr><td class="title"><a href="/program/ws2026/1506/">The Sheep Detectives</a></td></tr></table></div>`
		default:
			t.Fatalf("unexpected request: %s", r.URL)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}

	// Mark TU Film link as already imported so only Campuskino Heilbronn is inserted
	mock.ExpectQuery("SELECT DISTINCT").WillReturnRows(sqlmock.NewRows([]string{"Link"}).AddRow("https://www.tu-film.de/program/ws2026/1506/"))
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .* FROM `files`").WillReturnRows(sqlmock.NewRows([]string{"file"}))
	mock.ExpectExec("INSERT INTO `files`").WillReturnResult(sqlmock.NewResult(4, 1))
	mock.ExpectExec("INSERT INTO `movies`").WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	require.NoError(t, New(db).movieCron())
}
