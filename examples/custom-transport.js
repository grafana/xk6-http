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
