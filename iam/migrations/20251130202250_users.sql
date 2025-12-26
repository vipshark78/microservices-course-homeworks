-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    user_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
    login text NOT NULL unique,
    email text NOT NULL unique,
    password_hash text NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    PRIMARY KEY (user_uuid)
);

CREATE INDEX IF NOT EXISTS idx_user_uuid ON users (user_uuid);

create table notification_methods (
    id serial primary key,
    user_uuid UUID not null references users(user_uuid) on delete cascade ,
    provider_name text not null ,
    target text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
