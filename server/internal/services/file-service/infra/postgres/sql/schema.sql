create table tmp_file (
    file_hash text primary key,
    store_key text unique not null,
    file_name text not null,
    content_type text not null,
    file_size bigint not null,
    expires_at timestamptz not null,
    created_at timestamptz not null default NOW(),
    updated_at timestamptz not null default NOW()
);

CREATE INDEX idx_temporary_files_expires_at ON tmp_file(expires_at);