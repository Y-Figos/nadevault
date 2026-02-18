BEGIN;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
-- ---------------------------
-- ENUMS
-- ---------------------------
DO $$ BEGIN
  CREATE TYPE nade_type AS ENUM ('smoke','moly','frag','flashbang');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
  CREATE TYPE side_type AS ENUM ('T','CT');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
  CREATE TYPE mouse_click AS ENUM ('mouse1','mouse2','both');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- ---------------------------
-- MAPS
-- ---------------------------
CREATE TABLE IF NOT EXISTS cs_maps (
  id           SMALLSERIAL PRIMARY KEY,
  code         TEXT NOT NULL UNIQUE,     -- 'mirage', 'inferno'
  display_name TEXT NOT NULL,            -- 'Mirage'
  is_active    BOOLEAN NOT NULL DEFAULT TRUE,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------
-- IMMUTABLE SEARCH FUNCTION
-- ---------------------------
CREATE OR REPLACE FUNCTION nades_tsvector(
    name text,
    description text,
    from_callout text,
    to_callout text
) RETURNS tsvector AS $$
BEGIN
    RETURN to_tsvector('simple',
        coalesce(name, '') || ' ' ||
        coalesce(description, '') || ' ' ||
        coalesce(from_callout, '') || ' ' ||
        coalesce(to_callout, '')
    );
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- ---------------------------
-- NADES (lineups)
-- ---------------------------
CREATE TABLE IF NOT EXISTS nades (
  id            BIGSERIAL PRIMARY KEY,

  -- basic info
  name          TEXT NOT NULL,
  description   TEXT NOT NULL DEFAULT '',
  public_id     UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
  map_id        SMALLINT NOT NULL REFERENCES cs_maps(id) ON DELETE RESTRICT,
  nade_type     nade_type NOT NULL,
  common_side   side_type NOT NULL,

  -- key fields for fast lookup
  from_callout  TEXT NOT NULL,
  to_callout    TEXT NOT NULL,

  -- input modifiers
  mouse_click   mouse_click NOT NULL DEFAULT 'mouse1',
  is_jumping    BOOLEAN NOT NULL DEFAULT FALSE,
  is_running    BOOLEAN NOT NULL DEFAULT FALSE,
  is_walking    BOOLEAN NOT NULL DEFAULT FALSE,

  -- media payload
  images        JSONB NOT NULL DEFAULT '{}'::jsonb,

  -- community / governance
  is_public     BOOLEAN NOT NULL DEFAULT TRUE,
  created_by    TEXT,

  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  
  -- CORRECT DEFINITION (No trailing comma after STORED)
  search_tsv tsvector GENERATED ALWAYS AS (
        nades_tsvector(name, description, from_callout, to_callout)
    ) STORED
);

-- ---------------------------
-- UPDATED_AT trigger
-- ---------------------------
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_nades_updated_at ON nades;
CREATE TRIGGER trg_nades_updated_at
BEFORE UPDATE ON nades
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- ---------------------------
-- INDEXES
-- ---------------------------
-- Core filters
CREATE INDEX IF NOT EXISTS idx_nades_map_type_side
  ON nades (map_id, nade_type, common_side);

CREATE INDEX IF NOT EXISTS idx_nades_map_from
  ON nades (map_id, from_callout);

CREATE INDEX IF NOT EXISTS idx_nades_map_to
  ON nades (map_id, to_callout);

-- Public browsing
CREATE INDEX IF NOT EXISTS idx_nades_public_map
  ON nades (is_public, map_id);

-- Full-text search
CREATE INDEX IF NOT EXISTS idx_nades_search_tsv
  ON nades USING GIN (search_tsv);

-- JSONB index
CREATE INDEX IF NOT EXISTS idx_nades_images_gin
  ON nades USING GIN (images);

-- ---------------------------
-- SEED MAPS
-- ---------------------------
INSERT INTO cs_maps (code, display_name)
VALUES
  ('mirage','Mirage'),
  ('inferno','Inferno'),
  ('dust2','Dust II'),
  ('anubis','Anubis'),
  ('ancient','Ancient'),
  ('overpass','Overpass'),
  ('nuke','Nuke'),
  ('cache','Cache'),
  ('train','Train'),
  ('cobblestone','Cobblestone')
ON CONFLICT (code) DO NOTHING;

COMMIT;