create table if not exists "user_t" (
    "mail" varchar(255) not null,
    "phone" varchar(255),
    "created_at" date not null default now(),
    "plan" varchar(255) not null default 'free',
    primary key ("mail")
);

create table if not exists "data_t" (
    "user_mail" varchar(255),
    "id" uuid not null default uuid_generate_v4(),
    "name" varchar(255) not null, -- Shortened name of the data
    "type" varchar(4) not null, -- "file", "url", "dir" or "text"
    "content" text not null, -- URL to file or redirection, or text content of dir or text
    "ancestor_id" uuid, -- parent id null for root
    "descendant_id" uuid[], -- child id, only for dir
    "created_at" timestamp with time zone not null default now(),
    "expired_at" timestamp,
    primary key ("id")
);
