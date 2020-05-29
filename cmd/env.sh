#!/bin/bash
export GOOGLE_APPLICATION_CREDENTIALS=./serviceaccount.json
export PG_PROTO=tcp
export PG_ADDR=localhost:5432
export PG_DB=wakeup
export PG_USER=postgres
export PG_PASSWORD=postgres
export CALL_ROOM_ENDPOINT=https://europe-west1-wakeup-278613.cloudfunctions.net/CallRoom
export SCHEDULER_LOCATION=projects/wakeup-278716/locations/europe-west1
