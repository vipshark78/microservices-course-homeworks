-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders
(
    order_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
    user_uuid uuid NOT NULL,
    part_uuids uuid[] NOT NULL,
    total_price double precision NOT NULL,
    transaction_uuid uuid NULL,
    payment_method text NULL,
    status text NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    PRIMARY KEY (order_uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
