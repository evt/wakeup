# WakeUp

Hotel wake up service

# Deploy

```
gcloud functions deploy CallRoom --runtime go113 --allow-unauthenticated --trigger-http --project hotel-alarm --env-vars-file ./.env.yaml --region europe-west3 --memory 128MB
```