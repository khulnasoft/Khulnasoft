\set pguser `echo "$POSTGRES_USER"`

\c _khulnasoft
create schema if not exists _analytics;
alter schema _analytics owner to :pguser;
\c postgres
