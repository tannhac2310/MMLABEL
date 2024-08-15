
update production_plans set product_code = s.value
from (
	select
		t1.id as production_plan_id,
		t2.field as field,
		t2.value as value
	from production_plans as t1
	left join custom_fields as t2 on t1.id = t2.entity_id and t2.entity_type = 2
	where t2.field in ('ma_sp')
) as s
where id = s.production_plan_id;

update production_plans set product_name = s.value
from (
	select
		t1.id as production_plan_id,
		t2.field as field,
		t2.value as value
	from production_plans as t1
	left join custom_fields as t2 on t1.id = t2.entity_id and t2.entity_type = 2
	where t2.field in ('ten_sp')
) as s
where id = s.production_plan_id;