// Poll for new webhook events from the server. Display in console output
const webhookUrl = "/webhooks?from=";
let lastWebhook = 0;
setInterval(() => {
  fetch(webhookUrl + lastWebhook)
    .then(res => res.json())
    .then((logs) => {
      if(logs != null)
      {
        for(const log of logs)
        {
          lastWebhook = (log.Received > lastWebhook) ? log.Received : lastWebhook;
          let j = JSON.parse(atob(log.Data));
          if(j.hasOwnProperty("data"))
          {
            j.data = JSON.parse(j.data);
          }
          console.info("Received Webhook:", j);
        }
      }
    }).catch(err => console.error(err));
}, 20000);