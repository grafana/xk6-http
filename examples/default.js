import { Client } from 'k6/x/net/http';

export default async function () {
  const client = new Client();
  const response = await client.get('https://httpbin.test.k6.io/get');
  // This behavior will be in next PR:
  // const jsonData = await response.json();
  console.log(response);
}
