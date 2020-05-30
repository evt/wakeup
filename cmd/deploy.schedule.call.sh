#!/bin/bash
source ./env.sh
cd ..
gcloud functions deploy ScheduleCall --runtime go113 --allow-unauthenticated --trigger-http --project $WAKEUP_GC_PROJECT --env-vars-file ./.env.cloud.yaml --region $WAKEUP_GC_PROJECT_LOCATION --memory 128MB
