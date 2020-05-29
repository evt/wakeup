# WakeUp

Hotel wake up service

# Deploy

```
gcloud functions deploy CallRoom --runtime go113 --allow-unauthenticated --trigger-http --project wakeup-278716 --env-vars-file ./.env.cloud.yaml --region europe-west1 --memory 128MB
```