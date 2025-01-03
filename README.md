**NOTICE**: bridgeaddr.fiatjaf.com was shut down in January 2025.

---

            **bridgeaddr**
  bridge server for lightning addresses

This is a server that allows you to receive payments at `yourname@yourdomain.com` noncustodially (but not fully trustlessly[^trustless]).

It will serve the necessary JSON and then use RPC calls to connect to your node and fetch invoices on demand.

You don't have to do anything besides buying a domain and setting up some DNS records. HTTPS will be provided automatically for you.

# Supported Lightning Backends

  - LND
  - Eclair
  - Sparko
  - Commando
  - LNPay
  - LNbits

# Setup Guide

Considering you own the `domain.com` domain, you need to set up these DNS records:

| Record | Domain Name | Value                  |
|--------|-------------|------------------------|
| CNAME  | domain.com  | bridgeaddr.fiatjaf.com |

## To use with LND:

| Record | Domain Name           | Value                               |
|--------|-----------------------|-------------------------------------|
| TXT    | _kind.domain.com      | lnd                                 |
| TXT    | _host.domain.com      | http(s)://lnd-ip-or-domain.com:port |
| TXT    | _macaroon.domain.com  | invoice_macaroon_as_base64_or_hex   |

It is better to _bake_ a new macaroon with a single authorization to create invoices and nothing else. If you don't know how to do that it's fine to get the built-in "invoices" macaroon.

The host value here must be the address and port to your REST API, not your gRPC API nor your Lightning connection port.

## To use with Eclair:

| Record   | Domain Name        | Value                         |
| -------- | ------------------ | ----------------------------- |
| TXT      | _host.domain.com   | http(s)://eclair-domain.com   |

Follow [instructions here](https://gist.github.com/fiatjaf/8e74740d30763713154de15562e08789#file-exposing-eclair-md) on how to properly expose your Eclair to the external world.

## To use with CLN and Commando

| Record   | Domain Name        | Value                       |
| -------- | ------------------ | --------------------------- |
| TXT      | _kind.domain.com   | commando                    |
| TXT      | _host.domain.com   | node.ip.plus.port:9735      |
| TXT      | _nodeid.domain.com | nodeidlike_02c16cca44562... |
| TXT      | _rune.domain.com   | runeasbase64                |

## To use with CLN and [Sparko](https://github.com/fiatjaf/sparko):

| Record | Domain Name      | Value                                                    |
|--------|------------------|----------------------------------------------------------|
| TXT    | _kind.domain.com | sparko                                                   |
| TXT    | _host.domain.com | http(s)://sparko-ip-or-domain.com                        |
| TXT    | _key.domain.com  | key_with_permission_to_method_invoicewithdescriptionhash |

By default, your Sparko host will be something like http://your.ip:9737.


## To use with [LNPay](https://lnpay.co/):

| Record | Domain Name      | Value        |
|--------|------------------|--------------|
| TXT    | _pak.domain.com  | pak_oooooooo |
| TXT    | _waki.domain.com | waki_ooooooo |

See [keys docs](https://docs.lnpay.co/api/get-started/access-keys) for what "pak" and "waki" mean.

## To use with LNbits:

| Record | Domain Name      | Value                             |
|--------|------------------|-----------------------------------|
| TXT    | _kind.domain.com | lnbits                            |
| TXT    | _host.domain.com | http(s)://lnbits-ip-or-domain.com |
| TXT    | _key.domain.com  | lnbits_invoice_key                |

---

Just setup the records above and it's **done.** Now you can receive payments at `any_name@domain.com`.

# Warning

DNS records are public. Only put "invoice" keys there, never "payment"/"admin" keys.

# IPv6, .onion addresses, Tor, ZeroTier

If your node is listening on Tor, no problem, you can just use .onion addresses on the `_host` entry normally.

Some people have static IPv6 addresses pointing directly to their machines (instead of to their home router). You can use these directly.

If your node doesn't have a public address and it is also not listening on Tor, you can use https://zerotier.com/. It is very easy. Just download it, install it and join the public network `a0cbf4b62a1e645f`, then use the IP you'll be assigned and we will be able to connect.

# Optional extras:

If you want to specify a description for the wallet payment screen:

| Record | Domain Name             | Value     |
|--------|-------------------------|-----------|
| TXT    | _description.domain.com | free text |

If you want to specify an image for the wallet payment screen:

| Record | Domain Name       | Value                |
|--------|-------------------|----------------------|
| TXT    | _image.domain.com | https://url.to/image |

If you want to receive comments or payment notifications (if you don't know where to send these, I recommend https://t.me/incomingnotificationsbot or https://pipedream.com/):

| Record | Domain Name         | Value                          |
|--------|---------------------|--------------------------------|
| TXT    | _webhook.domain.com | https://url.to/receive/webhook |

The webhook will contain a JSON object like `{"comment": "...", "pr": "lnbc...", "amount:": 12345}`, amount in millisatoshis. The webhook is dispatched when an invoice is generated, not when it is paid, since we don't know when (or if) it was paid.

If you use a self-signed certificate and want that to be checked:

| Record | Domain Name      | Value                     |
|--------|------------------|---------------------------|
| TXT    | _cert.domain.com | -----BEGIN CERTIFICATE... |

If you want to reuse the domain root to redirect arbitrary pages to elsewhere (maybe to the `www.` subdomain?)(follows the same interface and rules found in [redirect.name](http://redirect.name)):

| Record | Domain Name          | Value                               |
|--------|----------------------|-------------------------------------|
| TXT    | _redirect.domain.com | Redirects to https://somewhere.else |

---

[^trustless]: bridgeaddr requires you to trust that the server won't just show their invoice instead of yours when someone tries to send you money. The server can do that and effectively steal the payments you receive until you notice that. It cannot however touch the money you have on your wallet ever.
