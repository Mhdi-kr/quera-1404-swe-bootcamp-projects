import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: 10,
  duration: '10s',
};

export default function() {
  let res = http.post('http://localhost:8080/reservations', JSON.stringify({
    "workspace_id": Math.floor(Math.random() * 10),
    "member_id": Math.floor(Math.random() * 10),
    "starts_at": "2026-02-20T22:30:58+12:00",
    "ends_at": "2026-02-21T22:30:58+12:00"
}));
  check(res, { "status is 201": (res) => res.status === 201 });
  sleep(1);
}
