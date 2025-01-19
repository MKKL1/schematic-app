create table users
(
    id   bigint not null
        constraint users_pk
            primary key,
    name text   not null
        constraint users_name_unique
            unique,
    oidc_sub uuid not null
        constraint users_oidc_sub_unique
            unique
);
