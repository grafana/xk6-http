import { Client } from 'k6/x/net/http';

export default async function () {
  const client = new Client({
    proxy: 'http://127.0.0.1:1080',
    headers: { 'User-Agent': 'saniyar' },  // set some global headers
  });
  const response = await client.get('https://httpbin.test.k6.io/get');
  const jsonData = await response.json();
  console.log(jsonData);
}
