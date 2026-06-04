-- SMDR MVP Database Schema

CREATE EXTENSION IF NOT EXISTS "citext";

-- Users table for Owner and Sales
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    username        CITEXT  NOT NULL UNIQUE,
    hashed_password TEXT    NOT NULL,
    name            TEXT    NOT NULL,
    phone           TEXT    NOT NULL DEFAULT '',
    role            TEXT    NOT NULL CHECK (role IN ('owner', 'sales')),
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    must_change_password BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Products table
CREATE TABLE products (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT    NOT NULL,
    sku        TEXT    NOT NULL UNIQUE,
    price      BIGINT  NOT NULL CHECK (price > 0),
    active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_by BIGINT  NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Unique active product names (partial unique index)
CREATE UNIQUE INDEX idx_products_active_name ON products (LOWER(name)) WHERE active = TRUE;

-- Warungs table
CREATE TABLE warungs (
    id             BIGSERIAL PRIMARY KEY,
    name           TEXT    NOT NULL,
    owner_name     TEXT    NOT NULL DEFAULT '',
    phone          TEXT    NOT NULL DEFAULT '',
    address        TEXT    NOT NULL DEFAULT '',
    latitude       DOUBLE PRECISION,
    longitude      DOUBLE PRECISION,
    active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_by     BIGINT  NOT NULL REFERENCES users(id),
    acquired_from_checkin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Visits (Kunjungan) table
CREATE TABLE visits (
    id               BIGSERIAL PRIMARY KEY,
    sales_id         BIGINT    NOT NULL REFERENCES users(id),
    warung_id        BIGINT    NOT NULL REFERENCES warungs(id),
    status           TEXT      NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'selesai', 'dibatalkan')),
    check_in_lat     DOUBLE PRECISION NOT NULL,
    check_in_lng     DOUBLE PRECISION NOT NULL,
    check_in_accuracy DOUBLE PRECISION,
    check_in_time    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    distance_meters  DOUBLE PRECISION,
    business_date    DATE      NOT NULL DEFAULT (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta'),
    notes            TEXT      NOT NULL DEFAULT '',
    cancelled_at     TIMESTAMPTZ,
    cancelled_by     BIGINT    REFERENCES users(id),
    cancel_reason    TEXT      NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- One draft per sales
CREATE UNIQUE INDEX idx_visits_one_draft_per_sales ON visits (sales_id) WHERE status = 'draft';

-- Visit deposit items (Titipan)
CREATE TABLE visit_deposit_items (
    id         BIGSERIAL PRIMARY KEY,
    visit_id   BIGINT NOT NULL REFERENCES visits(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity   INT    NOT NULL CHECK (quantity > 0),
    UNIQUE (visit_id, product_id)
);

-- Visit sale items (Penjualan)
CREATE TABLE visit_sale_items (
    id         BIGSERIAL PRIMARY KEY,
    visit_id   BIGINT NOT NULL REFERENCES visits(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity   INT    NOT NULL CHECK (quantity > 0),
    price      BIGINT NOT NULL CHECK (price > 0),
    UNIQUE (visit_id, product_id)
);

-- Visit return items (Retur)
CREATE TABLE visit_return_items (
    id         BIGSERIAL PRIMARY KEY,
    visit_id   BIGINT NOT NULL REFERENCES visits(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity   INT    NOT NULL CHECK (quantity > 0),
    reason     TEXT   NOT NULL DEFAULT '',
    UNIQUE (visit_id, product_id)
);

-- Visit payments (Pembayaran) — at most one per visit
CREATE TABLE visit_payments (
    id       BIGSERIAL PRIMARY KEY,
    visit_id BIGINT NOT NULL UNIQUE REFERENCES visits(id),
    amount   BIGINT NOT NULL CHECK (amount >= 0),
    method   TEXT   NOT NULL CHECK (method IN ('tunai', 'transfer', 'lainnya')),
    note     TEXT   NOT NULL DEFAULT ''
);

-- Audit logs
CREATE TABLE audit_logs (
    id         BIGSERIAL PRIMARY KEY,
    actor_id   BIGINT NOT NULL REFERENCES users(id),
    action     TEXT   NOT NULL,
    entity     TEXT   NOT NULL,
    entity_id  BIGINT NOT NULL,
    summary    TEXT   NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_visits_sales_id ON visits (sales_id);
CREATE INDEX idx_visits_warung_id ON visits (warung_id);
CREATE INDEX idx_visits_business_date ON visits (business_date);
CREATE INDEX idx_visits_status ON visits (status);
CREATE INDEX idx_visit_deposit_items_visit ON visit_deposit_items (visit_id);
CREATE INDEX idx_visit_sale_items_visit ON visit_sale_items (visit_id);
CREATE INDEX idx_visit_return_items_visit ON visit_return_items (visit_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity, entity_id);
CREATE INDEX idx_audit_logs_actor_id ON audit_logs (actor_id);
CREATE INDEX idx_warungs_name ON warungs (LOWER(name));
CREATE INDEX idx_warungs_active ON warungs (active) WHERE active = TRUE;

-- Updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_warungs_updated_at BEFORE UPDATE ON warungs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_visits_updated_at BEFORE UPDATE ON visits FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
