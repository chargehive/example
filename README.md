# Example ChargeHive Payment Page

The example uses a config file named `config.yaml` if this does not exist, then it will be created on first run with some default values. Edit the file with the
rest of the required values.

A webserver will be started with both http and https access on different ports. (all configurable in config.yaml)

### Run

Either run `./launch.sh` or run `go run ./` in the project directory.

### Troubleshooting

- If chargehive loads but `ChargeHive.onInit.then(function (event){})` is never fired, check that the `PaymentAuthCdn` is set correctly in configs. (
  probably https://cdn.paymentauth.me:8823)

### Config settings

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
```