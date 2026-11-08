package admin

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func testOIDC(t *testing.T) (Config, *http.Client) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	mux := http.NewServeMux()
	mux.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{map[string]string{
			"kty": "RSA", "kid": "test", "use": "sig", "alg": "RS256",
			"n": base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
			"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
		}}})
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		id, password, ok := r.BasicAuth()
		if !ok || id != "test-client" || password != "test-secret" || r.ParseForm() != nil || r.PostForm.Get("grant_type") != "authorization_code" {
			http.Error(w, "invalid client request", http.StatusBadRequest)
			return
		}
		payload, err := base64.RawURLEncoding.DecodeString(r.PostForm.Get("code"))
		var claims map[string]any
		if err != nil || json.Unmarshal(payload, &claims) != nil {
			http.Error(w, "invalid code", http.StatusBadRequest)
			return
		}
		challenge := sha256.Sum256([]byte(r.PostForm.Get("code_verifier")))
		if claims["_challenge"] != base64.RawURLEncoding.EncodeToString(challenge[:]) {
			http.Error(w, "invalid PKCE verifier", http.StatusBadRequest)
			return
		}
		delete(claims, "_challenge")
		if claims["_missing_id_token"] == true {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "test-access-token", "token_type": "Bearer"})
			return
		}
		badSignature := claims["_bad_signature"] == true
		delete(claims, "_bad_signature")
		payload, err = json.Marshal(claims)
		if err != nil {
			http.Error(w, "invalid claims", http.StatusBadRequest)
			return
		}
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","kid":"test"}`))
		signingInput := header + "." + base64.RawURLEncoding.EncodeToString(payload)
		hash := sha256.Sum256([]byte(signingInput))
		signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
		if err != nil {
			http.Error(w, "signing failed", http.StatusInternalServerError)
			return
		}
		if badSignature {
			signature[0] ^= 1
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "test-access-token", "token_type": "Bearer", "expires_in": 3600,
			"id_token": signingInput + "." + base64.RawURLEncoding.EncodeToString(signature),
		})
	})
	idp := httptest.NewTLSServer(mux)
	t.Cleanup(idp.Close)
	return Config{
		Enabled: true, BaseURL: "https://admin.example", ClientID: "test-client", ClientSecret: "test-secret",
		Issuer: "https://issuer.example/shibboleth", AuthorizationURL: idp.URL + "/authorize",
		TokenURL: idp.URL + "/token", JWKSURL: idp.URL + "/keys", AllowedSubjects: []string{"approved-subject"},
	}, idp.Client()
}

func adminRequest(handler http.Handler, method, path string, form url.Values, cookie *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, "https://admin.example"+path, strings.NewReader(form.Encode()))
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}

func findCookie(t *testing.T, response *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()
	for _, cookie := range response.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("missing cookie %s", name)
	return nil
}

func loginRequest(t *testing.T, handler http.Handler, cfg Config, alter func(map[string]any)) (*httptest.ResponseRecorder, string, *http.Cookie) {
	t.Helper()
	start := adminRequest(handler, http.MethodGet, "/admin/login", nil, nil)
	require.Equal(t, http.StatusSeeOther, start.Code)
	authURL, err := url.Parse(start.Header().Get("Location"))
	require.NoError(t, err)
	query := authURL.Query()
	require.Equal(t, "code", query.Get("response_type"))
	require.Equal(t, "S256", query.Get("code_challenge_method"))
	require.Equal(t, "openid", query.Get("scope"))
	require.Equal(t, cfg.BaseURL+"/admin/callback", query.Get("redirect_uri"))
	claims := map[string]any{
		"iss": cfg.Issuer, "aud": cfg.ClientID, "sub": "approved-subject",
		"exp": time.Now().Add(time.Hour).Unix(), "iat": time.Now().Unix(),
		"nonce": query.Get("nonce"), "_challenge": query.Get("code_challenge"),
	}
	if alter != nil {
		alter(claims)
	}
	code, err := json.Marshal(claims)
	require.NoError(t, err)
	callback := "/admin/callback?" + url.Values{
		"code": {base64.RawURLEncoding.EncodeToString(code)}, "state": {query.Get("state")},
	}.Encode()
	cookie := findCookie(t, start, loginCookie)
	return adminRequest(handler, http.MethodGet, callback, nil, cookie), callback, cookie
}

