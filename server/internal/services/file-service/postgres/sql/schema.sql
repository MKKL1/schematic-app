create table tmp_file (
    store_key text primary key,
    file_name text not null,
    content_type text not null,
    status text not null default 'pending',
    error_reason text,
    processing_attempts int default 0,
    final_hash text,
    expires_at timestamptz not null,
    created_at timestamptz not null default NOW(),
    updated_at timestamptz not null default NOW()
);

CREATE INDEX idx_temporary_files_expires_at ON tmp_file(expires_at);

create table file (
    hash text primary key,
    original_temp_id text not null references tmp_file(store_key),
    file_size int not null,
    content_type text not null,
    created_at timestamptz not null default NOW(),
    updated_at timestamptz not null default NOW()
);