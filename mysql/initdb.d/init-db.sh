# CreateTable
mysql -u root -psecret < "/docker-entrypoint-initdb.d/01-init.sql"
mysql -u root -psecret WorkTimeRecoder < "/docker-entrypoint-initdb.d/02-table.sql"
mysql -u root -psecret WorkTimeRecoder < "/docker-entrypoint-initdb.d/99-dummy.sql"