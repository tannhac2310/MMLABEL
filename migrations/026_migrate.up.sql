alter table public.production_plans add column product_name varchar(500) not null default '';
alter table public.production_plans add column product_code varchar(255) not null default '';

update public.production_plans set product_code = s.value
from (
	select
		t1.id as production_plan_id,
		t2.field as field,
		t2.value as value
	from public.production_plans as t1
	left join public.custom_fields as t2 on t1.id = t2.entity_id and t2.entity_type = 2
	where t2.field in ('ma_sp')
) as s
where id = s.production_plan_id;

update public.production_plans set product_name = s.value
from (
	select
		t1.id as production_plan_id,
		t2.field as field,
		t2.value as value
	from public.production_plans as t1
	left join public.custom_fields as t2 on t1.id = t2.entity_id and t2.entity_type = 2
	where t2.field in ('ten_sp')
) as s
where id = s.production_plan_id;