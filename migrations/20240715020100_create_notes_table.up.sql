CREATE TYPE note_target AS ENUM ('user', 'company');

--bun:split

CREATE TABLE notes (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    author_id         VARCHAR(255) NOT NULL,
    public_identifier VARCHAR(255) NOT NULL,
    target            note_target NOT NULL,

    content           TEXT NOT NULL,
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_notes_author_id UNIQUE (author_id, target, public_identifier)
);

--bun:split

CREATE INDEX notes_per_author ON notes (author_id);
