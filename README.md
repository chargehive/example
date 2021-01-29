# Example ChargeHive Payment Page

The example uses a config file named `config.yaml`

### Run

Either run `./launch.sh` or run `go run ./` in the project directory.

#### Features

- **Config file to control features**

  If the config file doesn't exist, one will be created on first run. Missing fields are always populated to show different available options
- **Http server**

  For testing ChargehiveJs on an insecure page
- **Https server**

  for testing ChargehiveJs on a secure page
- **Webhook receiver**

  for watching for chargehive webhooks, always responds with `{"message":"OK"}`
- **Chargehive API caller**

  Used to call `Cancel`, `Capture`, `Refund` for a given charge

### Chargehive API

This may need to be updated from time to time, run the `/chargehive/updateAPI.sh` script to download and build the latest api

### Troubleshooting

- If chargehive loads but `ChargeHive.onInit.then(function (event){})` is never fired, check that the `PaymentAuthCdn` is set correctly in configs. (
  probably https://cdn.paymentauth.me:8823)

### Config Descriptions

```
ProjectId         Name of the chargehive project
PaymentAuthCdn    FQDN and/or port for the PaymentAuth service (usually "https://cdn.paymentauth.me:8823")
PlacementToken    Placement token to authorize the use of the chargehiveJS on a specific domain
Currency          Initial currency
Country           Initial country
HttpListen        Port and/or address for the http webservice  (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
HttpsListen       Port and/or address for the https webservice (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
HttpsCertFilename File path to the pem certificate file
HttpsKeyFilename  File path to the pem key file
ApiHost           Port and/or address for the chargehive api server for requests (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
ApiAccessToken    Token used to authenticate requests to the api server
WebhookListen:    Port and/or address to receive chargehive webhooks (e.g. ":8080"/"localhost:8080"/"0.0.0.0:8080")
```