ALTER TABLE ink_return_detail
    ADD ink_export_detail_id VARCHAR(50) DEFAULT '':::STRING;
ALTER TABLE ink_return
    ADD ink_export_id VARCHAR(50) DEFAULT '':::STRING;
-- create index for ink_export_detail_id
CREATE INDEX ink_return_detail_ink_export_detail_id_idx ON ink_return_detail (ink_export_detail_id);
--- create index for ink_return_id of ink_return_detail
CREATE INDEX ink_return_detail_ink_return_id_idx ON ink_return_detail (ink_return_id);
-- create index for ink_export_id
CREATE INDEX ink_return_ink_export_id_idx ON ink_return (ink_export_id);
