# lightningaddr

A server that allows you to receive payments at `yourname@yourdomain.com` noncustodially.

It will serve the necessary JSON and then use RPC calls to connect to your node and fetch invoices on demand.

You don't have to do anything besides buying a domain and setting up some DNS records. HTTPS will be provided automatically for you.

## Supported Lightning Backends

  - Sparko
  - LND

## Setup Guide

Considering you own the `domain.com` domain, you need to set up these DNS records:

```
A domain.com -> 5.2.67.89
```

### To use with c-lightning and [Sparko](https://github.com/fiatjaf/sparko):
```
TXT _kind.domain.com -> sparko
TXT _host.domain.com -> http(s)://your-ip-or-whatever.com/rpc
TXT _key.domain.com -> key_with_permission_to_method_invoicewithdescriptionhash
```

### To use with LND:
```
TXT _kind.domain.com -> lnd
TXT _host.domain.com -> http(s)://your-ip-or-whatever.com
TXT _macaroon.domain.com -> macaroon_as_base64
```

**Done.** Now you can receive payments at `any_name@domain.com`.

### Optional extras:

If you want to specify a description for the wallet payment screen:
```
TXT _description.domain.com -> free text
```

If you want to specify an image for the wallet payment screen:
```
TXT _image.domain.com -> https://url.to/image
```

If you want to receive comments:
```
TXT _webhook.domain.com -> https://url.to/receive/webhook
```

The webhook will contain a JSON object like `{"comment": "...", "pr": "lnbc...", "amount:": 12345}`, amount in millisatoshis. The webhook is dispatched when an invoice is generated, not when it is paid, since we don't know when (or if) it was paid.

If you use a self-signed certificate and want that to be checked:
```
TXT _cert.domain.com -> -----BEGIN CERTIFICATE----- MIQT...TQIM -----END CERTIFICATE-----
```

If you want to reuse the domain root to redirect to elsewhere (maybe to the `www.` subdomain?):
```
TXT _redirect.domain.com -> Redirects to https://somewhere.else
```

The syntax for these redirects is exactly the same as in [redirect.name](http://redirect.name).
