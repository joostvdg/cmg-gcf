#!/usr/bin/env bash
# europe-west1 = Belgium
gcloud functions deploy cmg --entry-point Cmg --runtime go111 --trigger-http --memory=128MB --project=kearos-gcp --region=europe-west1 --timeout=15
