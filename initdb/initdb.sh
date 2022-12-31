#!/bin/sh
mongoimport --db=library-app --collection=books  --type csv --file /docker-entrypoint-initdb.d/initdata.csv --headerline --columnsHaveTypes -u ${MONGO_INITDB_ROOT_USERNAME:-root} -p ${MONGO_INITDB_ROOT_PASSWORD:-password} --authenticationDatabase admin
# borrowerの空配列を追加
/add_borrower
# docker buildの際に権限エラーにならないように権限付与
chmod -R 777 /data/db
