-- Fix typos and casing in seeded dish_name_tag_options labels
UPDATE `dish_name_tag_options` SET `EN` = 'Calamari' WHERE `EN` = 'Calamarie';
UPDATE `dish_name_tag_options` SET `EN` = 'Pumpkin' WHERE `EN` = 'Pumkin';
UPDATE `dish_name_tag_options` SET `EN` = 'Hummus' WHERE `EN` = 'Humus';
UPDATE `dish_name_tag_options` SET `DE` = 'Pilze' WHERE BINARY `DE` = 'pilze';
