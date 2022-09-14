#!/bin/bash

until nc -z mongodb "${MONGO_INITDB_PORT}"
do
  sleep "${NETCAT_PORT_CHECK_SLEEP_DURATION}"
done

go version
cat << EOF
   ______              _____
  / ____/________     / ___/______________ _____  ___  _____
 / / __/ ___/ __ \    \__ \/ ___/ ___/ __  / __ \/ _ \/ ___/
/ /_/ / /  / /_/ /   ___/ / /__/ /  / /_/ / /_/ /  __/ /
\____/_/   \____/   /____/\___/_/   \__,_/ .___/\___/_/
                                        /_/
EOF

echo "Getting go dependencies and syncing syncing them with 'dev_vendor' directory ..."
go get ./...
echo "DEV Go Container ready! Keeping DEV Go container alive using tail -f /dev/null"
echo "Use 'docker-compose exec go bash' in project root to access DEV Go container"
tail -f /dev/null