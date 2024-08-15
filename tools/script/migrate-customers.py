import psycopg2
from psycopg2 import sql
import uuid
from datetime import datetime
import sys

# make sure install psycopg2 before run
# pip3 install psycopg2
# python3 tools/script/migrate-customers.py ./tools/script/output.sql

# Database connection parameters
db_config = {
    'dbname': 'postgres',
    'user': 'root',
    'password': '',
    'host': 'localhost',
    'port': '26257'  # usually 5432 for PostgreSQL
}

def process_database(output: str):
    # Connect to the PostgreSQL database
    try:
        conn = psycopg2.connect(**db_config)
        conn.autocommit = True  # Start transaction
        cursor = conn.cursor()

        # Query to fetch the column values from the first table
        fetch_query = "SELECT id, customer_id FROM public.production_orders WHERE TRUE;"
        cursor.execute(fetch_query)
        rows = cursor.fetchall()

        po_customers = {}
        disctinct_customers = {}
        inserts = []
        updates = []
        for row in rows:
            old_id = row[0]  # Assuming the first column is 'id'
            name = row[1]    # Assuming the second column is 'name'

            # Generate a new UUID
            new_id = str(uuid.uuid4())
            po_customers[old_id] = name
            if name not in disctinct_customers:
                disctinct_customers[name] = new_id

        for name, new_id in disctinct_customers.items():
            # Insert into the second table
            inserts.append("\
INSERT INTO public.customers\r\
(id, 'name', tax, code, country, province, address, phone_number, fax, company_website, company_phone, contact_person_name, contact_person_email, contact_person_phone, contact_person_role, note, 'status', created_by, created_at, updated_at, deleted_at)\r\
VALUES ('{id}', '{name}', '{tax}', '{code}', 'Việt Nam', 'Hồ Chí Minh', 'Quận 1', 0, '{fax}', '', 0, '{contact_person_name}', '{contact_person_email}', 0, 'Nhân viên', '', 1, '{created_by}', now(), now(), NULL);"
                           .format(id=new_id, name=name, tax=name + '-tax', code=name + '-code',  fax=name + '-fax', contact_person_name=name, contact_person_email=name + '@emal.com', created_by='5782eebb-1311-4fea-9fd7-c00ad2489318'))
            
        for old_id, name in po_customers.items():
            if name not in disctinct_customers:
                print("error")
                return
            new_id = disctinct_customers[name]

            # Update the original table with the new UUID
            updates.append("\
UPDATE public.production_orders\r\
SET customer_id = '{}'\r\
WHERE id = '{}';" 
                .format(new_id, old_id))

        with open(output, 'w') as f:
            f.write("-- INSERT \n")

            for insert in inserts:
                f.write(insert)
                f.write("\n")
            print("wrote %d insert queries to %s" % (len(inserts), output))
            
            f.write("-- UPDATE \n")
            
            for update in updates:
                f.write(update)
                f.write("\n")
            print("wrote %d update queries to %s" % (len(updates), output))

    except Exception as e:
        print(f"An error occurred: {e}")

    finally:
        cursor.close()
        conn.close()

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("No output define")
    else:
        print("Write result to %s" % (sys.argv[1]))
        process_database(sys.argv[1])

    
