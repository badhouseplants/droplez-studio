CREATE TABLE projects(
  id UUID PRIMARY KEY,
  name TEXT,
  description TEXT,
  public BOOLEAN DEFAULT false,
  bpm INTEGER,
  key TEXT,
  genre TEXT,
);