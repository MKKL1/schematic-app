create table post
(
    id             bigint not null
        constraint schematic_pk
            primary key,
    "desc"         text,
    owner          bigint not null,
    author_known   bigint,
    author_unknown text
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

