import http from 'k6/http';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.2/index.js';

export const options = {
    scenarios: {
      constant_request_rate: {
        executor: 'constant-arrival-rate',
        rate: __ENV.RATE,
        timeUnit: '1s', // RATE iterations per second
        duration: '300s',
        preAllocatedVUs: 300, // how large the initial pool of VUs would be
        maxVUs: 500, // if the preAllocatedVUs are not enough, we can initialize more
      },
    },
  };

export default function () {
    http.get(__ENV.URL); // --env URL='http://localhost:8080/anything'
}

// https://k6.io/docs/using-k6/metrics/reference/
// https://k6.io/docs/results-output/end-of-test/custom-summary/
export function handleSummary(data) {
  delete data.metrics['http_req_duration{expected_response:true}'];
  delete data.metrics['http_req_tls_handshaking'];
  delete data.metrics['iterations'];
  delete data.metrics['vus'];

  return {
    stdout: textSummary(data, {enableColors: true }),
  };
}
