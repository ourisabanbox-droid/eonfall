ALTER TABLE civilizations
    ADD CONSTRAINT fk_civilizations_capital_region
        FOREIGN KEY (capital_region_id) REFERENCES regions(id) ON DELETE SET NULL;