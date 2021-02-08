/* Check ChargeHive JS was found and loaded*/
if(typeof ChargeHive === 'undefined')
{
  console.log("failed to load Chargehive JS!");
  displayCHError("Failed to load", document.getElementById("chScriptEle").src)
}

function displayCHError(title, message)
{
  document.getElementById("ch-failure-alert").style.display = "block";
  document.getElementById("ch-loading-spinner").style.display = "none";
  document.getElementById("ch-failure-alert-msg").innerText = message;
  document.getElementById("ch-failure-alert-title").innerText = title;
}

function showChOverlay()
{
  document.getElementsByClassName("overlay")[0].style.display = "flex";
}

function hideChOverlay()
{
  document.getElementsByClassName("overlay")[0].style.display = "none";
}

const resizeObserver = new ResizeObserver(entries => {
  let overlays = document.getElementsByClassName("overlay");
  for(let ol of overlays)
  {
    ol.style.width = (ol.parentNode.clientWidth - 1) + "px";
    ol.style.height = (ol.parentNode.clientHeight - 1) + "px";
  }
});
resizeObserver.observe(document.getElementById("paymentTabContent"));

function copyOnClick(e)
{
  e.preventDefault();
  let copyText = e.target.innerText;
  let textarea = document.createElement("textarea");
  textarea.textContent = copyText;
  textarea.style.position = "fixed"; // Prevent scrolling to bottom of page in MS Edge.
  document.body.appendChild(textarea);
  textarea.select();
  textarea.setSelectionRange(0, 99999);
  document.execCommand("copy");
  document.body.removeChild(textarea);
}

let counter = 0;

function clientEventAdd(event)
{
  counter++;
  let tbodyRef = document.getElementById('clientEventsTable').getElementsByTagName('tbody')[0];
  let newRow = tbodyRef.insertRow();

  let cellAdder = function (txt, long) {
    let c = newRow.insertCell();
    c.appendChild(document.createTextNode(txt));
    c.classList.add("px-1");
    c.classList.add("small");
    if(long)
    {
      c.classList.add("text-break");
    }
  }

  cellAdder(counter, false);
  cellAdder(event.type, false);
  cellAdder((event.detail === Object(event.detail)) ? JSON.stringify(event.detail) : event.detail, true);

  document.getElementById("clientRightTabContentEvents").scrollTop = newRow.offsetTop;
}

function clientEventClear()
{
  document.getElementById('clientEventsTable').getElementsByTagName('tbody')[0].innerHTML = "";
}