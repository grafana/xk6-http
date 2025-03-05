import { Client, Request } from 'k6/x/net/http';

export default async function () {
  const client = new Client({
    headers: { 'User-Agent': 'k6' },  // set some global headers
  });
  const request = new Request({
    url: 'https://httpbin.test.k6.io/get', 
    // These will be merged with the Client options.
    headers: { 'User-Agent': 'saniyar' },
  });
  // const response = await client.get(request, {
  //   // These will override any options for this specific submission.
  //   headers: { 'Case-Sensitive-Header': 'anothervalue' },
  // });
  
  // we don't implement override any options yet
  const response = await client.get(request)
  const jsonData = await response.json();
  console.log(jsonData);
}
