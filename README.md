# mm-printing backend server

## Overview

Service ports:
 - 4100: Hydra (manage account, user, authenticate of flamingo-group system)
 - 6100: Aurora (printing management, ...)
 - 8082: Swagger UI (api documents)


### Requirements

* Install [Docker Engine](https://docs.docker.com/engine/install/ubuntu/) (remember to [add Docker as a non-root user](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user) when using Linux)
* Install [docker-compose](https://docs.docker.com/compose/install/)

### Run local eviroment

* Prepare all dependecy: mongodb, ...
```sh
make all
```

* Run Hydra server
```sh
make run-hydra
```

Then access [hydra](http://localhost:8082/?url=http://localhost:4100/hydra/api-docs) to see the api documents

* Run Aurora server
```sh
make run-aurora
```

Then access [aurora](http://localhost:8082/?url=http://localhost:6100/aurora/api-docs) to see the api documents

### Run test

* Run both `hydra` and `gezu`
* Run: `make test`

helm uninstall mysql -n databases