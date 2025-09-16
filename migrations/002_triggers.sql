CREATE OR REPLACE FUNCTION log_item_changes() RETURNS TRIGGER AS $$
DECLARE
  uname TEXT;
BEGIN
  uname := current_setting('app.user', true);
  IF TG_OP = 'INSERT' THEN
    INSERT INTO history (item_id, action, changed_by, old_data, new_data)
    VALUES (NEW.id, 'insert', uname, NULL, to_jsonb(NEW));
    RETURN NEW;
  ELSIF TG_OP = 'UPDATE' THEN
    INSERT INTO history (item_id, action, changed_by, old_data, new_data)
    VALUES (NEW.id, 'update', uname, to_jsonb(OLD), to_jsonb(NEW));
    RETURN NEW;
  ELSIF TG_OP = 'DELETE' THEN
    INSERT INTO history (item_id, action, changed_by, old_data, new_data)
    VALUES (OLD.id, 'delete', uname, to_jsonb(OLD), NULL);
    RETURN OLD;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_history_item'
      AND table_name = 'history'
  ) THEN
    ALTER TABLE history DROP CONSTRAINT fk_history_item;
  END IF;
END$$;

ALTER TABLE history
  ADD CONSTRAINT fk_history_item
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE SET NULL;


DROP TRIGGER IF EXISTS trg_items_hist_ins ON items;
DROP TRIGGER IF EXISTS trg_items_hist_upd ON items;
DROP TRIGGER IF EXISTS trg_items_hist_del ON items;

CREATE TRIGGER trg_items_hist_ins AFTER INSERT ON items
FOR EACH ROW EXECUTE FUNCTION log_item_changes();

CREATE TRIGGER trg_items_hist_upd AFTER UPDATE ON items
FOR EACH ROW EXECUTE FUNCTION log_item_changes();

CREATE TRIGGER trg_items_hist_del BEFORE DELETE ON items
FOR EACH ROW EXECUTE FUNCTION log_item_changes();
