import kv from 'k6/x/kv';

export const options = {
  scenarios: {
    generator: {
      exec: 'generator',
      executor: 'per-vu-iterations',
      vus: 5,
    },
    results: {
      exec: 'results',
      executor: 'per-vu-iterations',
      startTime: '5s',
      vus: 1,
    },
  },
};

const client = new kv.Client();

export function generator() {
  client.set(`hello_${__VU}`, 'world');
}

export function results() {
  client.delete("hello_1");
  console.log(client.get("hello_1"));
  var r = client.viewPrefix("hello");
  for (var key in r) {
      console.log(key,r[key])
  }
}