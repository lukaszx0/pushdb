-- Create table in which keys are stored
CREATE TABLE keys (
    id SERIAL,
    name TEXT UNIQUE,
    value BYTEA,
    version INT DEFAULT 1
);

-- Function sending notification on key_change channel
CREATE OR REPLACE FUNCTION key_change_event() RETURNS TRIGGER AS $$
DECLARE
        payload json;
    BEGIN
        IF (TG_OP = 'DELETE') THEN
            payload = json_build_object(
                'action', TG_OP);
        ELSE
            payload = json_build_object(
                'action',TG_OP,
                'key', row_to_json(NEW));
        END IF;
        PERFORM pg_notify('key_change', payload::text);
        -- Result is ignored since this is an AFTER trigger
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

-- Trigger on keys table to execute key_change_event() on every table change
CREATE TRIGGER keys_change_event
AFTER INSERT OR UPDATE OR DELETE ON keys
FOR EACH ROW EXECUTE PROCEDURE key_change_event();
