#!/usr/bin/env bash
# Start NGROK in background
echo "⚡️ Starting ngrok"
ngrok http 3001 > /dev/null &

# waiting for the server to spin up
sleep 1

# Wait for ngrok to be available
while ! nc -z localhost 4040; do
  sleep 1/5 # wait Ngrok to be available
done

# Get NGROK dynamic URL from its own exposed local API
NGROK_REMOTE_URL="$(curl http://localhost:4040/api/tunnels | jq ".tunnels[0].public_url")"
echo "NGROK public address - ${NGROK_REMOTE_URL}."

if test -z "${NGROK_REMOTE_URL}"
then
  echo "❌ ERROR: ngrok doesn't seem to return a valid URL (${NGROK_REMOTE_URL})."
  exit 1
fi

# Trim double quotes from variable
NGROK_REMOTE_URL=$(echo ${NGROK_REMOTE_URL} | tr -d '"')

echo "NGROK public address - ${NGROK_REMOTE_URL}."




# kill all ngrok process
# killall ngrok
