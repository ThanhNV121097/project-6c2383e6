-- Greeting singleton schema.
CREATE TABLE greetings (
    id boolean PRIMARY KEY CHECK (id),
    text text NOT NULL CHECK (btrim(text) <> '')
);

INSERT INTO greetings (id, text) VALUES (true, 'Hello, World!')
ON CONFLICT (id) DO NOTHING;
