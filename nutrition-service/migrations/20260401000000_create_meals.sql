-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS meals(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    calories DOUBLE PRECISION NOT NULL DEFAULT 0,
    proteins DOUBLE PRECISION NOT NULL DEFAULT 0,
    fats DOUBLE PRECISION NOT NULL DEFAULT 0,
    carbohydrates DOUBLE PRECISION NOT NULL DEFAULT 0,
    weight_grams INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS meals;
-- +goose StatementEnd