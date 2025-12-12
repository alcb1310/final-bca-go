create table if not exists project (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    is_active boolean not null,
    gross_area numeric not null default 0,
    net_area numeric not null default 0,
    last_closure date default null,
    created_at timestamp default now()
);
