#!/bin/sh

ngrok http 8001 > /dev/null &
sleep 1

while ! nc -z localhost 4040; do
  sleep 1 # wait Ngrok to be available
done

NGROK_REMOTE_URL="$(curl http://localhost:4040/api/tunnels | jq ".tunnels[0].public_url")"

if test -z "${NGROK_REMOTE_URL}"
then
  echo "ERROR: ngrok doesn't seem to return a valid URL (${NGROK_REMOTE_URL})."
  exit 1
fi

NGROK_REMOTE_URL=$(echo ${NGROK_REMOTE_URL} | tr -d '"')
echo "NGROK public address - ${NGROK_REMOTE_URL}."

#if "${NGROK_REMOTE_URL}" == 'NULL'
#then
#  echo "ERROR: NGROK public url is NULL."
#fi

# remove current ngrok remote url
grep -v '^public' issuer/issuer_config.default.yaml > temp && mv temp issuer/issuer_config.default.yaml

# add new public url to config
to_add="public_url: ${NGROK_REMOTE_URL}"
echo "$to_add" >> issuer/issuer_config.default.yaml