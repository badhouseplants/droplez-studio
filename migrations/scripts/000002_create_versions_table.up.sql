CREATE TABLE versions (
  id UUID PRIMARY KEY,
  version INTEGER NOT NULL,
  project_id UUID NOT NULL,
  object_name TEXT NOT NULL,
  message TEXT NOT NULL,
  uploaded_at TIMESTAMP NOT NULL,
  UNIQUE (project_id, version)
);