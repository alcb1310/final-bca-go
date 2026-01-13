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

    created_at timestamp default now(),

    unique (name),
    unique (code)
);

create table if not exists item_material(
    item_id uuid references item(id) on delete restrict,
    material_id uuid references material(id) on delete restrict,
    quantity numeric not null default 0,
    created_at timestamp default now(),

    primary key (item_id, material_id)
);

create table if not exists budget(
    project_id uuid references project(id) on delete restrict,
    budget_item_id uuid references budget_item(id) on delete restrict,

    initial_quantity numeric,
    initial_cost numeric,
    initial_total numeric not null,

    spent_quantity numeric,
    spent_total numeric not null,

    remaining_quantity numeric,
    remaining_cost numeric,
    remaining_total numeric not null,

    updated_budget numeric not null,

    created_at timestamp default now(),

    primary key (project_id, budget_item_id)
);

create table if not exists invoice(
    id uuid primary key default gen_random_uuid(),
    project_id uuid references project(id) on delete restrict,
    supplier_id uuid references supplier(id) on delete restrict,
    invoice_number text not null,
    invoice_date date not null,
    invoice_total numeric not null,
    is_balanced boolean not null default false,
    created_at timestamp default now(),

    unique (supplier_id, project_id, invoice_number)
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

create or replace view vw_item_materials as
select
    m.id as material_id,
    m.code as material_code,
    m.name as material_name,
    m.unit as material_unit,
    i.id as item_id,
    i.code as item_code,
    i.name as item_name,
    i.unit as item_unit,
    im.quantity
from item_material im
left join material m on im.material_id = m.id
left join item i on im.item_id = i.id;

create or replace view vw_budget as
select
    bi.id as budget_item_id,
    bi.code as budget_item_code,
    bi.name as budget_item_name,
    bi.level as budget_item_level,
    bi.accumulate as budget_item_accumulate,
    p.id as project_id,
    p.name as project_name,
    b.initial_quantity,
    b.initial_cost,
    b.initial_total,
    b.spent_quantity,
    b.spent_total,
    b.remaining_quantity,
    b.remaining_cost,
    b.remaining_total,
    b.updated_budget
from budget b
left join project p on b.project_id = p.id
left join budget_item bi on b.budget_item_id = bi.id;
