# Example ChargeHive Payment Page

### N.B. Environmental Variables

- PAUTH_CDN - Full URI for PaymentAuth server (almost always https://cdn.paymentauth.me:8823 for local development)
- PAUTH_PLACEMENT_TOKEN - Your Sandbox Placement Token
- PAUTH_PROJECT_ID - Your Sandbox Project ID

### Run With Launch Script

`launch.sh` uses `launch.env` to set the environmental variables required. If the file doesn't exist then it will be created.

1. Run `./launch.sh`
2. If `launch.env` doesn't exist it will be created for you. Modify the values and re-run.

### Run Directly

Set the environmental variables:

    export PAUTH_CDN="https://cdn.paymentauth.me:8823"
    export PAUTH_PLACEMENT_TOKEN="CHANGE-ME"
    export PAUTH_PROJECT_ID="CHANGE-ME"

then start the server with:

    php -S 127.0.0.1:9180 -t ./src

or all as a single command:

    PAUTH_CDN="https://cdn.paymentauth.me:8823" PAUTH_PLACEMENT_TOKEN="CHANGE-ME" PAUTH_PROJECT_ID="CHANGE-ME" php -S 127.0.0.1:9180 -t ./src

### Troubleshooting

- If chargehive loads but `ChargeHive.onInit.then(function (event){})` is never fired, check that the `PAUTH_CDN` is set correctly. (
  probably https://cdn.paymentauth.me:8823)
