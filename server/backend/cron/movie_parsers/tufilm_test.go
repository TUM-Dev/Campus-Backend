package movie_parsers

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"

	"github.com/guregu/null"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

//go:embed test_data/oppenheimer.html
var htmlOppenheimer string

//go:embed test_data/supprise_film.html
var htmlSuppriseFilm string

//go:embed test_data/babylon.html
var htmlBabylon string

func TestOppenheimer(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlOppenheimer))
	require.NoError(t, err)
	header := "When the world changed forever"
	description := "In a hearing on his appeal against the revocation of his security clearance, physicist Julius Robert Oppenheimer looks back: on his beginnings, his private life and, above all, on the time when he is assigned scientific leadership of the Manhattan Project during World War II. At Los Alamos National Laboratory in New Mexico, he and his team are to develop a nuclear weapon under the supervision of Lt. Leslie Groves. Oppenheimer is proclaimed the &#34;father of the atomic bomb&#34;, but after his deadly invention is used with serious consequences in Hiroshima and Nagasaki, the just jubilant Oppenheimer is plunged into serious doubts.\n" +
		"In the further course of the film Lewis Strauss is to be confirmed as Secretary of Commerce in the cabinet of President Dwight D. Eisenhower. Soon it is also about his relations with Oppenheimer after the war, to whom he was superior as head of the Atomic Energy Authority. Thereby it concerns above all the reproaches to old connection to the communism which Oppenheimer is accused of."
	comment := "An intellectual thriller that is a masterpiece of image and word, of complex ideas made manifest, and all kept at a very human level to maintain constant immediacy. (Andrea Chase, Killer Movie Reviews)"
	expected := TuFilmWebsiteInformation{
		ImageUrl:             "https://www.tu-film.de/img/film/poster/.sized.Oppenheimer.jpg",
		ShortenedDescription: fmt.Sprintf("<b>%s<b>\n\n%s\n\n<i>%s<i>", header, description, comment),
		Director:             null.StringFrom("Christopher Nolan"),
		Actors:               null.StringFrom("Cillian Murphy, Emily Blunt, Matt Damon"),
		Runtime:              null.StringFrom("180 min"),
		ImdbID:               null.StringFrom("tt15398776"),
		ReleaseYear:          null.StringFrom("2023"),
		TrailerUrl:           null.String{},
	}
	info, err := parseWebsiteInformation(doc)
	require.NoError(t, err)
	require.Equal(t, &expected, info)
}

func TestBabylon(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBabylon))
	require.NoError(t, err)
	header := "Do you know where I can find some drugs?"
	description := "Its the 1920s, California, a chaotic world of parties and movie sets, crime and exuberance. Newcomers Nellie LaRoy (Margot Robbie) and Manny Torres (Diego Calva) will do everything to find their success. In front of the camera or behind it. When movies start getting lounder and sets get quiet, they and movie stars of old like Jack Conrad (Brad Pitt) will have to adapt or face extinction.\n" +
		"Dont let that story stop you tough, these three hours are filled with dance and drugs, music and montages, porn and poetry. You might not need that much of an attention span and you might not want it.\n" +
		"In the end, there is one thing Hollywood does best: Make movies about themselves. Did they forget to mention that these movies were made for nazi germany? Is all this crime and perversion really ok because the movies are just that good? Don&#39;t think about it too much, you will still enjoy it."
	comment := "Chazelle’s film commemorates the era’s hubris as it indulges in a bit of its own. This is how a world ends. Not with a whimper but a great deal of banging, baby. And vomiting. And snorting. (Irish Times)"
	expected := TuFilmWebsiteInformation{
		ImageUrl:             "https://www.tu-film.de/img/film/poster/.sized.Babylon.jpg",
		ShortenedDescription: fmt.Sprintf("<b>%s<b>\n\n%s\n\n<i>%s<i>", header, description, comment),
		Director:             null.StringFrom("Damien Chazelle"),
		Actors:               null.StringFrom("Brad Pitt, Margot Robbie, Jean Smart"),
		Runtime:              null.StringFrom("189 min"),
		ImdbID:               null.StringFrom("tt10640346"),
		ReleaseYear:          null.StringFrom("2022"),
		TrailerUrl:           null.StringFrom("https://www.youtube.com/watch?v=5muQK7CuFtY"),
	}
	info, err := parseWebsiteInformation(doc)
	require.NoError(t, err)
	require.Equal(t, &expected, info)
}

func TestSurpriseFilm(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlSuppriseFilm))
	require.NoError(t, err)
	header := "Überraschung!"
	description := "Seid gespannt und bereitet euch auf großes Kino vor – an diesem Termin spielen wir einen Überraschungsfilm für euch! Welcher Film das sein wird? Das wissen wir selbst noch nicht. Fest steht, die Kinobranche schläft nie und wir wollen uns nicht immer bereits ein halbes Jahr im Voraus festlegen, welche Filme bei uns laufen.\n" +
		"Seid also gespannt, welche Perle wir für euch aus dem Zauberhut ziehen! Der Film wird zwei Wochen davor über unsere üblichen Informationskanäle bekannt gegeben.\n" +
		"Habt ihr Meinungen dazu? Dann lasst es uns wissen und schreibt uns an ueberraschungsfilm@tu-film.de"
	expected := TuFilmWebsiteInformation{
		ImageUrl:             "https://www.tu-film.de/img/film/poster/.sized.berraschungsfilm.jpg",
		ShortenedDescription: fmt.Sprintf("<b>%s<b>\n\n%s", header, description),
	}
	info, err := parseWebsiteInformation(doc)
	require.NoError(t, err)
	require.Equal(t, &expected, info)
}

func TestNoExtration(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<html></html>"))
	require.NoError(t, err)
	expected := TuFilmWebsiteInformation{
		ImageUrl:             "https://www.tu-film.de/img/film/poster/.sized.berraschungsfilm.jpg",
		ShortenedDescription: "Surprise yourself",
	}
	info, err := parseWebsiteInformation(doc)
	require.NoError(t, err)
	require.Equal(t, &expected, info)
}
