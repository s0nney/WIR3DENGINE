-- Ideal table structure -- 

CREATE TABLE posts (
	id BIGSERIAL PRIMARY KEY,
	in_name VARCHAR(70),
	in_text VARCHAR(1024),
	ip VARCHAR(512),
	date_posted TIMESTAMP DEFAULT current_timestamp,
);

CREATE TABLE images (
	id BIGSERIAL PRIMARY KEY,
	filename VARCHAR(35),
	checksum VARCHAR(128),
	ip VARCHAR(512)
	date_posted TIMESTAMP DEFAULT current_timestamp,
);

CREATE TABLE bans (
	id BIGSERIAL PRIMARY KEY,
	ip VARCHAR(512),
	date_posted TIMESTAMP DEFAULT current_timestamp,
);
