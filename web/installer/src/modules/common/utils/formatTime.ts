export const formatTime = (time: string) => {
  let re = time;
  if (time.includes('.')) {
    let splits = time.split('.');
    re = splits[0];
  }
  return re;
};
