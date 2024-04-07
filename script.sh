docker build -t zane01/seblak-bombom .

docker container create --name seblak-bombom -p 8000:8000 zane01/seblak-bombom

docker cp seblak-bombom:/seblak-bombom .
