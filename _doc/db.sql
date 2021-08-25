-- sylabes definition

CREATE TABLE sylabes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	sylabe VARCHAR(16) NOT NULL DEFAULT ''
, audio_file varchar(255) NOT NULL DEFAULT '');

CREATE UNIQUE INDEX unique_sylabe on sylabes(sylabe, audio_file);

-- words definition

CREATE TABLE words (
	id INTEGER PRIMARY KEY AUTOINCREMENT
, word VARCHAR(128) NOT NULL, audio_file varchar(255) NOT NULL DEFAULT '');

CREATE UNIQUE INDEX unique_word on words(word, audio_file);

-- words_sylabes definition

CREATE TABLE words_sylabes (
	word_id INTEGER NOT NULL,
	sylabe_id INTEGER NOT NULL,
	position INTEGER NOT NULL,
	FOREIGN KEY (word_id) REFERENCES words (id),
	FOREIGN KEY (sylabe_id) REFERENCES sylabes (id)
);

-- Check:

SELECT 
	w.word, group_concat(s.sylabe, ',') 
FROM words as w 
	JOIN words_sylabes ws ON ws.word_id = w.id 
	JOIN sylabes s ON ws.sylabe_id = s.id 
GROUP BY w.id 
ORDER BY w.id, ws.position;
