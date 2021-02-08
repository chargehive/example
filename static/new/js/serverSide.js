window.addEventListener("load", function () {
  function sendData(formData, action)
  {
    const XHR = new XMLHttpRequest();
    const resultElement = document.getElementById("resultOutput");
    // Define what happens on successful data submission
    XHR.addEventListener("load", function (event) {
      console.log("result event:", event);
      resultElement.innerHTML = event.target.responseText;
    });
    XHR.addEventListener("error", function (event) {
      console.error("Failed to get result:", event);
    });
    XHR.open("POST", action);
    XHR.send(formData);
  }

  for(let i = 0; i < document.forms.length; i++)
  {
    let form = document.forms[i];
    form.addEventListener("submit", function (event) {
      sendData(new FormData(form), form.action);
      event.preventDefault();
    }, false);
  }

  let tabButtons = [].slice.call(document.querySelectorAll('ul.tab-nav li a.button'));
  tabButtons.map(function (button) {
    button.addEventListener('click', function () {
      document.querySelector('li a.active.button').classList.remove('active');
      button.classList.add('active');

      document.querySelector('.tab-pane.active').classList.remove('active');
      document.querySelector(button.getAttribute('href')).classList.add('active');
    })
  })


  // Poll for new webhook events from the server. Display in console output
//  const webhookUrl = "/webhooks?from=";
//  let lastWebhook = 0;
//  setInterval(() => {
//    fetch(webhookUrl + lastWebhook)
//      .then(res => res.json())
//      .then((logs) => {
//        if(logs != null)
//        {
//          for(const log of logs)
//          {
//            lastWebhook = (log.Received > lastWebhook) ? log.Received : lastWebhook;
//            let j = JSON.parse(atob(log.Data));
//            if(j.hasOwnProperty("data"))
//            {
//              j.data = JSON.parse(j.data);
//            }
//            console.info("Received Webhook:", j);
//          }
//        }
//      }).catch(err => console.error(err));
//  }, 2000);
});
