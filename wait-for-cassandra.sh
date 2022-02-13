#!/bin/sh

set -eu

echo "$(date) Connecting to database ..."

sleep=1
index=0
limit=100

until [ $index -ge $limit ]
do
  nc -z "${DB_HOST}" "${DB_PORT}" && break

  index=$(( index + 1 ))

  echo "$(date) Waiting for database on "${DB_HOST}":"${DB_PORT}"  ($index/$limit) ..."
  sleep $sleep
done

if [ $index -eq $limit ]
then
  echo "$(date) Failed to connect to database, terminating ..."
  exit 1
fi

echo "$(date) Database Server is ready ..."
echo "$(date) Allowing database 60 seconds to complete migrations ..."
sleep 60

./main