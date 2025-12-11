-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    user_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
    login text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    notification_methods jsonb NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    PRIMARY KEY (user_uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
