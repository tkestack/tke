let index = 10000;
const timeLead = 1e12;
export function uuid() {
  return 'app-tke-fe-' + (++index * timeLead + Math.random() * timeLead).toString(36);
}
