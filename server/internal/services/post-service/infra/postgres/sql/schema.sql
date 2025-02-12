create table post
(
    id             bigint not null
        constraint schematic_pk
            primary key,
    name           text not null,
    "desc"         text,
    owner          bigint not null,
    author_id   BIGINT
);

create table gallery_image
(
    image_id bigint not null
        constraint gallery_image_pk
            primary key,
    post_id  bigint
        constraint gallery_image_post_id_fk
            references post,
    "desc"   text
);

create table categories (
    name text primary key,
    metadata_schema jsonb not null
);

CREATE TABLE post_category_metadata (
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