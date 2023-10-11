package cron

import (
	"strings"
	"testing"
)

func TestIMDBExtration(t *testing.T) {
	reader := strings.NewReader(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<body class="white">
<div class="container">
  <div id="main">
	<a name="top"></a>
	<div id="content">
<div class="widgets">
<div class="widget right info programm">
	<h3>Vorstellung</h3>
	<div>Do, 6. April 2023<br />
	um 20:00 Uhr<br />
	<a href="https://goo.gl/maps/WA54uV2bmsABudTi8">Hörsaal MW1801, Campus Garching</a></div>
<h1>6. April: Babylon<i> (Garching, OV)</i></h1>
<div class="widget" id="films">
<table>
<tr class="top">
<td colspan="2" class="title">
<h3><a href="http://www.imdb.com/title/tt10640346/" target="_blank" title="In der IMDB nachschlagen">Babylon</a> (Digital)</h3>
<h4>USA (2022)</h4>
</td>
</tr>
<tr class="bottom">
<td class="film">
<img src="/img/film/poster/.sized.Babylon.jpg" class="poster" alt="" /><br />
<a href="http://www.youtube.com/watch?v=5muQK7CuFtY">Zum Trailer</a><br />
<img src="/img/film/fsk_16.png" class="icon" alt="ab 16" />
<img src="/img/film/sf_dolbydigital.png" class="icon" alt="Dolby Digital" />
<img src="/img/film/pf_cinemascope.png" class="icon" alt="CinemaScope" />
<br /><i>Regie: </i>Damien Chazelle
<br /><i>Schauspieler: </i>Brad Pitt, Margot Robbie, Jean Smart
<br />189 Minuten
</td>
<td class="text">
<div class="teaser">Do you know where I can find some drugs?</div>
<div class="description">
<p>Its the 1920s, California, a chaotic world of parties and movie sets, crime and exuberance. Newcomers Nellie LaRoy (Margot Robbie) and Manny Torres (Diego Calva) will do everything to find their success. In front of the camera or behind it. When movies
<img src="/img/film/scenes/.thumb.Babylon-1.jpg" title="Szene aus Babylon" class="scene first" 0="0" alt="" /><img src="/img/film/scenes/.thumb.Babylon-2.jpg" title="Szene aus Babylon" class="scene" 0="0" alt="" /><img src="/img/film/scenes/.thumb.Babylon.jpg" title="Szene aus Babylon" class="scene" 0="0" alt="" /> start getting lounder and sets get quiet, they and movie stars of old like Jack Conrad (Brad Pitt) will have to adapt or face extinction.</p>
<p>

Dont let that story stop you tough, these three hours are filled with dance and drugs, music and montages, porn and poetry. You might not need that much of an attention span and you might not want it.
</p><p>
In the end, there is one thing Hollywood does best: Make movies about themselves. Did they forget to mention that these movies were made for nazi germany? Is all this crime and perversion really ok because the movies are just that good? Don't think about it too much, you will still enjoy it.</p>
</div>
<div class="comment">Chazelle’s film commemorates the era’s hubris as it indulges in a bit of its own. This is how a world ends. Not with a whimper but a great deal of banging, baby. And vomiting. And snorting. (Irish Times)</div>
</td>
</tr>
</table></div>
</div>
</div>
</body>
</html>
`)
	imdbID, err := parseImdbIDFromReader(reader)
	if err != nil {
		t.Error(err)
	}
	if imdbID != "tt10640346" {
		t.Error("imdbID is not correct")
	}
}
