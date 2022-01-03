# xk6-kv

This is a [k6](https://go.k6.io/k6) extension using the [xk6](https://github.com/grafana/xk6) system.

| :exclamation: This is a proof of concept, isn't supported by the k6 team, and may break in the future. USE AT YOUR OWN RISK! |
|------|

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Install `xk6`:
  ```shell
  $ go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```shell
  $ xk6 build --with github.com/dgzlopes/xk6-kv@latest
  ```

## Example

```javascript
// script.js
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
      startTime: '1s',
      maxDuration: '2s',
      vus: 1,
    },
    ttl: {
      exec: 'ttl',
      executor: 'constant-vus',
      startTime: '3s',
      vus: 1,
      duration: '2s',
    },
  },
};

const client = new kv.Client();

export function generator() {
  client.set(`hello_${__VU}`, 'world');
  client.setWithTTLInSecond(`ttl_${__VU}`, `ttl_${__VU}`, 5);
}

export function results() {
  console.log(client.get("hello_1"));
  client.delete("hello_1");
  try {
    var keyDeleteValue = client.get("hello_1");
    console.log(typeof (keyDeleteValue));
  }
  catch (err) {
    console.log("empty value", err);
  }
  var r = client.viewPrefix("hello");
  for (var key in r) {
    console.log(key, r[key])
  }
}

export function ttl() {
  try {
    console.log(client.get('ttl_1'));
  }
  catch (err) {
    console.log("empty value", err);
  }
}
```

Result output:

```
$ ./k6 run script.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: example.js
     output: -

  scenarios: (100.00%) 3 scenarios, 7 max VUs, 10m30s max duration (incl. graceful stop):
           * generator: 1 iterations for each of 5 VUs (maxDuration: 10m0s, exec: generator, gracefulStop: 30s)
           * results: 1 iterations for each of 1 VUs (maxDuration: 2s, exec: results, startTime: 1s, gracefulStop: 30s)
           * ttl: 1 looping VUs for 2s (exec: ttl, startTime: 3s, gracefulStop: 30s)

INFO[0001] world                                         source=console
INFO[0001] empty value error in get value with key hello_1  source=console
INFO[0001] hello_12 world                                source=console
INFO[0001] hello_2 world                                 source=console
INFO[0001] hello_7 world                                 source=console
INFO[0001] hello_9 world                                 source=console
INFO[0001] hello_5 world                                 source=console
INFO[0001] hello_8 world                                 source=console
INFO[0001] hello_4 world                                 source=console
INFO[0001] hello_6 world                                 source=console
INFO[0001] hello_10 world                                source=console
INFO[0001] hello_11 world                                source=console
INFO[0001] hello_13 world                                source=console
INFO[0001] hello_14 world                                source=console
INFO[0001] hello_15 world                                source=console
INFO[0001] hello_3 world                                 source=console
INFO[0003] ttl_1                                         source=console
INFO[0005] empty value error in get value with key ttl_1  source=console
INFO[0005] empty value error in get value with key ttl_1  source=console
INFO[0005] empty value error in get value with key ttl_1  source=console
INFO[0005] empty value error in get value with key ttl_1  source=console

running (00m05.0s), 0/7 VUs, 47297 complete and 0 interrupted iterations
generator ✓ [======================================] 5 VUs  00m00.0s/10m0s  5/5 iters, 1 per VU
results   ✓ [======================================] 1 VUs  0.0s/2s         1/1 iters, 1 per VU
ttl       ✓ [======================================] 1 VUs  2s             

     data_received........: 0 B   0 B/s
     data_sent............: 0 B   0 B/s
     iteration_duration...: avg=36.8µs min=15.66µs med=22.64µs max=53.39ms p(90)=79.34µs p(95)=95.68µs
     iterations...........: 47297 9457.107597/s
     vus..................: 1     min=0         max=1
     vus_max..............: 7     min=7         max=7
```
