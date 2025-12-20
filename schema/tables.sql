create table if not exists project (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    is_active boolean not null,
    gross_area numeric not null default 0,
    net_area numeric not null default 0,
    last_closure date default null,
    created_at timestamp default now(),

    unique (name)
);

create table if not exists supplier (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    supplier_id text not null,
    contact_name text,
    contact_email text,
    contact_phone text,
    created_at timestamp default now(),

    unique (supplier_id),
    unique (name)
);

create table if not exists budget_item(
    id uuid primary key default gen_random_uuid(),
    code text not null,
    name text not null,
    level smallint not null default 1,
    accumulate boolean not null default true,
    parent_id uuid references budget_item(id) on delete restrict,
    created_at timestamp default now(),

    unique (code),
    unique (name)
);
