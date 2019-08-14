<?php
$template = file_get_contents('index.html');

function replaceEnv($var, $template)
{
  return str_replace($var . '_HERE', getenv('PAUTH_' . $var) ?: $var . '_HERE', $template);
}

$template = replaceEnv('PLACEMENT_TOKEN', $template);
$template = replaceEnv('PROJECT_ID', $template);

if(!empty(getenv('PAUTH_CDN')))
{
  $template = str_replace('https://cdn.paymentauth.com', getenv('PAUTH_CDN'), $template);
}

$cards = '<div id="card-templates">CARDS</div>';

const VI = 'Visa';
const MC = 'Mastercard';
const AM = 'American Express';

$brands = [VI, MC, AM];

$cards = [
  'Sandbox'                     => [
    'Success'                   => [VI => '8235 0000 0000 0000', MC => '8235 0000 0000 0042'],
    'Reject'                    => [
      VI => '8235 0000 0000 0208',
      MC => '8235 0000 0000 0240',
      AM => '8235 0000 0000 0281 ',
    ],
    '3ds2 Identify'             => [VI => '8235 0100 0000 0009', MC => '8235 0100 0000 0041'],
    '3ds2 Challenge'            => [VI => '8235 0200 0000 0008', MC => '8235 0200 0000 0040'],
    '3ds2 Identify & Challenge' => [VI => '8235 0300 0000 0007', MC => '8235 0300 0000 0049'],
  ],
  'Paysafe'                     => [
    'Success'                          => [VI => '4037112233000001', MC => '5100400000000000', AM => '370123456789017'],
    '3ds2 Identify'                    => [VI => '4000000000001000', MC => '5200000000001005', AM => '340000000001007'],
    '3ds2 Identify & Fail'             => [VI => '4000000000001067', MC => '5200000000001047'],
    '3ds2 Identify & Challenge'        => [VI => '4000000000001091', MC => '', AM => '340000000001098'],
    '3ds2 Identify & Challenge & Fail' => [MC => '5200000000001104'],
    '3ds1 Challenge'                   => [VI => '4000000000000002', MC => '5200000000000007', AM => '340000000003961'],
  ],
  'Braintree (Expiry: 01/20**)' => [
    'Success'      => [VI => '4000000000001000', MC => '5200000000001005', AM => '340000000001007'],
    'Reject'       => [VI => '4000000000001018', MC => '5200000000001013', AM => '340000000001015'],
    'Bypass'       => [VI => '4000000000001083', MC => '5200000000001088', AM => '340000000001080'],
    '3ds2 Success' => [VI => '4000000000001091', MC => '5200000000001096', AM => '340000000001098'],
  ],
];

$cardTable = '<table>';

foreach($cards as $provider => $types)
{
  $cardTable .= '<tr>';
  $cardTable .= '<th colspan="4" class="provider">' . $provider . '</th>';
  $cardTable .= '</tr>';
  $cardTable .= '<tr>';
  $cardTable .= '<th>Result</th>';
  foreach($brands as $brand)
  {
    $cardTable .= '<th>' . $brand . '</th>';
  }
  $cardTable .= '</tr>';
  foreach($types as $type => $cardNumbers)
  {
    $cardTable .= '<tr>';
    $cardTable .= '<td>' . $type . '</td>';
    foreach($brands as $brand)
    {
      $card = $cardNumbers[$brand] ?? '-';
      $card = str_replace(' ', '', $card);
      $card = implode('<span class="cs"></span>', str_split($card, 4));
      $cardTable .= '<td class="test-card-number">' . $card . '</td>';
    }
    $cardTable .= '</tr>';
  }
}

$cardTable .= '</table>';

echo str_replace('<div id="card-templates"></div>', '<div id="card-templates">' . $cardTable . '</div>', $template);
