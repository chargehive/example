# Example ChargeHive Payment Page

### Environmental Variables

These must be set in either case:

- PAUTH_CDN - Full URI for PaymentAuth server (almost always https://cdn.paymentauth.me:8823 for local development)
- PAUTH_PLACEMENT_TOKEN - Your Sandbox Placement Token
- PAUTH_PROJECT_ID - Your Sandbox Project ID

### Run With Launch Script

Modify `PAUTH_CDN`, `PAUTH_PLACEMENT_TOKEN`, `PAUTH_PROJECT_ID` **in** `launch.sh` 
to the correct values for your setup. Then run:

    ./launch.sh

### Run Directly

Set the environmental variables:

    export PAUTH_CDN="https://cdn.paymentauth.me:8823"
    export PAUTH_PLACEMENT_TOKEN="CHANGE-ME"
    export PAUTH_PROJECT_ID="CHANGE-ME2"

then start the server with:

    php -S 127.0.0.1:9180 -t ./src

or all as a single command:

    PAUTH_CDN="https://cdn.paymentauth.me:8823" PAUTH_PLACEMENT_TOKEN="CHANGE-ME" PAUTH_PROJECT_ID="CHANGE-ME2" php -S 127.0.0.1:9180 -t ./src

### Troubleshooting

- If chargehive loads but `ChargeHive.onInit.then(function (event){})` is never fired, check that the `PAUTH_CDN` is set correctly. (probably https://cdn.paymentauth.me:8823)
