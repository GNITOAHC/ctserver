create table if not exists "user_t" (
    "mail" varchar(255) not null,
    "username" varchar(255) not null unique,
    "name" varchar(255),
    "phone" varchar(255),
    "created_at" date not null default now(),
    "plan" varchar(255) not null default 'free',
    primary key ("mail")
);

create table if not exists "data_t" (
    "username" varchar(255),
    "id" uuid not null default uuid_generate_v4(),
    "path" varchar(255) not null, -- Path of the data
    "type" varchar(4) not null, -- "file", "url", "dir" or "text"
    "content" text not null, -- URL to file or redirection, directory's name or text content
    "description" varchar(255), -- Short description of the data
    "ancestor_id" uuid, -- parent id null for root
    "descendant_id" uuid[], -- child id, only for dir
    "created_at" timestamp with time zone not null default now(),
    "expired_at" timestamp with time zone,
    primary key ("id")
);

create table if not exists "shortened_t" (
    "data_id" uuid not null primary key,
    "shortened" text not null unique
);

alter table "shortened_t" add foreign key ("data_id") references "data_t" ("id") on delete cascade;
