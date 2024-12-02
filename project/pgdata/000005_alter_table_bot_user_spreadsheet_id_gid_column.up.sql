ALTER TABLE bot_user
ADD COLUMN spreadsheet_id          VARCHAR(255) DEFAULT '',
ADD COLUMN spreadsheet_gid         VARCHAR(255) DEFAULT '',
ADD COLUMN spreadsheet_base_gid    VARCHAR(255) DEFAULT '',
ADD COLUMN spreadsheet_column      VARCHAR(255) DEFAULT '';