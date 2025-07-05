-- postgresql

CREATE TABLE players
(
    id   SERIAL PRIMARY KEY UNIQUE NOT NULL,
    name VARCHAR(40) UNIQUE NOT NULL
);

-- logically, there should be a single table, but do this to upgrade several table management skills
CREATE TABLE scores
(
    player_id INT REFERENCES players (id) UNIQUE NOT NULL,
    score     INT                                NOT NULL DEFAULT 0
);
