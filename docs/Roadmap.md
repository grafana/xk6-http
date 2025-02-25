# Current roadmap

This document has been created due to document the roadmap as we develop the features and link related issues and decisions, so everyone can get onboard on what happened before and what's going to happen easily.




## Things we've done:

- Nothing yet so far, we will edit this section as we continue to develop new features or make new decisions. Also it would be great to link related issue to each item.
## Usage/Examples to cover:

### client:
- Using a client with default transport settings, and making a GET request:
```javascript
import { Client } from 'k6/x/net/http';

export default async function () {
  const client = new Client();
  const response = await client.get('https://httpbin.test.k6.io/get');
  const jsonData = await response.json();
  console.log(jsonData);
}
```
- Creating a client with custom transport settings, some HTTP options, and making a POST request:
```javascript
import { TCP } from 'k6/x/net';
import { Client } from 'k6/x/net/http';

export default async function () {
  const client = new Client({
    dial: async address => {
      return await TCP.open(address, { keepAlive: true });
    },
    proxy: 'https://myproxy',
    headers: { 'User-Agent': 'k6' },  // set some global headers
  });
  await client.post('http://10.0.0.10/post', {
    json: { name: 'k6' }, // automatically adds 'Content-Type: application/json' header
  });
}
```
- Configuring TLS with a custom CA certificate and forcing HTTP/2:
```javascript
import { TCP } from 'k6/x/net';
import { Client } from 'k6/x/net/http';
import { open } from 'k6/x/file';

const caCert = await open('./custom_cacert.pem');

export default async function () {
  const client = new Client({
    dial: async address => {
      return await TCP.open(address, {
        tls: {
          alpn: ['h2'],
          caCerts: [caCert],
        }
      });
    },
  });
  await client.get('https://10.0.0.10/');
}
```
- Forcing unencrypted HTTP/2 (h2c):
```javascript
import { TCP } from 'k6/x/net';
import { Client } from 'k6/x/net/http';

export default async function () {
  const client = new Client({
    dial: async address => {
      return await TCP.open(address, { tls: false });
    },
    version: [2],
  });
  await client.get('http://10.0.0.10/');
```
### Host name resolution.
[Read more.](https://github.com/grafana/k6/blob/master/docs/design/018-new-http-api.md#host-name-resolution)
```javascript
import { TCP } from 'k6/x/net';
import dns from 'k6/x/net/dns';

const hosts = {
  'hostA': '10.0.0.10',
  'hostB': '10.0.0.11',
};

export default async function () {
  const socket = await TCP.open('myhost', {
    lookup: async hostname => {
      // Return either the IP from the static map, or do an OS lookup,
      // or fallback to making a DNS query to specific servers.
      return hosts[hostname] || await dns.lookup(hostname) ||
        await dns.resolve(hostname, {
          rrtype: 'A',
          servers: ['1.1.1.1:53', '8.8.8.8:53'],
        });
    },
  });
}
```

### request and responses
HTTP requests can be created declaratively, and sent only when needed. This allows reusing request data to send many similar requests.
```javascript
import { Client, Request } from 'k6/x/net/http';

export default async function () {
  const client = new Client({
    headers: { 'User-Agent': 'k6' },  // set some global headers
  });
  const request = new Request('https://httpbin.test.k6.io/get', {
    // These will be merged with the Client options.
    headers: { 'Case-Sensitive-Header': 'somevalue' },
  });
  const response = await client.get(request, {
    // These will override any options for this specific submission.
    headers: { 'Case-Sensitive-Header': 'anothervalue' },
  });
  const jsonData = await response.json();
  console.log(jsonData);
}
```

### More to be added
