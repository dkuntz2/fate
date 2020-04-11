BEGIN;

CREATE TABLE character (
  pk bigserial PRIMARY KEY,

  player text NOT NULL UNIQUE,
  name   text NOT NULL,

  refresh     int8 NOT NULL DEFAULT 1,
  fate_points int8 NOT NULL DEFAULT 0,

  high_concept text   NOT NULL,
  trouble      text   NOT NULL,
  aspects      text[] NOT NULL DEFAULT '{}',

  stress       int8   NOT NULL DEFAULT 0,
  consequences text[] NOT NULL DEFAULT '{}',

  careful  int8 NOT NULL DEFAULT 0,
  clever   int8 NOT NULL DEFAULT 0,
  flashy   int8 NOT NULL DEFAULT 0,
  forceful int8 NOT NULL DEFAULT 0,
  quick    int8 NOT NULL DEFAULT 0,
  sneaky   int8 NOT NULL DEFAULT 0
);

COMMIT;
