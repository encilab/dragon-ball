
CREATE TABLE IF NOT EXISTS character_dragonball (
	id INT NOT NULL,
	name VARCHAR(64) NOT NULL,
	ki VARCHAR(256) NOT NULL,
	race VARCHAR(64) NOT NULL,
    image VARCHAR(256) NOT NULL,
	PRIMARY KEY (id)
);

CREATE INDEX idx_character_name ON character_dragonball (name);
