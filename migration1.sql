CREATE TABLE band (
	id SERIAL PRIMARY KEY,
	name TEXT
);

CREATE TABLE album (
	id SERIAL PRIMARY KEY,
	band_id integer REFERENCES band (id),
	name TEXT
);

CREATE TABLE track (
	id SERIAL PRIMARY KEY,
	name TEXT,
	album_id integer REFERENCES album (id)
);
