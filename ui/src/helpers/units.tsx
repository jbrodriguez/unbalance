const phpRound = (value: number, precision = 0) => {
  const factor = Math.pow(10, precision);
  const tmp = value * factor;
  const roundedTmp = Math.round(tmp);
  return roundedTmp / factor;
};

const toFixedFix = (n: number, prec: number) => {
  const k = Math.pow(10, prec);
  return `${(Math.round(n * k) / k).toFixed(prec)}`;
};

const numberFormat = (
  v: number,
  decimals: number,
  decPoint: string,
  thousandsSep: string,
) => {
  const value = `${v}`.replace(/[^0-9+\-Ee.]/g, '');

  const n = !isFinite(+value) ? 0 : +value;
  const prec = !isFinite(+decimals) ? 0 : Math.abs(decimals);
  const sep = typeof thousandsSep === 'undefined' ? ',' : thousandsSep;
  const dec = typeof decPoint === 'undefined' ? '.' : decPoint;
  let s = [];

  // Fix for IE parseFloat(0.55).toFixed(0) = 0;
  s = (prec ? toFixedFix(n, prec) : `${Math.round(n)}`).split('.');

  if (s[0].length > 3) {
    s[0] = s[0].replace(/\B(?=(?:\d{3})+(?!\d))/g, sep);
  }

  if ((s[1] || '').length < prec) {
    s[1] = s[1] || '';
    s[1] += new Array(prec - s[1].length + 1).join('0');
  }

  return s.join(dec);
};

const k = 1000;
const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

const formatBytes = (bytes: number) => {
  if (bytes === 0) return { value: '0', unit: 'Byte' };

  let base = bytes ? Math.floor(Math.log(bytes) / Math.log(k)) : 0;
  bytes = bytes / Math.pow(k, base);

  let precision =
    bytes >= 100
      ? 0
      : bytes >= 10
        ? 1
        : phpRound(bytes * 100) % 100 === 0
          ? 0
          : 2;

  if (phpRound(bytes, precision) === k) {
    bytes = 1;
    precision = 2;
    base += 1;
  }

  return {
    value: `${numberFormat(bytes, precision, '.', bytes >= 10000 ? ',' : '')}`,
    unit: `${sizes[base]}`,
  };
};

export const humanBytes = (bytes: number) => {
  const { value, unit } = formatBytes(bytes);
  return `${value} ${unit}`;
};
