docker build -t zane01/seblak-bombom .
# tanpa menggunakan docker network juga tidak apa-apa
# docker network
docker network create --driver bridge seblak-bombom-network

# app
docker container create --name seblak-bombom-app --network seblak-bombom-net -p 8000:8000 seblak-bombom-img

docker cp seblak-bombom-app:/seblak-bombom .

# database mariadb
docker container create --name seblak-bombom-db --network seblak-bombom-network -e MARIADB_ALLOW_EMPTY_ROOT_PASSWORD=true -p 3306:3306 mariadb

# konek ke mariadb
mysql -h localhost -P 3306 --protocol=tcp -u root

docker run -v /migrations:/migrations --network seblak-bombom-network migrate/migrate
    -path=/migrations/ -database mysql://host.docker.internal:3306/seblak_bombom up

# untuk melihat path direktori kontainer
docker container exec -i -t seblak-bombom-app /bin/sh

# menjalankan compose
docker compose up