CREATE TABLE projects(
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  daw TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  public BOOLEAN NOT NULL DEFAULT false,
  bpm INTEGER DEFAULT 0,
  key TEXT NOT NULL DEFAULT 'none',
  genre TEXT NOT NULL DEFAULT 'none',
  template BOOLEAN NOT NULL DEFAULT false
);