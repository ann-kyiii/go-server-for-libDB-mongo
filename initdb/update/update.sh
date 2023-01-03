#!/bin/sh
mongoimport --db=library-app --collection=books --type csv --file /docker-entrypoint-initdb.d/update/updatedata.csv --mode=merge --upsertFields=id --headerline --columnsHaveTypes -u ${MONGO_INITDB_ROOT_USERNAME:-root} -p ${MONGO_INITDB_ROOT_PASSWORD:-password} --authenticationDatabase admin
# 新しい本にborrowerの空配列を追加
/add_borrower
