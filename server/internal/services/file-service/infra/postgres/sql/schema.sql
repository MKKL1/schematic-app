create table tmp_file (
    store_key text primary key,
    file_name text not null,
    expires_at timestamptz not null,
    created_at timestamptz not null default NOW(),
    updated_at timestamptz not null default NOW()
);

CREATE INDEX idx_temporary_files_expires_at ON tmp_file(expires_at);