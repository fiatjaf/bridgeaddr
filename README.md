# lightningaddr

To use, set up these DNS requests on a domain or subdomain you own:
```
CNAME domain.com -> lightningaddr.fiatjaf.com
```

To use with c-lightning and [Sparko](https://github.com/fiatjaf/sparko):
```
TXT _kind.domain.com -> sparko
TXT _host.domain.com -> http(s)://your-ip-or-whatever.com/rpc
TXT _key.domain.com -> key_with_permission_to_method_invoicewithdescriptionhash
```

To use with LND:
```
TXT _kind.domain.com -> lnd
TXT _host.domain.com -> http(s)://your-ip-or-whatever.com
TXT _macaroon.domain.com -> macaroon_as_base64
```

Optionally if you want to specify a description for the wallet payment screen:
```
TXT _description.domain.com -> free text
```

Optionally if you want to specify an image for the wallet payment screen:
```
TXT _image.domain.com -> https://url.to/image
```

Optionally if you want to receive comments:
```
TXT _webhook.domain.com -> https://url.to/receive/webhook
```

The webhook will contain a JSON object like `{"comment": "...", "pr": "lnbc...", "amount:": 12345}`, amount in millisatoshis. The webhook is dispatched when an invoice is generated, not when it is paid, since we don't know when (or if) it was paid.

Optionally if you use a self-signed certificate and want that to be checked:
```
TXT _cert.domain.com -> -----BEGIN CERTIFICATE----- MIQT...TQIM -----END CERTIFICATE-----
```
