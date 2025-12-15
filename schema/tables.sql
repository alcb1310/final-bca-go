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
