
brew install golang-migrate

migrate -version

migrate create -seq -ext=.sql -dir=./migrations create_movies_table


# In this command:
# The -seq flag indicates that we want to use sequential numbering like 0001, 
# 0002, ... for the migration files (instead of a Unix timestamp, which is the default).
# The -ext flag indicates that we want to give the migration files the extension .sql. 
# The -dir flag indicates that we want to store the migration files in the ./migrations
# directory (which will be created automatically if it doesn’t already exist).
# The name create_movies_table is a descriptive label that we give the migration
# files to signify their contents.

# This will create two files in the migrations directory:
# ./migrations/
#   ├── 000001_create_movies_table.down.sql 
#   └── 000001_create_movies_table.up.sql


# These two new files will be completely empty. Now edit the files with the 
# necessary SQL commands to be executed.


migrate create -seq -ext=.sql -dir=./migrations add_movies_check_constraints



# To execute the migrations:
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
# 1/u create_movies_table (38.19761ms)
# 2/u add_movies_check_constraints (63.284269ms)


# In postgresql run the following:
# \dt
# You should see that the movies table as well as schema_migrations table has been created.
# The schema_migrations table is automatically generated by the migrate tool and 
# used to keep track of which migrations have been applied.


migrate -path=./migrations -database=$GREENLIGHT_DB_DSN version
# 2


# Migrate up or down to a specific version by using the goto command:
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN goto 1


# To roll-back by a specific number of migrations. For example, to rollback the
# most recent migration you would run:
migrate -path=./migrations -database =$EXAMPLE_DSN down 1

# To apply all down migrations
migrate -path=./migrations -database=$EXAMPLE_DSN down


# To force the database version number to 1 you should use the force command like so:
migrate -path=./migrations -database=$EXAMPLE_DSN force 1


# The migrate tool also supports reading migration files from remote sources including 
# Amazon S3 and GitHub repositories. For example:
migrate -source="s3://<bucket>/<path>" -database=$EXAMPLE_DSN up
$ migrate -source="github://owner/repo/path#ref" -database=$EXAMPLE_DSN up
$ migrate -source="github://user:personal-access-token@owner/repo/path#ref" -database=$EXAMPLE_DSN up



migrate create -seq -ext .sql -dir ./migrations add_movies_indexes

migrate create -seq -ext=.sql -dir=./migrations create_users_table

migrate create -seq -ext .sql -dir ./migrations create_tokens_table

migrate create -seq -ext .sql -dir ./migrations add_permissions