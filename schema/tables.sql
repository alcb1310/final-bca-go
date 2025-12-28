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

create table if not exists category(
    id uuid primary key default gen_random_uuid(),
    name text not null,
    created_at timestamp default now(),

    unique (name)
);

create table if not exists material(
    id uuid primary key default gen_random_uuid(),
    code text not null,
    name text not null,
    unit text not null,

    category_id uuid references category(id) on delete restrict,
    created_at timestamp default now(),

    unique (name),
    unique (code)
);

create table if not exists item(
    id uuid primary key default gen_random_uuid(),
    code text not null,
    name text not null,
    unit text not null,

    category_id uuid references category(id) on delete restrict,
    created_at timestamp default now(),

    unique (name),
    unique (code)
);

----------------------------------------
--                VIEWS               --
----------------------------------------

create or replace view vw_budget_item as
select
    b.id,
    b.code,
    b.name,
    b.level,
    b.accumulate,
    p.id as parent_id,
    p.code as parent_code,
    p.name as parent_name
from budget_item b
left join budget_item p on b.parent_id = p.id;

create or replace view vw_materials as
select
    m.id,
    m.code,
    m.name,
    m.unit,
    c.id as category_id,
    c.name as category_name
from material m
left join category c on m.category_id = c.id;
