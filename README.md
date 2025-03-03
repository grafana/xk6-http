# xk6-http
Work in progress implementation of k6's new HTTP module

This extension goal is to change the way users used to make HTTP connection with k6.

As you can see in [k6 repo](https://github.com/grafana/k6) there is [lots of issues](https://github.com/grafana/k6/issues?q=is%3Aissue%20state%3Aopen%20label%3Anew-http) with old/standard [HTTP API]() implemented in k6 originally.
So we came with an idea to design a complete new HTTP API (you can see [the design document here](https://github.com/grafana/k6/blob/master/docs/design/018-new-http-api.md)).

## Requirements

- [Goland 1.20+](https://go.dev/)
- [Git](https://git-scm.com/)
- [xk6](https://github.com/grafana/xk6) (`go install go.k6.io/xk6/cmd/xk6@latest`)

## Getting started

1. Build the k6 binary:
`make build`

2. Run an example:
`./k6 run ./examples/test.js`

## Usage/Examples

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
***
- Creating a client with custom transport settings, some HTTP options, and making a POST request:
  
  > This example is on todo list and doesn't work now
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
***
- see `examples/` and `docs/user/` dir for more examples
  
  > Some examples are on todo list and are not implemented yet.

## Contribute
If you want to contribute or help with the development of k6, start by reading [CONTRIBUTING.md]().
Before you start coding, it might be a good idea to first discuss your plans and implementation details with the k6 maintainers—especially when it comes to big changes and features.
You can do this in the GitHub issue for the problem you're solving (create one if it doesn't exist).
  
  > **NOTE:** To disclose security issues, refer to [SECURITY.md]().

## Support
To get help, report bugs, suggest features, and discuss k6 with others, refer to [SUPPORT.md]().

## License
xk6-http is distributed under the AGPL-3.0 license.

