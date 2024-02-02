create database if not exists images;
create table app_file
(
    id         char(36)                              not null
        primary key,
    name       varchar(255)                          not null,
    mime_type  varchar(50)                           not null,
    size       bigint                                not null,
    path       text                                  not null,
    created_at timestamp default current_timestamp() null,
    updated_at timestamp default current_timestamp() null on update current_timestamp()
);

