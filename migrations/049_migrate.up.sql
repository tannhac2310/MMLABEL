-- alter inspection_forms add columns
ALTER TABLE inspection_forms
ADD COLUMN IF NOT EXISTS stage_ids TEXT,
ADD COLUMN IF NOT EXISTS qa_note TEXT,
ADD COLUMN IF NOT EXISTS data JSONB;