func TestAdminAuthenticationAndManagement(t *testing.T) {
	db := utils.SetupTestContainer(context.Background(), t)
	cfg, client := testOIDC(t)
	jobs := cron.New(db)
	handler, err := newHandler(db, jobs, "test-build", cfg, client)
	require.NoError(t, err)
	icon := model.File{Name: "icon.png", Path: "news/sources/", Downloaded: null.BoolFrom(true)}
	require.NoError(t, db.Create(&icon).Error)
	source := model.NewsSource{Title: "<script>unsafe</script>", URL: null.StringFrom("https://example.com/rss"), FileID: icon.File, Hook: null.StringFrom("newspread")}
	require.NoError(t, db.Omit("File").Create(&source).Error)
	job := model.Crontab{Type: null.StringFrom(cron.NewsType), ID: null.IntFrom(source.Source), Interval: 7200, LastRun: 1234567}
	require.NoError(t, db.Create(&job).Error)
	groupImage := model.File{Name: "group logo.png", Path: "student_club/", Downloaded: null.BoolFrom(true)}
	pendingImage := model.File{Name: "pending.png", Path: "student_club/", Downloaded: null.BoolFrom(false)}
	councilImage := model.File{Name: "council.svg", Path: "student_council/", Downloaded: null.BoolFrom(true)}
	require.NoError(t, db.Create(&groupImage).Error)
	require.NoError(t, db.Create(&pendingImage).Error)
	require.NoError(t, db.Create(&councilImage).Error)
	groupCollection := model.StudentClubCollection{Name: "University initiatives", Language: "English"}
	councilCollection := model.StudentCouncilCollection{Name: "Fachschaften", Language: "English"}
	require.NoError(t, db.Create(&groupCollection).Error)
	require.NoError(t, db.Create(&councilCollection).Error)
	require.NoError(t, db.Omit("Image", "StudentClubCollection").Create(&[]model.StudentClub{
		{Name: "Club <script>", Language: "English", ImageID: null.IntFrom(groupImage.File), ImageCaption: null.StringFrom("Credit <photo>"), StudentClubCollectionID: groupCollection.ID},
		{Name: "Missing image club", Language: "English", StudentClubCollectionID: groupCollection.ID},
		{Name: "Pending image club", Language: "English", ImageID: null.IntFrom(pendingImage.File), StudentClubCollectionID: groupCollection.ID},
	}).Error)
	council := model.StudentCouncil{
		Name: "Engineering council", Language: "English", Description: null.StringFrom("Council description"),
		LinkUrl: null.StringFrom("https://example.com/council"), ImageID: null.IntFrom(councilImage.File),
		ImageCaption: null.StringFrom("Council image credit"), StudentCouncilCollectionID: councilCollection.ID,
	}
	require.NoError(t, db.Omit("Image", "StudentCouncilCollection").Create(&council).Error)

	t.Run("identity headers do not authenticate", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/feeds", nil)
		req.Header.Set("Remote-User", "approved-subject")
		req.Header.Set("X-Device-Id", "approved-subject")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusSeeOther, response.Code)
		require.Equal(t, "/admin/login", response.Header().Get("Location"))
		require.Equal(t, "no-store", response.Header().Get("Cache-Control"))
		require.Contains(t, response.Header().Get("Content-Security-Policy"), "frame-ancestors 'none'")
		require.Equal(t, "same-origin", response.Header().Get("Referrer-Policy"))
	})

	t.Run("invalid identity tokens fail closed", func(t *testing.T) {
		cases := []struct {
			name   string
			alter  func(map[string]any)
			status int
		}{
			{"issuer", func(c map[string]any) { c["iss"] = "https://wrong.example" }, http.StatusUnauthorized},
			{"audience", func(c map[string]any) { c["aud"] = "wrong-client" }, http.StatusUnauthorized},
			{"expired", func(c map[string]any) { c["exp"] = time.Now().Add(-time.Minute).Unix() }, http.StatusUnauthorized},
			{"nonce", func(c map[string]any) { c["nonce"] = "wrong-nonce" }, http.StatusUnauthorized},
			{"unauthorized", func(c map[string]any) { c["sub"] = "other-user"; c["email"] = "admin@tum.de" }, http.StatusForbidden},
			{"case-sensitive subject", func(c map[string]any) { c["sub"] = "Approved-Subject" }, http.StatusForbidden},
			{"authorized party", func(c map[string]any) { c["azp"] = "other-client" }, http.StatusUnauthorized},
			{"multiple audiences", func(c map[string]any) { c["aud"] = []string{cfg.ClientID, "another-client"} }, http.StatusUnauthorized},
			{"signature", func(c map[string]any) { c["_bad_signature"] = true }, http.StatusUnauthorized},
			{"missing ID token", func(c map[string]any) { c["_missing_id_token"] = true }, http.StatusBadGateway},
		}
		for _, test := range cases {
			t.Run(test.name, func(t *testing.T) {
				response, _, _ := loginRequest(t, handler, cfg, test.alter)
				require.Equal(t, test.status, response.Code, response.Body.String())
				for _, cookie := range response.Result().Cookies() {
					require.NotEqual(t, sessionCookie, cookie.Name)
				}
			})
		}
	})

	response, callback, flowCookie := loginRequest(t, handler, cfg, nil)
	require.Equal(t, http.StatusSeeOther, response.Code, response.Body.String())
	cookie := findCookie(t, response, sessionCookie)
	require.True(t, cookie.Secure)
	require.True(t, cookie.HttpOnly)
	require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	require.Equal(t, "/", cookie.Path)
	require.Empty(t, cookie.Domain)
	require.Greater(t, cookie.MaxAge, 0)
	require.LessOrEqual(t, cookie.MaxAge, 3600)
	var session model.AdminSession
	require.NoError(t, db.First(&session, "token_hash = ?", digest(cookie.Value)).Error)
	require.NotEqual(t, cookie.Value, session.TokenHash)
	require.Equal(t, http.StatusBadRequest, adminRequest(handler, http.MethodGet, callback, nil, flowCookie).Code)
	mismatch := *flowCookie
	mismatch.Value = secret()
	require.Equal(t, http.StatusBadRequest, adminRequest(handler, http.MethodGet, callback, nil, &mismatch).Code)

	t.Run("dashboard and escaped read-only pages", func(t *testing.T) {
		for _, path := range []string{"/admin/", "/admin/feeds", "/admin/groups", "/admin/councils", "/admin/audit"} {
			response := adminRequest(handler, http.MethodGet, path, nil, cookie)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		}
		page := adminRequest(handler, http.MethodGet, "/admin/feeds", nil, cookie)
		require.NotContains(t, page.Body.String(), "<script>unsafe</script>")
		require.Contains(t, page.Body.String(), "&lt;script&gt;unsafe&lt;/script&gt;")
		require.Equal(t, http.StatusMethodNotAllowed, adminRequest(handler, http.MethodPost, "/admin/groups", nil, cookie).Code)
		require.Equal(t, http.StatusBadRequest, adminRequest(handler, http.MethodGet, "/admin/groups?page=-1", nil, cookie).Code)
	})
	t.Run("directory displays downloaded images and councils without leaking remote URLs", func(t *testing.T) {
		groups := adminRequest(handler, http.MethodGet, "/admin/groups", nil, cookie)
		require.Contains(t, groups.Body.String(), `src="/files/student_club/group%20logo.png"`)
		require.Contains(t, groups.Body.String(), `alt="Club &lt;script&gt;"`)
		require.Contains(t, groups.Body.String(), "Credit &lt;photo&gt;")
		require.Contains(t, groups.Body.String(), "No image")
		require.NotContains(t, groups.Body.String(), "pending.png")
		require.Contains(t, groups.Body.String(), `href="/admin/councils"`)
		require.Contains(t, groups.Header().Get("Content-Security-Policy"), "img-src 'self'")
		councils := adminRequest(handler, http.MethodGet, "/admin/councils", nil, cookie)
		require.Contains(t, councils.Body.String(), "Engineering council")
		require.Contains(t, councils.Body.String(), "Fachschaften")
		require.Contains(t, councils.Body.String(), "Council description")
		require.Contains(t, councils.Body.String(), `src="/files/student_council/council.svg"`)
		require.Contains(t, councils.Body.String(), "Council image credit")
		require.Equal(t, http.StatusSeeOther, adminRequest(handler, http.MethodGet, "/admin/councils", nil, nil).Code)
		require.Equal(t, http.StatusMethodNotAllowed, adminRequest(handler, http.MethodPost, "/admin/councils", nil, cookie).Code)
		require.Equal(t, http.StatusBadRequest, adminRequest(handler, http.MethodGet, "/admin/councils?page=-1", nil, cookie).Code)
	})
	t.Run("council pagination is independent of clubs", func(t *testing.T) {
		var extra []model.StudentCouncil
		for i := range pageSize {
			extra = append(extra, model.StudentCouncil{Name: fmt.Sprintf("Z council %02d", i), Language: "English", StudentCouncilCollectionID: councilCollection.ID})
		}
		require.NoError(t, db.Omit("Image", "StudentCouncilCollection").Create(&extra).Error)
		first := adminRequest(handler, http.MethodGet, "/admin/councils", nil, cookie)
		require.Contains(t, first.Body.String(), `href="?page=2"`)
		second := adminRequest(handler, http.MethodGet, "/admin/councils?page=2", nil, cookie)
		require.Equal(t, http.StatusOK, second.Code)
		require.Contains(t, second.Body.String(), "Z council 49")
		require.NotContains(t, second.Body.String(), "Engineering council")
		groups := adminRequest(handler, http.MethodGet, "/admin/groups", nil, cookie)
		require.NotContains(t, groups.Body.String(), `href="?page=2"`)
	})

	form := url.Values{
		"id": {strconv.FormatInt(source.Source, 10)}, "title": {"Updated feed"},
		"url": {"https://example.com/updated"}, "icon": {strconv.FormatInt(icon.File, 10)},
		"interval": {"3600"}, "enabled": {"on"}, "csrf": {session.CSRFToken},
	}
	t.Run("mutations require the session CSRF token and correct origin", func(t *testing.T) {
		invalid := url.Values{}
		for key, values := range form {
			invalid[key] = values
		}
		invalid.Set("csrf", "wrong")
		require.Equal(t, http.StatusForbidden, adminRequest(handler, http.MethodPost, "/admin/feeds", invalid, cookie).Code)
		require.Equal(t, http.StatusUnauthorized, adminRequest(handler, http.MethodPost, "/admin/feeds", form, nil).Code)
		request := httptest.NewRequest(http.MethodPost, "/admin/feeds", strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Set("Origin", "https://evil.example")
		request.AddCookie(cookie)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		require.Equal(t, http.StatusForbidden, response.Code)
	})
	t.Run("invalid feed configuration is rejected", func(t *testing.T) {
		for _, test := range []struct{ key, value string }{
			{"url", "http://example.com/rss"}, {"url", "https://user:secret@example.com/rss"},
			{"url", "javascript:alert(1)"}, {"icon", "0"}, {"interval", "4294967296"},
		} {
			invalid := url.Values{}
			for key, values := range form {
				invalid[key] = values
			}
			invalid.Set(test.key, test.value)
			require.Equal(t, http.StatusBadRequest, adminRequest(handler, http.MethodPost, "/admin/feeds", invalid, cookie).Code)
		}
	})
	t.Run("feed editing preserves hooks and last start time", func(t *testing.T) {
		response := adminRequest(handler, http.MethodPost, "/admin/feeds", form, cookie)
		require.Equal(t, http.StatusSeeOther, response.Code, response.Body.String())
		require.NoError(t, db.First(&source, source.Source).Error)
		require.NoError(t, db.First(&job, job.Cron).Error)
		require.Equal(t, "Updated feed", source.Title)
		require.Equal(t, "newspread", source.Hook.String)
		require.EqualValues(t, 1234567, job.LastRun)
		require.EqualValues(t, 3600, job.Interval)
		var audit model.AdminAudit
		require.NoError(t, db.Order("id DESC").First(&audit).Error)
		require.Equal(t, session.Subject, audit.Subject)
		require.Contains(t, audit.Details, "Updated feed")
		require.NotContains(t, audit.Details, session.CSRFToken)
		history := adminRequest(handler, http.MethodGet, "/admin/audit", nil, cookie)
		require.Equal(t, http.StatusOK, history.Code)
		require.Contains(t, history.Body.String(), "Previous value")
		require.Contains(t, history.Body.String(), "New value")
		require.Contains(t, history.Body.String(), "1 hour")
		require.Contains(t, history.Body.String(), "2 hours")
		require.Contains(t, history.Body.String(), "&lt;script&gt;unsafe&lt;/script&gt;")
		require.NotContains(t, history.Body.String(), "<pre>")
	})
	t.Run("malformed audit data is reported without hiding valid changes", func(t *testing.T) {
		entry := model.AdminAudit{Subject: session.Subject, Action: "feed.update", SourceID: source.Source, Details: "{invalid"}
		require.NoError(t, db.Create(&entry).Error)
		response := adminRequest(handler, http.MethodGet, "/admin/audit", nil, cookie)
		require.Equal(t, http.StatusOK, response.Code)
		require.Contains(t, response.Body.String(), "Details unavailable")
		require.Contains(t, response.Body.String(), "Updated feed")
		require.NotContains(t, response.Body.String(), "{invalid")
	})
	t.Run("ambiguous duplicate schedules cannot be changed", func(t *testing.T) {
		duplicate := model.Crontab{Type: null.StringFrom(cron.NewsType), ID: null.IntFrom(source.Source), Interval: 7200}
		require.NoError(t, db.Create(&duplicate).Error)
		t.Cleanup(func() { require.NoError(t, db.Delete(&duplicate).Error) })
		require.Equal(t, http.StatusConflict, adminRequest(handler, http.MethodPost, "/admin/feeds", form, cookie).Code)
	})
	t.Run("new disabled feed does not inherit the GORM enabled default", func(t *testing.T) {
		form.Set("id", "0")
		form.Set("title", "Disabled new feed")
		form.Del("enabled")
		require.Equal(t, http.StatusSeeOther, adminRequest(handler, http.MethodPost, "/admin/feeds", form, cookie).Code)
		var created model.NewsSource
		require.NoError(t, db.Where("title = ?", "Disabled new feed").First(&created).Error)
		var schedule model.Crontab
		require.NoError(t, db.Where("type = ? AND id = ?", cron.NewsType, created.Source).First(&schedule).Error)
		require.Zero(t, schedule.Interval)
		require.False(t, created.Hook.Valid)
	})
	t.Run("audit failure rolls back the source and schedule", func(t *testing.T) {
		var schedulesBefore int64
		require.NoError(t, db.Model(&model.Crontab{}).Count(&schedulesBefore).Error)
		require.NoError(t, db.Exec("CREATE TRIGGER reject_admin_audit BEFORE INSERT ON admin_audits FOR EACH ROW SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'audit unavailable'").Error)
		t.Cleanup(func() { require.NoError(t, db.Exec("DROP TRIGGER reject_admin_audit").Error) })
		form.Set("title", "Must roll back")
		response := adminRequest(handler, http.MethodPost, "/admin/feeds", form, cookie)
		require.Equal(t, http.StatusInternalServerError, response.Code)
		require.NotContains(t, response.Body.String(), "audit unavailable")
		var count int64
		require.NoError(t, db.Model(&model.NewsSource{}).Where("title = ?", "Must roll back").Count(&count).Error)
		require.Zero(t, count)
		require.NoError(t, db.Model(&model.Crontab{}).Count(&count).Error)
		require.Equal(t, schedulesBefore, count)
	})
	t.Run("revoked and expired sessions cannot view admin data", func(t *testing.T) {
		revoked := cfg
		revoked.AllowedSubjects = []string{"someone-else"}
		other, err := newHandler(db, jobs, "test", revoked, client)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, adminRequest(other, http.MethodGet, "/admin/feeds", nil, cookie).Code)
		require.NoError(t, db.Model(&session).UpdateColumn("expires_at", time.Now().Add(-time.Minute)).Error)
		require.Equal(t, http.StatusSeeOther, adminRequest(handler, http.MethodGet, "/admin/feeds", nil, cookie).Code)
		require.NoError(t, db.Model(&session).UpdateColumn("expires_at", time.Now().Add(time.Hour)).Error)
	})
	t.Run("logout revokes the database session", func(t *testing.T) {
		require.Equal(t, http.StatusForbidden, adminRequest(handler, http.MethodPost, "/admin/logout", url.Values{"csrf": {"wrong"}}, cookie).Code)
		response := adminRequest(handler, http.MethodPost, "/admin/logout", url.Values{"csrf": {session.CSRFToken}}, cookie)
		require.Equal(t, http.StatusSeeOther, response.Code)
		require.Equal(t, "/admin/signed-out", response.Header().Get("Location"))
		signedOut := adminRequest(handler, http.MethodGet, "/admin/signed-out", nil, nil)
		require.Equal(t, http.StatusOK, signedOut.Code)
		require.Contains(t, signedOut.Body.String(), "Sign in with TUM Login")
		require.Equal(t, -1, findCookie(t, response, sessionCookie).MaxAge)
		require.ErrorIs(t, db.First(&model.AdminSession{}, "token_hash = ?", session.TokenHash).Error, gorm.ErrRecordNotFound)
		require.Equal(t, http.StatusSeeOther, adminRequest(handler, http.MethodGet, "/admin/feeds", nil, cookie).Code)
	})
}

func TestAdminDisabledAndConfiguration(t *testing.T) {
	handler, err := New(nil, nil, "", Config{})
	require.NoError(t, err)
	for _, path := range []string{"/admin", "/admin/", "/admin/login", "/admin/feeds", "/admin/style.css"} {
		require.Equal(t, http.StatusNotFound, adminRequest(handler, http.MethodGet, path, nil, nil).Code)
	}
	_, err = New(nil, nil, "", Config{Enabled: true})
	require.Error(t, err)
	cfg, _ := testOIDC(t)
	require.NoError(t, cfg.validate())
	cfg.BaseURL = "http://admin.example"
	require.Error(t, cfg.validate())
	cfg.BaseURL = "https://admin.example/"
	require.Error(t, cfg.validate())
}
