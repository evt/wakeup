# WakeUp

Hotel wake up service

# Deploy

```
gcloud functions deploy GetTrace --allow-unauthenticated --runtime go113 --trigger-http --project businessclass-stage --env-vars-file ./.env.yaml --region europe-west1 --memory 128MB
```