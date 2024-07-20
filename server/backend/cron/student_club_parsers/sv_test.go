package student_club_parsers

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
)

//go:embed sv_studentclubs_07-24.html
var htmlStudentClubs string

func TestParseStudentClubs(t *testing.T) {
	expectedClubs := []SVStudentClub{
		{Name: "AG JLC", Collection: "Akademisch", Description: null.StringFrom("Die AG besteht aus Studierenden, Promovierenden und Absolv\u00adent*in\u00adnen der Lebensmittelchemie. Wir organisieren Treffen &amp; Exkursionen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/AG_JLC.jpg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.ag-jlc.de/")},
		{Name: "Applied Mathematics Student Club Munich", Collection: "Akademisch"},
		{Name: "Debattierclub München", Collection: "Akademisch", Description: null.StringFrom("Seit 2001 bieten wir einen Ort für freundschaftliche Debatten und nehmen an deutschen und internationalen Wettkämpfen teil."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Debattierclub_Muenchen.png"), ImageCaption: null.StringFrom("Bild: Debattierclub München"), LinkUrl: null.StringFrom("https://www.debattierclub-muenchen.de/")},
		{Name: "Desorientierungstage", Collection: "Akademisch", Description: null.StringFrom("Unsere Veranstaltungen führen zu Diskussionen, sprengen Denkmuster und ermöglichen Bekanntschaften außerhalb der eigenen Bubble."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Desorientierungstage.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.desorientierungstage.de/")},
		{Name: "Edventure@TUM", Collection: "Akademisch", Description: null.StringFrom("Tauche in reale Bildungsprojekte ein, denke unternehmerisch und löse kollaborativ Herausforderungen im Bildungsbereich."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Edventure.png"), ImageCaption: null.StringFrom("Bild: Razin Abdullah / TUM"), LinkUrl: null.StringFrom("https://edventuretum.wordpress.com/")},
		{Name: "EESTEC LC Munich", Collection: "Akademisch", Description: null.StringFrom("Unser Ziel als EESTEC ist es, internationale Kontakte zu knüpfen und den Gedankenaustausch unter MINT-Studenten zu fördern."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/EESTEC_LC.svg"), ImageCaption: null.StringFrom("Bild: EESTEC LC Munich"), LinkUrl: null.StringFrom("https://eestec-muc.de/")},
		{Name: "EMSA", Collection: "Akademisch", Description: null.StringFrom("Als 1991 in Brüssel gegründete NGO vertritt EMSA die Interessen und Meinungen euro\u00adpäischer Medizinstudierenden."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/EMSA.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.sv.tum.de/med/projekte/emsa/")},
		{Name: "Enactus München", Collection: "Akademisch", Description: null.StringFrom("Wir fokussieren uns auf soziales Unternehmertum. Dazu arbeiten wir mit 6 Projekten weltweit an den Sustainable Development Goals."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Enactus.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://enactus-muenchen.de/")},
		{Name: "ISPE SC München", Collection: "Akademisch", Description: null.StringFrom("Wir bieten Studenten und Berufseinsteigern die Möglichkeit ihre Netzwerke zu erweitern und an Events teilzunehmen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/ISPE_SC_Muenchen.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://ispe-dach.org/student-chapter-muenchen/")},
		{Name: "MTP München", Collection: "Akademisch", Description: null.StringFrom("Durch Treffen mit Unternehmens\u00adreferenten und Beratungsprojekte gestalten wir ein Bildungsnetzwerk für Marketing und Unternehmertum."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/MTP.png"), ImageCaption: null.StringFrom("Bild: Marketing zwischen Theorie und Praxis"), LinkUrl: null.StringFrom("https://www.mtp.org/muenchen")},
		{Name: "Philosophia Munich", Collection: "Akademisch", Description: null.StringFrom("Durch Diskussion moderner Philosophie fördern wir das Verständnis für Ethik, Gesellschaft und dem Universum."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Philosophia_Mu__nchen.jpg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.philosophiamunich.org/")},
		{Name: "PushQuantum", Collection: "Akademisch"},
		{Name: "TEDxTUM", Collection: "Akademisch", Description: null.StringFrom("Unter dem Motto Ideas Worth Spreading haben wir das Ziel, neu\u00adgierige Geister zu inspirieren und zu ermutigen, Ideen anderer zu nutzen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TEDxTUM.png"), ImageCaption: null.StringFrom("Bild: Impact x Ideas e.V."), LinkUrl: null.StringFrom("https://www.tedxtum.com")},
		{Name: "TUM Case Club", Collection: "Akademisch", Description: null.StringFrom("Bei uns dreht sich alles um Case Studies: Workshops mit erfahrenen Vortragenden, inspirier\u00adendes Netzwerk und eigene Wettbewerbe."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TUM_Case_Club.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.tumcaseclub.de/")},
		{Name: "TUM-ABES", Collection: "Akademisch", Description: null.StringFrom("TUM-ABES organisiert Exkursionen zu Firmen, Startups und Lehrstühlen in der Biotechnologie und zeigt Studierenden Karrieremöglichkeiten."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TUM_ABES.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://m.facebook.com/tumabes/")},
		{Name: "TUMuchData", Collection: "Akademisch"},
		{Name: "BSH München", Collection: "Interessensgruppen", Description: null.StringFrom("Wir fördern den Austausch über sicherheitspolitische Themen durch vielfältige Veranstal\u00adtungen &amp; Dis\u00adkussionen über internationale Politik."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Bundesverband_Sicherheitspolitik_an_den_Hochschulen.jpg"), ImageCaption: null.StringFrom("Bild: BSH München"), LinkUrl: null.StringFrom("https://muenchen.sicherheitspolitik.de/aktuelles")},
		{Name: "EA München", Collection: "Interessensgruppen", Description: null.StringFrom("Wir versuchen Wege zu finden soviel Gutes zu tun wie möglich. Dabei verlassen wir uns auf Evidenz und sorgfältige Analysen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Effektiver_Altruismus.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://bit.ly/ea_munich_signup")},
		{Name: "Init. Stoppt Studiengebühren", Collection: "Interessensgruppen", Description: null.StringFrom("Gleiche Bildungschancen für alle: Wir stellen uns gegen Gebühren für internationale Studierende und stehen für eine weltoffene TUM ein."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Initiative_Stoppt_Studiengebuehren.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://stopptstudiengebuehren.vercel.app/")},
		{Name: "Munich Animal Rights Club", Collection: "Interessensgruppen"},
		{Name: "PAN UG TU München", Collection: "Interessensgruppen", Description: null.StringFrom("Unser Ziel ist es, den Stellenwert der Ernährungsmedizin zu verbessern, ihre Rolle in Prävention/Kuration zu betonen und Projekte zu erarbeiten."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/PAN.png"), ImageCaption: null.StringFrom("Bild: PAN International"), LinkUrl: null.StringFrom("https://pan-int.org/de/university-groups/")},
		{Name: "Responsible Tech Hub", Collection: "Interessensgruppen", Description: null.StringFrom("Wir möchten die Debatte um technische Verantwortung entmystifizieren und einen Raum für integrative Co-Kreation schaffen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Responsible_Technology_Hub.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.responsibletechhub.com/")},
		{Name: "School of Circularity", Collection: "Interessensgruppen"},
		{Name: "School of Transformation", Collection: "Interessensgruppen", Description: null.StringFrom("SOFT setzt sich als intersektional feminist\u00adisches Kollektiv mit feminist\u00adischer Praxis in Bildung, Arbeit und Selbstorganisation auseinander."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/SOFT.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.schooloftransformation.eu/")},
		{Name: "Studentischer Automobilverband München", Collection: "Interessensgruppen"},
		{Name: "VACCtion", Collection: "Interessensgruppen"},
		{Name: "Women in CS @ TUM", Collection: "Interessensgruppen", Description: null.StringFrom("Als Initiative an der TUM-Informatik setzen wir uns für gleichberech\u00adtigte Beteiligung von Frauen und anderen unter\u00adrepräsentierten Gruppen ein."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/IFF.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://womenincstum.github.io/")},
		{Name: "100 Voices - One Planet", Collection: "Interessensgruppen", Description: null.StringFrom("Wir sammeln Stimmen aus den 100 am stärksten vom Klimawandel betroffenen Ländern und vereinen sie in einer Dokumentation."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/100_Voices_One_Planet.png"), ImageCaption: null.StringFrom("Bild: 100VOP"), LinkUrl: null.StringFrom("https://www.100vop.org")},
		{Name: "LBV-HSG Freising", Collection: "Interessensgruppen", Description: null.StringFrom("Wir möchten gemeinsam die Natur erkunden, Artenkenntnisse vertiefen und dabei eine gute Zeit verbringen – mit Exkursionen, Vorträgen, u.v.m."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/LBV-HSG-Freising-Logo_460x284px.png"), ImageCaption: null.StringFrom("Bild: LBV"), LinkUrl: null.StringFrom("https://freising.lbv.de/hochschulgruppe/")},
		{Name: "LBV-HSG Straubing", Collection: "Interessensgruppen", Description: null.StringFrom("Bei uns kommen Naturbegeisterte zusammen, die sich für Natur- und Umweltschutz in Straubing und Umgebung engagieren."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/LBV_Straubing.png"), ImageCaption: null.StringFrom("Bild: LBV"), LinkUrl: null.StringFrom("https://straubing-bogen.lbv.de/über-uns/hochschulgruppe/")},
		{Name: "Öko-AK Landbau Weihenstephan", Collection: "Interessensgruppen"},
		{Name: "The Green Team", Collection: "Interessensgruppen", Description: null.StringFrom("Unser Ziel ist es, heutige Studierende bei der Ermöglichung einer nachhaltigeren Welt zu unterstützen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/The_Green_Team.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("http://thegreenteamtum.com/")},
		{Name: "TUM Renewable Energies Initiative", Collection: "Interessensgruppen"},
		{Name: "AcTUM", Collection: "Kunst &amp; Kultur", Description: null.StringFrom("AcTUM ist eine bunt gemischte Gruppe aus Mitgliedern der TUM, die einfach Lust auf Theater haben. Wir suchen immer begeisterte Leute!"), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/AcTUM.png"), ImageCaption: null.StringFrom("Bild: Tim Kuhr"), LinkUrl: null.StringFrom("https://ac-tum.de")},
		{Name: "Campus-Cneipe", Collection: "Kunst &amp; Kultur", Description: null.StringFrom("Die Campus-Cneipe sorgt seit 2005 für Freizeitspaß in Garching. Sie wird von TUM-Studenten geleitet, versorgt und am Laufen gehalten."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Campus-Cneipe.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://c2.tum.de")},
		{Name: "der tu film", Collection: "Kunst &amp; Kultur", Description: null.StringFrom("Hörsaal bei Tag - Kino bei Nacht: Dienstags und donnerstags zeigen wir die besten aktuellen Filme und Klassiker auf der großen Leinwand."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/der_tu_film.png"), ImageCaption: null.StringFrom("Bild: der tu film"), LinkUrl: null.StringFrom("https://www.tu-film.de/")},
		{Name: "Fusian Dance Crew", Collection: "Kunst &amp; Kultur", Description: null.StringFrom("Wir sind eine internationale Tanzcrew, die gemischten Hip-Hop trainiert, auftritt und dabei eine lustige Gemeinschaft aufbaut."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Fusian_Dance_Crew.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://msha.ke/fusiandance")},
		{Name: "Kulturklub", Collection: "Kunst &amp; Kultur", Description: null.StringFrom("Kultur genießen ist unser Motto. Wir gehen zusammen ins Theater, Konzert oder Museum."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/_processed_/9/a/csm_Kulturklub_badc149f97.jpg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://kulturklub.wixsite.com/kulturklub")},
		{Name: "NACHHOELZER", Collection: "Kunst &amp; Kultur", Description: null.StringFrom("Initiative für studentische Kultur am Stammgelände. Ein Freiraum und Austauschort, der überfakultäre und interdisziplinäre Interaktion fördert."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Nachhoelzer.png"), ImageCaption: null.StringFrom("Bild: NACHHOELZER"), LinkUrl: null.StringFrom("https://www.instagram.com/nachhoelzer/")},
		{Name: "Praias do Isar", Collection: "Kunst &amp; Kultur", Description: null.StringFrom("Praias do Isar ist eine Sambaschule im Rio-Stil, die sowohl Percussion- als auch Tanzunterricht anbietet."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Praias_do_Isar.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://praiasdoisar.de")},
		{Name: "Student DnD Club Heilbronn", Collection: "Kunst &amp; Kultur"},
		{Name: "TUdesign", Collection: "Kunst &amp; Kultur"},
		{Name: "TUM Buchclub", Collection: "Kunst &amp; Kultur", Description: null.StringFrom("Wir lesen nonfiction und fiction Bücher. Ziel dabei ist es, in den Diskussionen (3-4 Mal / Semester) unseren Horizont zu erweitern."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TUM_Buchclub.jpg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.instagram.com/tum.buchclub/")},
		{Name: "TUM Jazzband", Collection: "Kunst &amp; Kultur", Description: null.StringFrom("Wir sind eine Big Band von Studierenden für Studierende und heißen Musiker:innen von allen Münchner Hochschulen willkommen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TUM_Jazzband.jpg"), ImageCaption: null.StringFrom("Bild: Xaver Cox"), LinkUrl: null.StringFrom("https://www.tumjazzband.de")},
		{Name: "180DC München", Collection: "Karriere", Description: null.StringFrom("Als weltweit größte studentische Unternehmensberatung unterstützen wir gemeinnützige Organisationen und Social Start-Ups strategisch."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/180_Degrees_Consulting.png"), ImageCaption: null.StringFrom("Bild: 180 Degrees Consulting Munich e. V."), LinkUrl: null.StringFrom("https://180dcmunich.org/")},
		{Name: "Academy Consult", Collection: "Karriere", Description: null.StringFrom("Bei uns erweitern Studierende ihre theoretischen Kenntnisse aus dem Studium, indem sie Unternehmen auf externen Projekten beraten."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Academy_Consult.svg"), ImageCaption: null.StringFrom("Bild: Academy Consult München e. V."), LinkUrl: null.StringFrom("https://academyconsult.de/")},
		{Name: "EUROAVIA München", Collection: "Karriere", Description: null.StringFrom("EUROAVIA München fördert Luft- und Raumfahrtwissen, berufliche Entwicklung und internationales Networking bei Studierenden."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Euroavia.png"), ImageCaption: null.StringFrom("Bild: EUROAVIA München e.V."), LinkUrl: null.StringFrom("https://linktr.ee/euroavia.munich")},
		{Name: "Family Business Club", Collection: "Karriere", Description: null.StringFrom("Wir sind eine Gruppe an der TUM School of Management, die sich an Studier\u00adende mit Interesse an Familienunternehmen richtet."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Family_Business_Club.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.family-business-club.de")},
		{Name: "IAESTE an der TUM e. V.", Collection: "Karriere"},
		{Name: "IKOM", Collection: "Karriere", Description: null.StringFrom("Wir knüpfen Kontakte. Persönlich. Deutschlands größtes studentisches Karriereforum."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/IKOM.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.ikom-tum.de")},
		{Name: "Innovis VC", Collection: "Karriere", Description: null.StringFrom("Innovis VC ist Europas am schnellsten wachsende studentische Initiative an der Schnittstelle von Startups, Studenten und Investoren."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Innovis_VC.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.innovis.vc/")},
		{Name: "Nucleate Germany", Collection: "Karriere", Description: null.StringFrom("Nucleate ist eine globale Organi\u00adsation, die die Gründung von Firmen in der Biotechnologie unterstützt und künftige Führungskräfte fördert."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Nucleate_Germany.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://nucleate.xyz/")},
		{Name: "TUM Blockchain Club", Collection: "Karriere", Description: null.StringFrom("Wir stärken das Europäische Blockchain-Ökosystem unter Studierenden und bilden ein Forum für positive Nutzung von Blockchain."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TUM_Blockchain_Club.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.tum-blockchain.com/")},
		{Name: "TUM Club for Sustainable Entrepreneurship", Collection: "Karriere"},
		{Name: "VDE-IEEE Hochschulgruppe", Collection: "Karriere"},
		{Name: "Abasha e. V.", Collection: "Allgemeinwohl", Description: null.StringFrom("Abasha unterstützt junge Sport- und Bildungsinitiativen aus der ganzen Welt durch eine Plattform für digitales Ehrenamt."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Abasha.png"), ImageCaption: null.StringFrom("Bild: Abasha e. V."), LinkUrl: null.StringFrom("https://abasha.de/")},
		{Name: "Campus for Change e. V.", Collection: "Allgemeinwohl", Description: null.StringFrom("Wir realisieren weltweit soziale Projekte in den Bereichen Entwick\u00adlungshilfe, Bildung, Gesundheits\u00adversor\u00adgung und Flüchtlingshilfe."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Campus_for_Change.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.campusforchange.org")},
		{Name: "Ingenieure ohne Grenzen e. V.", Collection: "Allgemeinwohl", Description: null.StringFrom("Entwicklungszusammenarbeit: Mit technischem Wissen Lebens\u00adbeding\u00adungen verbessern, um das Zusam\u00admenwachsen der Welt zu fördern."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Ingenieure_Ohne_Grenzen.jpg"), ImageCaption: null.StringFrom("Bild: Ingenieure ohne Grenzen e. V."), LinkUrl: null.StringFrom("https://www.ingenieure-ohne-grenzen.org/de/mitmachen/regionalgruppe-muenchen")},
		{Name: "Munich Roots in Education", Collection: "Allgemeinwohl", Description: null.StringFrom("Wir möchten die Bildung von sozial benachteiligten Kindern und Jugend\u00adlichen durch Erlöse unserer Veran\u00adstalt\u00adungen fördern &amp; unterstützen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Munich_Roots_in_Education.png"), ImageCaption: null.StringFrom("Bild: Munich Roots in Education"), LinkUrl: null.StringFrom("https://munich-roots.de/")},
		{Name: "Nightline München", Collection: "Allgemeinwohl", Description: null.StringFrom("Wir sind ein Zuhör\u00adtelefon für Studier\u00adende. Am Telefon sitzen andere Studierende, die bei Problemen an der Uni und anderswo beistehen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Nightline_Muenchen.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.nightline-muc.de/de")},
		{Name: "ROCK YOUR LIFE! München", Collection: "Allgemeinwohl", Description: null.StringFrom("Unser Mentoring-Programm für Studierende und Mittelschüler*innen fördert Bildungsgerechtigkeit und Chancengleichheit."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Rock_Your_Life.png"), ImageCaption: null.StringFrom("Bild: ROCK YOUR LIFE! gGmbH"), LinkUrl: null.StringFrom("https://muenchen.rockyourlife.de/")},
		{Name: "Sailsetters", Collection: "Allgemeinwohl"},
		{Name: "Talente Spenden", Collection: "Allgemeinwohl", Description: null.StringFrom("Ehrenamtliches Engagement stärken: unsere Projekte sind etwa Mentoring für Geflüchtete sowie die Biotoppflege geschützter Wiesen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Talente_Spenden.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.tum.de/studium/studienfinanzierung/stipendien/stipendien-der-tum/deutschlandstipendium/was-wir-machen/wir-engagieren-uns/")},
		{Name: "TU eMpower Africa e. V.", Collection: "Allgemeinwohl", Description: null.StringFrom("TU eMpower Africa ist ein non-profit, das sich dem Ziel verschreibt, afrikanische Gemeinschaften zu stärken."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TU_Empower_Africa.png"), ImageCaption: null.StringFrom("Bild: TU eMpower Africa e. V."), LinkUrl: null.StringFrom("https://tu-empower-africa.org/")},
		{Name: "UNICEF-Hochschulgruppe München", Collection: "Allgemeinwohl"},
		{Name: "AEGEE München", Collection: "Internationales", Description: null.StringFrom("AEGEE ist eine der größten Jugend\u00adorganisationen Europas. Wir stehen für kulturelle Verständigung und Vereinigung aller jungen Europäer."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/AEGEE.jpg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://aegee-muenchen.de/")},
		{Name: "AIS-München", Collection: "Internationales", Description: null.StringFrom("Die AIS-München will Studierende unterstützen, ein ganzheitliches Bild der iranischen Kultur vermitteln und bei Integration &amp; Austausch helfen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/AIS_Munich.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.ais-munich.org/")},
		{Name: "LSM", Collection: "Internationales", Description: null.StringFrom("Der Verein ist eine luxemburgische Studentengruppe in München mit über 200 Aktiven, die sich austauschen und kennen lernen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Letzeburger_Studenten_zu_Mu__nchen.jpg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://lsm.lu/")},
		{Name: "Munich Egyptian Students Club", Collection: "Internationales"},
		{Name: "MUNTUM", Collection: "Internationales", Description: null.StringFrom("Als Gruppe für politisch Interessierte aller Studiengänge nehmen wir an internationalen Konferenzen Teil und erweitern Rhetorik und Horizonte."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/MUNTUM.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.muntum.org/")},
		{Name: "TUM Canadian Students Association", Collection: "Internationales"},
		{Name: "TUMHN Debate", Collection: "Internationales"},
		{Name: "AIESEC in München", Collection: "Mentoring"},
		{Name: "Lern-Fair Hub München", Collection: "Mentoring", Description: null.StringFrom("Wir setzen uns dafür ein, allen Schüler:innen in Deutschland dieselben Chancen auf Bildung zu ermöglichen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Lern-Fair.png"), ImageCaption: null.StringFrom("Bild: Lern-Fair e. V."), LinkUrl: null.StringFrom("https://www.lern-fair.de")},
		{Name: "she.codes by TEC", Collection: "Mentoring", Description: null.StringFrom("Um Mädchen für technische Themen zu begeistern und zu ermutigen, organisieren wir ehrenamtlich kostenlose Programmierworkshops."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/She_Codes.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://codes.education")},
		{Name: "Studenten bilden Schüler", Collection: "Mentoring", Description: null.StringFrom("Als ehrenamtlicher Student bei uns hilfst du sozial benachteiligten Schülern – per Nachhilfe in Fächern deiner Wahl oder in der Orga."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Studenten_bilden_Schueler-opt-min.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://studenten-bilden-schueler.de/standorte/muenchen")},
		{Name: "ESN TUMi e. V.", Collection: "Networking", Description: null.StringFrom("Wir organisieren für Austausch- und internationale Studenten Events: von kulturellen Veranstaltungen und Ausflügen bis zu Sport und Partys."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/ESN_TUMi.svg"), ImageCaption: null.StringFrom("Bild: ESN TUMi e. V."), LinkUrl: null.StringFrom("https://tumi.esn.world/events")},
		{Name: "Munich FinanceCircle", Collection: "Networking"},
		{Name: "Society of Sommeliers", Collection: "Networking", Description: null.StringFrom("Als Münchens erster studentischer Weinclub bilden wir Mitglieder in den Bereichen Wein, Weinverkostung und -herstellung weiter."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Society_of_Sommeliers.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.societyofsommeliers.de/")},
		{Name: "SRM Talks", Collection: "Networking", Description: null.StringFrom("Monatliche Vorträge mit Experten aus der Nachhaltigkeitsbranche, gefolgt von einem Networking-Essen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/SRM_Talks.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.linkedin.com/company/srmtalks")},
		{Name: "START Munich", Collection: "Networking", Description: null.StringFrom("Als größte Münchner studentische Initiative mit Fokus Entrepreneurship ermöglicht unser Netzwerk den Zugang zum Start-up-Ökosystem."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/START_Munich.png"), ImageCaption: null.StringFrom("Bild: START Munich"), LinkUrl: null.StringFrom("https://www.startmunich.de/")},
		{Name: "Stipendiennetzwerk München", Collection: "Networking", Description: null.StringFrom("Wir fördern den Austausch zwischen Stipendiat:innen der 13 Förderwerke in München und informieren Jugendliche über Stipendien."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Stipendiennetzwerk_Mu__nchen.svg"), ImageCaption: null.StringFrom("Bild: Stipendiennetzwerk München e.V."), LinkUrl: null.StringFrom("https://stipendiennetzwerk.de/")},
		{Name: "The Entrepreneurial Group", Collection: "Networking", Description: null.StringFrom("Die Entrepreneurial Group ist die studentische Initiative, die Studierenden den Start ins Unter\u00adneh\u00ad\u00admertum in München ermöglicht."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TEG.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.teg-ev.de")},
		{Name: "TUM Business Club e. V.", Collection: "Networking"},
		{Name: "Muslimische Studenteninitiative", Collection: "Religionsverbände"},
		{Name: "Electrochemical Society Student Chapter Munich", Collection: "Forschung und Anwendung"},
		{Name: "Future Foods", Collection: "Forschung und Anwendung"},
		{Name: "iGEM Munich", Collection: "Forschung und Anwendung", Description: null.StringFrom("Unser gemeinsames Team von TUM und LMU für den iGEM – dem größten internationalen Wettbewerb für synthetische Biologie."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/iGEM.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://igemmunich.mwn.de")},
		{Name: "neuroTUM", Collection: "Forschung und Anwendung", Description: null.StringFrom("Wir entwickeln ein Brain-Computer-Interface. Dazu analysieren wir eigene EEG-Tests mit modernsten Machine-Learning-Ansätzen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/NeuroTUM.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.neurotum.com/")},
		{Name: "TU Investment Club", Collection: "Forschung und Anwendung", Description: null.StringFrom("Der TU Investment Club ist ein Non-Profit, das sich der Förderung von Studenten mit einem Interesse an den Kapitalmärkten widmet."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TU_Investment_Club.svg"), ImageCaption: null.StringFrom("Bild: TU Investment Club e. V."), LinkUrl: null.StringFrom("https://tuinvest.de/")},
		{Name: "TUM Dev", Collection: "Forschung und Anwendung", Description: null.StringFrom("Wir entwickeln open source Websites und Apps für Studierende der TUM."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TUM_Dev.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://tum.dev")},
		{Name: "TUM.ai", Collection: "Forschung und Anwendung"},
		{Name: "Akaflieg München", Collection: "Forschung und Anwendung", Description: null.StringFrom("Wir sind ein studentischer Verein für alle Münchner Hochschulen. Seit 1924 konstruieren, bauen und fliegen wir Segel- &amp; Motorflugzeuge."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Akaflieg.svg"), ImageCaption: null.StringFrom("Bild: Akaflieg München"), LinkUrl: null.StringFrom("https://www.akaflieg-muenchen.de")},
		{Name: "Bioinformatics Munich Student Lab", Collection: "Forschung und Anwendung"},
		{Name: "DASH", Collection: "Forschung und Anwendung", Description: null.StringFrom("DASH entwickelt ein Exoskelett für die unteren Gliedmaßen von vollständig Querschnittsgelähmten, um neue Mobilität zu ermöglichen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/DASH.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://tumdash.com/")},
		{Name: "Elara Aerospace", Collection: "Forschung und Anwendung", Description: null.StringFrom("In den Weltraum &amp; darüber hinaus: Wir wollen als erstes Team eine von Studierenden gebaute Methalox-Rakete auf 100 km Höhe schießen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Elara_Aerospace.png"), ImageCaption: null.StringFrom("Bild: Elara Aerospace"), LinkUrl: null.StringFrom("https://elara-aerospace.com/")},
		{Name: "EnHands", Collection: "Forschung und Anwendung", Description: null.StringFrom("Unser Ziel als Hochschulgruppe ist die Entwicklung von zugänglichen, bezahlbaren und anpassbaren Prothesen für Entwicklungsländer."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/EnHands.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://enhands.de")},
		{Name: "Falcon Vision", Collection: "Forschung und Anwendung"},
		{Name: "Impetus Sailing Team", Collection: "Forschung und Anwendung", Description: null.StringFrom("Seit 2023 entwickeln, bauen und segeln wir Hochleistungssegelboote mit Fokus auf Strömungsmechanik, Leichtbau &amp; nachhaltige Bauweise."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Impetus_Sailing_Team.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.impetussailing.de/")},
		{Name: "LEVITUM", Collection: "Forschung und Anwendung", Description: null.StringFrom("Wir möchten mit Wasserstoffzellen die weltweit reichweitenstärkste eVTOL-Drohne mit Abfluggewicht unter 25 kg zu entwickeln."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/LEVITUM.jpg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.levitum.de")},
		{Name: "MedTech OneWorld Students", Collection: "Forschung und Anwendung", Description: null.StringFrom("Wir entwickeln mit lokalen Partnern medizintechnische Projekte für Entwicklungsländer. Alle Fachrichtungen sind willkommen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/MedTech_OneWorld_Students.png"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://medtechoneworldstudents.wordpress.com/")},
		{Name: "Phantum", Collection: "Forschung und Anwendung"},
		{Name: "RoboTUM", Collection: "Forschung und Anwendung"},
		{Name: "TUfast", Collection: "Forschung und Anwendung"},
		{Name: "TUM Boring", Collection: "Forschung und Anwendung", Description: null.StringFrom("Wir arbeiten an Technik für unterirdischen Verkehr. Unsere Tunnel\u00adbohrmaschine hat die Not-a-Boring Competition 2x gewonnen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TUM_Boring.svg"), ImageCaption: null.StringFrom("Bild: TUM Boring - Innovation in Tunneling"), LinkUrl: null.StringFrom("https://tum-boring.com/")},
		{Name: "TUM Carbon", Collection: "Forschung und Anwendung", Description: null.StringFrom("Wir möchten den Klimawandel lösen. Dazu entwicklen wir eine „Carbon Removal“-Technologie zur aktiven Entfernung atmosphärischen CO2s."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TUM_Carbon.jpg"), ImageCaption: null.StringFrom("Bild: Luisa Wunderlich"), LinkUrl: null.StringFrom("https://www.tumcarbon.com")},
		{Name: "TUM Phoenix Robotics", Collection: "Forschung und Anwendung", Description: null.StringFrom("Wir bearbeiten Projekte in den Bereichen des autonomen Fahrens und Fliegens."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/TUM_Phoenix_Robotics.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.epc.ed.tum.de/phoenix/tum-phoenix-robotics/")},
		{Name: "WARR", Collection: "Forschung und Anwendung", Description: null.StringFrom("Wir sind die WARR, die größte Studentengruppe der TUM. Wir arbeiten an ambitionierten Projekten rund um das Thema Weltraum."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/WARR.svg"), ImageCaption: null.StringFrom("Bild: WARR"), LinkUrl: null.StringFrom("https://www.warr.de")},
		{Name: "Munich eSports e. V.", Collection: "Sport", Description: null.StringFrom("Wir setzen uns an der TUM für den E-Sport ein: mit Mannschaften in der deutsche Uniliga, Gaming-Abenden, Viewing Parties, Turniere und mehr!"), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Munich_eSports.svg"), ImageCaption: null.StringFrom("Bild: Munich eSports e.V."), LinkUrl: null.StringFrom("https://munich-esports.de/")},
		{Name: "Munich Student Athletes Club", Collection: "Sport", Description: null.StringFrom("Wir fördern Austausch studentischer Leistungssportler*innen aller Sportarten und unterstützen sie in verschiedenen Bereichen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Munich_Student_Athlete_Club.jpg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.munich-sac.de/")},
		{Name: "Student AirRace", Collection: "Sport", Description: null.StringFrom("Das Ziel von Student AirRace ist es, Studierenden durch Wettflüge von Drohnen einen praktischen Zugang zur Luftfahrt zu ermöglichen."), ImageUrl: null.StringFrom("https://sv.tum.de/fileadmin/w00bdj/www/Hochschulgruppen/Logos/Student_AirRace.svg"), ImageCaption: null.StringFrom("Bild: Privat"), LinkUrl: null.StringFrom("https://www.student-airrace.com")},
		{Name: "TUM Chess Club", Collection: "Sport"},
		{Name: "TUM HN Hiking Club", Collection: "Sport"},
	}
	expectedCollections := []SVStudentClubCollection{
		{Name: "Akademisch", Description: "Hochschulgruppen, die sich direkt mit akademischen Inhalten beschäftigen – ob in einem Feld oder interdisziplinär"},
		{Name: "Interessensgruppen", Description: "Hochschulgruppen, die sich als Interessensvertretung für Teile der Studierendenschaft engagieren"},
		{Name: "Kunst &amp; Kultur", Description: "Gruppen mit einem Fokus auf Kunst- und Kulturveranstaltungen"},
		{Name: "Karriere", Description: "Gruppen, die Studierenden beim Einstieg in die Karriere oder bei einer studiumsbeglietenden Beschäftigung helfen"},
		{Name: "Allgemeinwohl", Description: "Gruppen, die sich vor Ort oder im Ausland für das Allgemeinwohl engagieren"},
		{Name: "Internationales", Description: "Gruppen, die den Austausch mit anderen Ländern fördern oder ein Netzwerk für internationale Studierende bieten"},
		{Name: "Mentoring", Description: "Gruppen, bei denen Studierende als Mentor*innen und/oder Mentees teilnehmen können"},
		{Name: "Networking", Description: "Gruppen, deren Hauptziel es ist, Studierende mit Gleichgesinnten zu vernetzen"},
		{Name: "Religionsverbände", Description: "Gruppen, die Austausch innerhalb oder zwischen Religionen fördern"},
		{Name: "Forschung und Anwendung", Description: "Gruppen, die ein Behandlung der Studiengangsinhalte durch eigene Forschung oder Anwendungsarbeit anstreben"},
		{Name: "Sport", Description: "Gruppen, die den gemeinsamen Sport fördern wollen"},
	}
	doc := strings.NewReader(htmlStudentClubs)
	clubs, collections, err := ParseStudentClubs(doc)
	require.NoError(t, err)
	require.Equal(t, expectedCollections, collections)
	require.Equal(t, expectedClubs, clubs)
}
