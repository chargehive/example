###Setting Environment Variables
The following environment variables can be set
- PAUTH_PLACEMENT_TOKEN - Your Sandbox Placement Token
- PAUTH_PROJECT_ID - Your Sandbox Project ID

### Launch the server

You can launch the server using

    `php -S 127.0.0.1:9180 -t ./src`

You can launch the server with the environment variables using the following command

    `PAUTH_PLACEMENT_TOKEN="enter_token_here" PAUTH_PROJECT_ID="enter_project_id_here" php -S 127.0.0.1:9180 -t ./src`
