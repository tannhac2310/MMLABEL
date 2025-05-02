
DELETE FROM options WHERE entity = 'QC';
INSERT INTO options (id, entity, code, name, created_by, created_at, updated_at)
VALUES
    ('L16', 'QC', 'L16', 'Khác', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L15', 'QC', 'L15', 'Đổ keo không đạt', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L14', 'QC', 'L14', 'Dập không đạt', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L13', 'QC', 'L13', 'Bế không đạt', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L12', 'QC', 'L12', 'Cắt Laser, vi tính không đạt', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L11', 'QC', 'L11', 'Cắt không đạt', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L10', 'QC', 'L10', 'Cán không đạt', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L09', 'QC', 'L09', 'Keo nhăn', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L08', 'QC', 'L08', 'Lỗi bong tróc mực', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L07', 'QC', 'L07', 'Lỗi in do nguyên vật liệu (Vết bẩn, Bị Gãy, Trầy xước)', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L06', 'QC', 'L06', 'Lỗi in do mực (Bong hột, In lem mực nhòe)', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L05', 'QC', 'L05', 'Lỗi in do khung(Mất nét ,Chữ răng cưa, xì khung...)', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L04', 'QC', 'L04', 'Lỗi in do con người (in lệch)', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L03', 'QC', 'L03', 'Lỗi in do môi trường (Chấm đen, Chấm trắng, Bụi nền, bụi chữ)', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L02', 'QC', 'L02', 'In sai nội dung', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW()),
    ('L01', 'QC', 'L01', 'Lỗi in khác màu', 'adfd41b8-9bd9-4cbc-bfae-6d7b25ea0141', NOW(), NOW());