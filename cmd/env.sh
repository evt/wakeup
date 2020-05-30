#!/bin/bash
export GOOGLE_APPLICATION_CREDENTIALS=./serviceaccount.json
export WAKEUP_PG_PROTO=tcp
export WAKEUP_PG_ADDR=localhost:5432
export WAKEUP_PG_DB=wakeup
export WAKEUP_PG_USER=postgres
export WAKEUP_PG_PASSWORD=postgres
export WAKEUP_GC_PROJECT=wakeup-278716
export WAKEUP_GC_PROJECT_LOCATION=europe-west1
export WAKEUP_CALL_ROOM_ENDPOINT=https://$WAKEUP_GC_PROJECT_LOCATION-$WAKEUP_GC_PROJECT.cloudfunctions.net/CallRoom
export WAKEUP_SCHEDULER_LOCATION=projects/$WAKEUP_GC_PROJECT/locations/$WAKEUP_GC_PROJECT_LOCATION
