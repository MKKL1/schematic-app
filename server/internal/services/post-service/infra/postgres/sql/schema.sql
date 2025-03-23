create table post
(
    id             bigint not null
        constraint schematic_pk
            primary key,
    name           text not null,
    "desc"         text,
    owner          bigint not null,
    author_id   bigint,
    created_at timestamptz not null default NOW(),
    updated_at timestamptz not null default NOW()
);

create table gallery_image
(
    post_id  bigint
        constraint gallery_image_post_id_fk
            references post,
    file_id text not null,
    "order" smallint not null,
    "desc"   text
);

create table attached_files
(
    hash text,
    temp_id uuid,
    post_id bigint not null references post,
    name text not null default '',
    file_size int not null default 0,
    downloads int not null default 0,
    created_at timestamptz not null default NOW(),
    updated_at timestamptz not null default NOW()
);

create table categories (
    name text primary key,
    metadata_schema jsonb not null
);

create table post_category_metadata (
    post_id bigint not null references post(id) on delete cascade,
    category text not null references categories(name) on delete cascade,
    metadata jsonb not null,
    primary key (post_id, category)
);

create index idx_all_values_search on post_category_metadata
    using GIN (metadata jsonb_path_ops);

create index idx_category_lookup on post_category_metadata (category);

create table post_tags (
    post_id bigint not null ,
    tag text not null ,
    primary key (post_id, tag)
);

create index idx_post_tags_tag on post_tags(tag);

CREATE TYPE category_metadata_pair AS (
    category text,
    metadata jsonb
);