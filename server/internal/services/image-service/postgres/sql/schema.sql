create table image (
    file_hash text not null,
    image_type text not null, --Gallery/avatar...
    created_at timestamptz not null default NOW(),
    primary key(file_hash, image_type)
);