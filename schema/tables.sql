CREATE TABLE if NOT EXISTS project(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    gross_area NUMERIC NOT NULL DEFAULT 0,
    net_area NUMERIC NOT NULL DEFAULT 0,
    last_closure DATE DEFAULT NULL,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(name)
);
CREATE TABLE if NOT EXISTS supplier(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    supplier_id TEXT NOT NULL,
    contact_name TEXT,
    contact_email TEXT,
    contact_phone TEXT,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(supplier_id),
    UNIQUE(name)
);
CREATE TABLE if NOT EXISTS budget_item(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    level SMALLINT NOT NULL DEFAULT 1,
    accumulate BOOLEAN NOT NULL DEFAULT TRUE,
    parent_id UUID REFERENCES budget_item(id)
ON DELETE restrict,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(code),
    UNIQUE(name)
);
CREATE TABLE if NOT EXISTS category(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(name)
);
CREATE TABLE if NOT EXISTS material(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    unit TEXT NOT NULL,
    category_id UUID REFERENCES category(id)
ON DELETE restrict,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(name),
    UNIQUE(code)
);
CREATE TABLE if NOT EXISTS item(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    unit TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(name),
    UNIQUE(code)
);
CREATE TABLE if NOT EXISTS item_material(
    item_id UUID REFERENCES item(id)
ON DELETE restrict,
    material_id UUID REFERENCES material(id)
ON DELETE restrict,
    quantity NUMERIC NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY(item_id,
    material_id)
);
CREATE TABLE if NOT EXISTS budget(
    project_id UUID REFERENCES project(id)
ON DELETE restrict,
    budget_item_id UUID REFERENCES budget_item(id)
ON DELETE restrict,
    initial_quantity NUMERIC,
    initial_cost NUMERIC,
    initial_total NUMERIC NOT NULL,
    spent_quantity NUMERIC,
    spent_total NUMERIC NOT NULL,
    remaining_quantity NUMERIC,
    remaining_cost NUMERIC,
    remaining_total NUMERIC NOT NULL,
    updated_budget NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY(project_id,
    budget_item_id)
);
CREATE TABLE if NOT EXISTS invoice(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES project(id)
ON DELETE restrict,
    supplier_id UUID REFERENCES supplier(id)
ON DELETE restrict,
    invoice_number TEXT NOT NULL,
    invoice_date DATE NOT NULL,
    invoice_total NUMERIC NOT NULL,
    is_balanced BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(supplier_id,
    project_id,
    invoice_number)
);
ALTER TABLE invoice ALTER COLUMN invoice_total
SET DEFAULT 0;
CREATE TABLE if NOT EXISTS invoice_details(
    invoice_id UUID NOT NULL REFERENCES invoice(id)
ON DELETE restrict,
    budget_item_id UUID NOT NULL REFERENCES budget_item(id)
ON DELETE restrict,
    quantity NUMERIC NOT NULL,
    cost NUMERIC NOT NULL,
    total NUMERIC NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    UNIQUE(invoice_id,
    budget_item_id),
    PRIMARY KEY(invoice_id,
    budget_item_id)
);
CREATE TABLE if NOT EXISTS historic(
    project_id UUID NOT NULL REFERENCES project(id)
ON DELETE restrict,
    budget_item_id UUID NOT NULL REFERENCES budget_item(id)
ON DELETE restrict,
    historic_date DATE NOT NULL,
    initial_quantity NUMERIC,
    initial_cost NUMERIC,
    initial_total NUMERIC NOT NULL,
    spent_quantity NUMERIC,
    spent_total NUMERIC NOT NULL,
    remaining_quantity NUMERIC,
    remaining_cost NUMERIC,
    remaining_total NUMERIC NOT NULL,
    updated_budget NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(
        project_id,
        budget_item_id,
        historic_date
    ),
    PRIMARY KEY(
        project_id,
        budget_item_id,
        historic_date
    )
);
----------------------------------------
--                VIEWS               --
----------------------------------------
CREATE
OR REPLACE VIEW vw_budget_item AS SELECT
    b.id,
    b.code,
    b.name,
    b.level,
    b.accumulate,
    p.id AS parent_id,
    p.code AS parent_code,
    p.name AS parent_name
FROM
    budget_item b
LEFT JOIN budget_item p
    ON b.parent_id = p.id;
CREATE
OR REPLACE VIEW vw_materials AS SELECT
    m.id,
    m.code,
    m.name,
    m.unit,
    c.id AS category_id,
    c.name AS category_name
FROM
    material m
LEFT JOIN category c
    ON m.category_id = c.id;
CREATE
OR REPLACE VIEW vw_item_materials AS SELECT
    m.id AS material_id,
    m.code AS material_code,
    m.name AS material_name,
    m.unit AS material_unit,
    i.id AS item_id,
    i.code AS item_code,
    i.name AS item_name,
    i.unit AS item_unit,
    im.quantity
FROM
    item_material im
LEFT JOIN material m
    ON im.material_id = m.id
LEFT JOIN item i
    ON im.item_id = i.id;
CREATE
OR REPLACE VIEW vw_budget AS SELECT
    bi.id AS budget_item_id,
    bi.code AS budget_item_code,
    bi.name AS budget_item_name,
    bi.level AS budget_item_level,
    bi.accumulate AS budget_item_accumulate,
    p.id AS project_id,
    p.name AS project_name,
    b.initial_quantity,
    b.initial_cost,
    b.initial_total,
    b.spent_quantity,
    b.spent_total,
    b.remaining_quantity,
    b.remaining_cost,
    b.remaining_total,
    b.updated_budget
FROM
    budget b
LEFT JOIN project p
    ON b.project_id = p.id
LEFT JOIN budget_item bi
    ON b.budget_item_id = bi.id;
CREATE
OR REPLACE VIEW vw_invoice AS SELECT
    i.id,
    s.id AS supplier_id,
    s.supplier_id AS supplier_number,
    s.name AS supplier_name,
    s.contact_name,
    s.contact_email,
    s.contact_phone,
    p.id AS project_id,
    p.name AS project_name,
    p.is_active,
    i.invoice_number,
    i.invoice_date,
    i.invoice_total,
    i.is_balanced
FROM
    invoice i
JOIN supplier s
    ON i.supplier_id = s.id
JOIN project p
    ON i.project_id = p.id;
CREATE
OR REPLACE VIEW vw_invoice_details AS SELECT
    id.invoice_id,
    i.invoice_number,
    i.invoice_total AS invoice_total,
    i.invoice_date AS invoice_date,
    p.id AS project_id,
    p.name AS project_name,
    s.id AS supplier_id,
    s.supplier_id AS supplier_number,
    s.name AS supplier_name,
    id.budget_item_id,
    b.code AS budget_item_code,
    b.name AS budget_item_name,
    b.level AS budget_item_level,
    id.quantity,
    id.cost,
    id.total
FROM
    invoice_details id
JOIN budget_item b
    ON id.budget_item_id = b.id
JOIN invoice i
    ON id.invoice_id = i.id
JOIN supplier s
    ON i.supplier_id = s.id
JOIN project p
    ON i.project_id = p.id;
