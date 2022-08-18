[![Actions Status](https://github.com/jfbramlett/go-datastate-builder/workflows/Go/badge.svg)](https://github.com/jfbramlett/go-datastate-builder/actions)

# Database Loader

Database loader is a tool that can be used to extract a subset of data from one database
(like staging) and load into another (like a blank local database).

Instead of trying to make some abstract process to load based on config this is very
much a "programmed" solution where a developer writes the process. An attempt has been
made to provide tooling to make manually writing easier but it still requires an understanding
of the database model and how information relates.

The code is made up of a "driver" program that is configured with a set of loaders,
the loaders represent a logical service in echo system. For instance, there is a
loader that will load catalog content, a loader for merchandising content, a loader
for taxonomy content. While there is a relationship amongst these various loaders
managing it is left up to the individual loaders and how they are ordered.

To add a new service you just need to add a new loader (following the examples of
existing loaders) and then wire it into the driver.

One thing worth noting, the current loaders delete existing content before executing
the load, this is done to ensure fk's to static reference data (like enums) align.

## Running the Loader

To use the loader you need build a blank local database - this is a bit tricky given the state of
db migrations so the best thing to do is to force local changes to our tooling, run the loader, and
then revert the changes.

### Step 1: Modify the docker-compose.db.yml

Modify the docker-compose.db.yml to not use our data volumes - essentially make the file look like:

```go
version: "3.7"
services:
  mysql:
    image: 118139069697.dkr.ecr.us-west-1.amazonaws.com/hub/mysql:5.7.25
    command: mysqld --datadir=/var/lib/mysql-no-volume
    environment:
      - MYSQL_ALLOW_EMPTY_PASSWORD=yes
      - MYSQL_DATABASE=splice_local
    ports:
      - ${SPLICE_PORT_mysql}:3306
```

### Step 2: verify your data state is configured password-less

Make sure your data state is configured without a database password. Check the file `services/config/data/local/setting.json` and make sure the setting for
`data_password` is an empty string.

### Step 3: start the database

Use the regular `svc` commands to start the database:

```shell
svc-up mysql
```

### Step 4: create the `splice_local` database

Connect to the running mysql instance using your tool of choice. Check to see if the database `splice_local` exists, if not then issue this command:

```sql
CREATE DATABASE splice_local;
```

### Step 5: set default charset

We need to make sure the database is configured to use the right default charset or else we have
issues in the migrations. So open your db tool of choice and execute:

```sql
ALTER DATABASE `splice_local` CHARACTER SET utf8 COLLATE utf8_unicode_ci;
```

### Step 6: load the db migrations

Run the migrations against the database

```shell
svc-init-db
svc-api-migrate
```

### Step 7: start database tunnel to staging (or prod) replica instance

Start an ssh tunnel to connect to the staging replica database.

```shell
$SPLICE_PLATFORM/tunnel-staging-replica.sh
```

### Step 8: run the loader

Now you can run the loader providing the necessary args:

```shell
etl --source_db_dsn <your username>:<your pwd>@tcp(localhost:57527)/splice_staging?parseTime=true
```

### Step 9: commit image

Create a snapshot of the image.

```shell
docker commit <container id> 118139069697.dkr.ecr.us-west-1.amazonaws.com/cleanroom-data:cleanroom_localdev_v<increment number>
```

### Step 10: push the image to ECR

Share the image:

```shell
dcocker push 118139069697.dkr.ecr.us-west-1.amazonaws.com/cleanroom-data:cleanroom_localdev_v<increment number>
```

### Step 11: update create example

Modify the `svc-data-create` help to reference this new version. The file is
`service/config/lib/handlers/data/create.js` and should be around line 34.

### Step 12: let people know about the change!

Announce the update to everyone.

### Step 13: revert your local changes to docker-compose.db.yml

Revert your changes to the docker-compose.db.yml to put everything back the way it was.
