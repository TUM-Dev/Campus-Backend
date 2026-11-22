package movie_parsers

import (
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

const sampleUnifilmHeilbronnHTML = `
<div class="anzeigebereich-film">
	<div class="film-showcase" data-id="880" data-sid="30801" data-cid="194">
		<div class="film-showcase-container">
			<figure class="image_container_placeholder">
				<img src="assets/images/f/Ghostbusters-Die_Geisterjaeger_DE-3175f26f.jpg" width="300" height="450" alt="">
			</figure>
			<div class="text_container">
				<h1 class="headline-h3">
					<span class="">Terminverschiebung: Ghostbusters (1984)</span>&nbsp;&nbsp;<span class="headline-normalcase">[OmU]</span>
				</h1>
				<h2 class="headline-h4 film-subtitel">Die Geisterjäger</h2>
				<ul class="film-info-termin">
					<li class="vertical-divider"><span class="film-info-text datum">Sonntag 18.10.2026</span></li>
					<li class="vertical-divider"><span class="film-info-text uhrzeit">18:00</span></li>
					<li><span class="film-info-text raum">NXTSpace - Bildungscampus</span></li>
				</ul>
				<span>Terminverschiebung: vom 15.10. auf den 18.10. um 18 Uhr</span>
				<p>Vier exzentrische Parapsychologen machen sich in New York selbstständig.</p>
				<p>Three eccentric parapsychologists set up a unique ghost removal service in New York City.</p>
				<ul class="film-info-filmdaten">
					<li class="vertical-divider">R: Ivan Reitman</li>
					<li class="vertical-divider">USA</li>
					<li class="vertical-divider">FSK 12</li>
					<li class="vertical-divider">Komödie</li>
					<li>105 Min.</li>
				</ul>
			</div>
		</div>
	</div>
	<div class="film-showcase" data-id="5228" data-sid="30847" data-cid="194">
		<div class="film-showcase-container">
			<figure class="image_container_placeholder">
				<img src="assets/images/c/good-luck-have-fun-dont-die-DE-dfc9ee8c.jpg" width="300" height="450" alt="">
			</figure>
			<div class="text_container">
				<h1 class="headline-h3">
					<span class="">Good Luck, Have Fun, Don't Die</span>&nbsp;&nbsp;<span class="headline-normalcase">[OmU]</span>
				</h1>
				<h2 class="headline-h4 film-subtitel"></h2>
				<ul class="film-info-termin">
					<li class="vertical-divider"><span class="film-info-text datum">Donnerstag 29.10.2026</span></li>
					<li class="vertical-divider"><span class="film-info-text uhrzeit">18:00</span></li>
					<li><span class="film-info-text raum">NXTSpace - Bildungscampus</span></li>
				</ul>
				<span></span>
				<p>Ein Mann aus der Zukunft taucht in einem Diner in Los Angeles auf.</p>
				<ul class="film-info-filmdaten">
					<li class="vertical-divider">R: Gore Verbinski</li>
					<li class="vertical-divider">DEU, USA</li>
					<li class="vertical-divider">FSK 16</li>
					<li class="vertical-divider">SciFi</li>
					<li>134 Min.</li>
				</ul>
			</div>
		</div>
	</div>
</div>`

func TestParseCampuskinoHeilbronnProgram(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sampleUnifilmHeilbronnHTML))
	require.NoError(t, err)

	screenings, err := ParseCampuskinoHeilbronnProgram(doc)
	require.NoError(t, err)
	require.Len(t, screenings, 2)

	first := screenings[0]
	require.Equal(t, "Terminverschiebung: Ghostbusters (1984) [OmU] – Die Geisterjäger", first.Title)
	require.Equal(t, "https://www.unifilm.de/studentenkinos/Heilbronn#sid=30801", first.Link)
	require.Equal(t, "2026-10-18T18:00:00+02:00", first.Info.ScreeningDate.Format(time.RFC3339))
	require.Equal(t, "NXTSpace - Bildungscampus, Heilbronn", first.Info.Location)
	require.Equal(t, "1984", first.Info.ReleaseYear.String)
	require.Equal(t, "Ivan Reitman", first.Info.Director.String)
	require.Equal(t, "Komödie", first.Info.Genre.String)
	require.Equal(t, "105 min", first.Info.Runtime.String)
	require.Equal(t, "https://www.unifilm.de/assets/images/f/Ghostbusters-Die_Geisterjaeger_DE-3175f26f.jpg", first.Info.ImageUrl)
	require.Contains(t, first.Info.ShortenedDescription, "<b>Terminverschiebung: vom 15.10. auf den 18.10. um 18 Uhr<b>")
	require.Contains(t, first.Info.ShortenedDescription, "Vier exzentrische Parapsychologen")
	require.Contains(t, first.Info.ShortenedDescription, "Three eccentric parapsychologists")

	second := screenings[1]
	require.Equal(t, "Good Luck, Have Fun, Don't Die [OmU]", second.Title)
	require.Equal(t, "https://www.unifilm.de/studentenkinos/Heilbronn#sid=30847", second.Link)
	require.Equal(t, "2026-10-29T18:00:00+01:00", second.Info.ScreeningDate.Format(time.RFC3339))
	require.Equal(t, "NXTSpace - Bildungscampus, Heilbronn", second.Info.Location)
	require.False(t, second.Info.ReleaseYear.Valid)
	require.Equal(t, "Gore Verbinski", second.Info.Director.String)
	require.Equal(t, "SciFi", second.Info.Genre.String)
	require.Equal(t, "134 min", second.Info.Runtime.String)
}

func TestCleanMovieTitleForSearch(t *testing.T) {
	require.Equal(t, "Ghostbusters", CleanMovieTitleForSearch("Terminverschiebung: Ghostbusters (1984) [OmU] – Die Geisterjäger"))
	require.Equal(t, "Good Luck, Have Fun, Don't Die", CleanMovieTitleForSearch("Good Luck, Have Fun, Don't Die [OmU]"))
	require.Equal(t, "The Drama", CleanMovieTitleForSearch("The Drama [OmU] – Noch mal auf Anfang"))
}
