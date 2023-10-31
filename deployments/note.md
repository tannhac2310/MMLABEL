# backup database

```sql
CREATE
SCHEDULE core_schedule_label
FOR
BACKUP INTO
's3://databases/mm-printing?AWS_ENDPOINT=http://s3_storage:9000&AWS_REGION=us-east-1&AWS_ACCESS_KEY_ID=bX983KRuvA7dF8MYthfK&AWS_SECRET_ACCESS_KEY=YOwWxiBiountG3xwcild'
RECURRING '@daily'
FULL BACKUP ALWAYS
WITH SCHEDULE OPTIONS first_run = 'now';
```

# show backups

```sql
SHOW BACKUPS IN  's3://databases/mm-printing?AWS_ENDPOINT=https://mmlabel.buonho.vn&AWS_REGION=us-east-1&AWS_ACCESS_KEY_ID=bX983KRuvA7dF8MYthfK&AWS_SECRET_ACCESS_KEY=YOwWxiBiountG3xwcild'
```

# restore

*before restore, clean database*

```sql
DROP DATABASE IF EXISTS postgres CASCADE;
CREATE DATABASE IF NOT EXISTS postgres;
```

```sql
RESTORE public.*
    FROM 's3://databases/mm-printing/{VERSION}?AWS_ENDPOINT=https://mmlabel.buonho.vn&AWS_REGION=us-east-1&AWS_ACCESS_KEY_ID=bX983KRuvA7dF8MYthfK&AWS_SECRET_ACCESS_KEY=YOwWxiBiountG3xwcild'
WITH into_db = 'postgres';
```