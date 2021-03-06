#!/bin/bash
source ./env.sh
cd ..
gcloud functions deploy CallRoom \
--runtime go113 \
--allow-unauthenticated \
--trigger-http \
--project $WAKEUP_GC_PROJECT \
--set-env-vars WAKEUP_PG_PROTO=$WAKEUP_PG_PROTO,WAKEUP_PG_ADDR=$WAKEUP_PG_ADDR,WAKEUP_PG_DB=$WAKEUP_PG_DB,WAKEUP_PG_USER=$WAKEUP_PG_USER,WAKEUP_PG_PASSWORD=$WAKEUP_PG_PASSWORD,WAKEUP_GC_PROJECT=$WAKEUP_GC_PROJECT,WAKEUP_GC_PROJECT_LOCATION=$WAKEUP_GC_PROJECT_LOCATION,WAKEUP_CALL_ROOM_ENDPOINT=$WAKEUP_CALL_ROOM_ENDPOINT,WAKEUP_SCHEDULER_LOCATION=$WAKEUP_SCHEDULER_LOCATION,WAKEUP_SCHEDULER_TIMEZONE=$WAKEUP_SCHEDULER_TIMEZONE,WAKEUP_CALL_ENDPOINT=$WAKEUP_CALL_ENDPOINT,WAKEUP_SCHEDULER_MAX_RETRY_COUNT=$WAKEUP_SCHEDULER_MAX_RETRY_COUNT,WAKEUP_SCHEDULER_RETRY_PERIOD=$WAKEUP_SCHEDULER_RETRY_PERIOD \
--region $WAKEUP_GC_PROJECT_LOCATION \
--memory 128MB \
--max-instances 10
cd cmd