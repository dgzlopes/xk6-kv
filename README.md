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
  console.log(client.get("hello_1"));
  var r = client.viewPrefix("hello");
  for (var key in r) {
      console.log(key,r[key])
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
     script: ../example.js
     output: -

  scenarios: (100.00%) 2 scenarios, 6 max VUs, 10m35s max duration (incl. graceful stop):
           * generator: 1 iterations for each of 5 VUs (maxDuration: 10m0s, exec: generator, gracefulStop: 30s)
           * results: 1 iterations for each of 1 VUs (maxDuration: 10m0s, exec: results, startTime: 5s, gracefulStop: 30s)

INFO[0005] world                                         source=console
INFO[0005] hello_3 world                                 source=console
INFO[0005] hello_4 world                                 source=console
INFO[0005] hello_5 world                                 source=console
INFO[0005] hello_6 world                                 source=console
INFO[0005] hello_1 world                                 source=console
INFO[0005] hello_2 world                                 source=console

running (00m05.0s), 0/6 VUs, 6 complete and 0 interrupted iterations
generator ✓ [======================================] 5 VUs  00m00.0s/10m0s  5/5 iters, 1 per VU
results   ✓ [======================================] 1 VUs  00m00.0s/10m0s  1/1 iters, 1 per VU

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=145.98µs min=34.94µs med=67.46µs max=550.58µs p(90)=321.92µs p(95)=436.25µs
     iterations...........: 6   1.199378/s
     vus..................: 0   min=0 max=0
     vus_max..............: 6   min=6 max=6
```
