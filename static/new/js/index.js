let CURRENCY = 'USD';
const CART = [];

/* Event handlers */
ChargeHive.addEventListener(ChargeHive.events.CHARGE_ID, function (event) {
  clientEventAdd(event);
  console.info("CHARGE_ID event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_ALL_VALID, function (event) {
  clientEventAdd(event);
  console.info("ON_ALL_VALID event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_AUTOFILL, function (event) {
  clientEventAdd(event);
  console.info("ON_AUTOFILL event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_BLUR, function (event) {
  clientEventAdd(event);
  console.info("ON_BLUR event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_CANCEL, function (event) {
  clientEventAdd(event);
  console.info("ON_CANCEL event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_CHARGE_CREATED, function (event) {
  clientEventAdd(event);
  console.info("ON_CHARGE_CREATED event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_DECLINED, function (event) {
  clientEventAdd(event);
  console.info("ON_DECLINED event:", event.detail)
  setNormal();
})
ChargeHive.addEventListener(ChargeHive.events.ON_ERROR, function (event) {
  clientEventAdd(event);
  console.info("ON_ERROR event:", event.detail)
  setNormal();
  console.error((!event.detail) ? 'Error Without Detail:' : 'Error:', (!event.detail) ? event : event.detail);
  if(event.detail.type && event.detail.type === 'card')
  {
    document.querySelector('.ch-field-placeholder').style.border = '1px solid red';
  }
});
ChargeHive.addEventListener(ChargeHive.events.ON_FIELD_CHANGE, function (event) {
  clientEventAdd(event);
  console.info("ON_FIELD_CHANGE event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_FOCUS, function (event) {
  clientEventAdd(event);
  console.info("ON_FOCUS event:", event.detail)
})
ChargeHive.onInit.then(function (event) {
  clientEventAdd(event);
  hideChOverlay();
  console.info("ON_INIT event:", event.detail)
  ChargeHive.setPaymentMethodType('PLACEMENT_CAPABILITY_CARD_FORM');
  document.querySelectorAll('[capability]').forEach(
    function (e) {
      if(event.detail.Capabilities.indexOf(e.getAttribute('capability')) === -1)
      {
        e.classList.add('unavailable');
      }
    });

  // your merchant reference
  ChargeHive.prepareCharge(randomString("MerchRef_", 12));

  /* customer info */
  ChargeHive.setCustomerInfo({firstName: 'Test', lastName: 'Customer', email: 'test@example.com'});
  ChargeHive.setBillingAddress(
    {
      address1: 'address1',
      city:     'city',
      county:   'state',
      country:  'GB',
      postal:   'zip',
    }
  );

  //Simplify Setup
  ChargeHive.setNameOnCard('John Smith');
  ChargeHive.setCardExpiry(12, 21);
  addToCart('Initial Product', 0.05);
});
ChargeHive.addEventListener(ChargeHive.events.ON_METHOD_TYPE_CHANGE, function (event) {
  clientEventAdd(event);
  console.info('ON_METHOD_TYPE_CHANGE event:', event.detail);
  document.querySelector('input[value="' + event.detail.type + '"]').checked = true;
  document.querySelectorAll('[capability].inputContainer').forEach(
    function (e) {
      if(e.getAttribute('capability') !== event.detail.type)
      {
        e.classList.remove('selected');
      }
    });
  let container = document.querySelector('[capability=' + event.detail.type + '].inputContainer');
  if(container)
  {
    container.classList.add('selected');
  }
});
ChargeHive.addEventListener(ChargeHive.events.ON_PASTE, function (event) {
  clientEventAdd(event);
  console.info("ON_PASTE event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_PENDING, function (event) {
  clientEventAdd(event);
  console.info("ON_PENDING event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_READY, function (event) {
  clientEventAdd(event);
  console.info("ON_READY event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_SUBMIT, function (event) {
  clientEventAdd(event);
  console.info("ON_SUBMIT event:", event.detail)
})
ChargeHive.addEventListener(ChargeHive.events.ON_SUCCESS, function (event) {
  clientEventAdd(event);
  console.info("ON_SUCCESS event:", event.detail)
  document.querySelector('.container').innerText = 'Thanks for your purchase';
});
ChargeHive.addEventListener(ChargeHive.events.ON_TOKEN, function (event) {
  clientEventAdd(event);
  console.info("ON_TOKEN event:", event.detail)
});

ChargeHive.addEventListener(ChargeHive.events.ON_VERIFY, function (event) {
  clientEventAdd(event);
  console.info("ON_VERIFY event:", event.detail)
})

/* Custom event handlers */
document.addEventListener('change', function (e) {
  if(e.target.id === 'currencySelector')
  {
    CURRENCY = e.target.value;
    updateTransaction();
  }
  else if(e.target.matches('[name="method"]'))
  {
    ChargeHive.setPaymentMethodType(e.target.value);
  }
});

ChargeHive.setStyle(
  {
    all:     {
      default:  {color: 'black'},
      complete: {color: 'green'},
      empty:    {color: 'black'},
      invalid:  {color: 'red'},
    },
    cardNum: {
      invalid: {color: 'orange'},
    },
    cardExp: {
      invalid: {color: 'pink'},
    },
    cardCvv: {
      invalid: {color: 'purple'},
    },
  }
);

function addToCart(name, price)
{
  CART.push({name: name, price: price});

  updateTransaction();
}

function updateTransaction()
{
  ChargeHive.clearOrderItems();
  let total = 0;
  for(let i in CART)
  {
    if(CART.hasOwnProperty(i))
    {
      let itemPrice = Math.round(CART[i].price * 100);
      ChargeHive.addOrderItem({name: CART[i].name, unitPrice: itemPrice});
      total += itemPrice;
    }
  }

  ChargeHive.updateCharge({amount: total, currency: CURRENCY});
}

function doIt()
{
  ChargeHive.authorizeCharge();
  setProcessing();
}

function doTokenize()
{
  ChargeHive.tokenizeCard();
}

function setProcessing()
{
  document.querySelector('.ch-field-placeholder').style.border = '';
  document.querySelector('#doIt').setAttribute('disabled', 'disabled');
  document.querySelector('#doIt').innerText = 'Processing...';
}

function setNormal()
{
  document.querySelector('#doIt').removeAttribute('disabled');
  document.querySelector('#doIt').innerText = 'Pay Now';
}

const letters = '0123456789ABCDEF';

function randomString(prefix, chars)
{
  for(let i = 0; i < chars; i++)
  {
    prefix += letters[Math.floor(Math.random() * 16)];
  }
  return prefix
}
