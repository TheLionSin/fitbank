-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS activities (
                                        id UUID PRIMARY KEY,
                                        type TEXT NOT NULL,
                                        weight DOUBLE PRECISION DEFAULT 0,
                                        reps INTEGER DEFAULT 0,
                                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS activities;
-- +goose StatementEnd