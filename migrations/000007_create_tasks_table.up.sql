CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- COULD BE ADDED MORE GROUPS

INSERT INTO groups (name) VALUES
    ('all'),
    ('user')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS tasks_statuses (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

INSERT INTO tasks_statuses (name) VALUES
    ('in progress'),
    ('completed'),
    ('cancelled')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status_id INTEGER REFERENCES tasks_statuses(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by INTEGER REFERENCES admins(id),
    amount DECIMAL(10, 2) NOT NULL,
    user_id INTEGER REFERENCES users(id),
    for_group_id INTEGER REFERENCES groups(id)
);

CREATE OR REPLACE FUNCTION check_user_group() RETURNS TRIGGER AS $$
DECLARE
    user_group_id INTEGER;
BEGIN
    SELECT id INTO user_group_id FROM groups WHERE name = 'user';

    IF (NEW.for_group_id = user_group_id) THEN
        IF NEW.user_id IS NULL THEN
            RAISE EXCEPTION 'user_id must be set for tasks in the "user" group';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_user_group
    BEFORE INSERT OR UPDATE ON tasks
    FOR EACH ROW EXECUTE FUNCTION check_user_group();