import http from 'k6/http';

export const options = {
    scenarios: {
      constant_request_rate: {
        executor: 'constant-arrival-rate',
        rate: 10000,
        timeUnit: '1s', // 1 iterations per second, 10000 RPS
        duration: '30s',
        preAllocatedVUs: 20, // how large the initial pool of VUs would be
        maxVUs: 40, // if the preAllocatedVUs are not enough, we can initialize more
      },
    },
  };

export default function () {
    http.get(__ENV.URL); // --env URL='http://localhost:8080/anything'
}
