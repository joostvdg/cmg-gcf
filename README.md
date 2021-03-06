# cmg-gcf

Google Cloud Function wrapper for CMG

## Create Google Cloud Function

```bash
gcloud functions deploy cmg --entry-point Cmg --runtime go111 --trigger-http
```

## use

### Normal Game

```bash
http https://europe-west1-kearos-gcp.cloudfunctions.net/cmg
```

### 6 Player Game

```bash
http "https://europe-west1-kearos-gcp.cloudfunctions.net/cmg?type=large&max=365&min=156&minr=65&max300=22"
```

## Go mod

Make sure the `go.mod` is compatible with Go 1.11, as this is currently the only supported runtime.

When you create a `.mod` configuration with Go 1.12, it will add `go 1.12` in the `go.mod` file for no apparent reason.
This will fail the deployment of the GCF.

## Resources

* https://medium.com/google-cloud/google-cloud-functions-for-go-57e4af9b10da
* https://itnext.io/writing-google-cloud-functions-in-go-fb711f33459a
