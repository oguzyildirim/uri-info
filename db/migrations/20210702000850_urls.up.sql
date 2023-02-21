CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE urls (
  id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  HTML_version    VARCHAR NOT NULL,
  page_title    VARCHAR NOT NULL,
  headings_count    VARCHAR NOT NULL,
  links_count    INTEGER NOT NULL,
  inaccessible_links_count    INTEGER NOT NULL,
  have_login_form    BOOLEAN NOT NULL
);