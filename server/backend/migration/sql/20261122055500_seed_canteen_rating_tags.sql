-- Seed canteen_rating_tag_options if empty (idempotent & retry-safe)
INSERT INTO `canteen_rating_tag_options` (`DE`, `EN`)
SELECT seed.`DE`, seed.`EN`
FROM (
  SELECT 'Ausgabe gut koordiniert' AS `DE`, 'Output Well Coordinated' AS `EN`
  UNION ALL SELECT 'Wartezeit', 'Waiting Time'
  UNION ALL SELECT 'Sauberkeit', 'Cleanliness'
  UNION ALL SELECT 'Verfügbarkeit der Gerichte', 'Availability'
  UNION ALL SELECT 'Genügend freie Plätze', 'Enough Free Tables'
  UNION ALL SELECT 'Auswahl Vegan', 'Variety Vegan'
  UNION ALL SELECT 'Auswahl Vegetarisch', 'Variety Vegetarian'
  UNION ALL SELECT 'Auswahl Allgemein', 'Variety General'
  UNION ALL SELECT 'Glaubhaftigkeit des Speiseplans', 'Credibility of the menu'
) AS seed
WHERE NOT EXISTS (
  SELECT 1 FROM (SELECT 1 FROM `canteen_rating_tag_options` LIMIT 1) AS existing
);
