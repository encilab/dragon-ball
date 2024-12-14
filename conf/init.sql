
CREATE TABLE IF NOT EXISTS character_dragonball (
	id INT CHECK (value >= 0) NOT NULL,
	name VARCHAR(64) NOT NULL,
	key VARCHAR(256) NOT NULL,
	race VARCHAR(64) NOT NULL,
    url_image VARCHAR(256) NOT NULL,
	PRIMARY KEY (id)
);

CREATE INDEX idx_character_name ON character_dragonball (name);
