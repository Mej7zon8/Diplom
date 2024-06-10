export function rgbToHsl(r: number, g: number, b: number) {
  r /= 255, g /= 255, b /= 255;
  let max = Math.max(r, g, b), min = Math.min(r, g, b);
  let h = 0, s = 0, l = (max + min) / 2;

  if (max == min) {
    h = s = 0; // achromatic
  } else {
    let d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
    switch (max) {
      case r:
        h = (g - b) / d + (g < b ? 6 : 0);
        break;
      case g:
        h = (b - r) / d + 2;
        break;
      case b:
        h = (r - g) / d + 4;
        break;
    }
    h /= 6;
  }

  return [h, s, l];
}

export function hslToRgb(h: number, s: number, l: number) {
  let r, g, b;

  if (s === 0) {
    r = g = b = l; // achromatic
  } else {
    const hue2rgb = (p: number, q: number, t: number) => {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1 / 6) return p + (q - p) * 6 * t;
      if (t < 1 / 2) return q;
      if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
      return p;
    }

    let q = l < 0.5 ? l * (1 + s) : l + s - l * s;
    let p = 2 * l - q;

    r = hue2rgb(p, q, h + 1 / 3);
    g = hue2rgb(p, q, h);
    b = hue2rgb(p, q, h - 1 / 3);
  }

  return [Math.round(r * 255), Math.round(g * 255), Math.round(b * 255)];
}

export function uuid2Colors(uuid: string | undefined) {
  let v = {"bg": "#d2d2d2", "fg": "gray"}
  if (!uuid)
    return v

  // Normalize the UUID by removing any non-hex characters (like dashes)
  const normalizedUuid = uuid.replace(/[^a-f0-9]/gi, '');

  // Split the UUID into two halves
  const uuidHalfLength = Math.ceil(normalizedUuid.length / 2);
  const firstHalf = normalizedUuid.substr(0, uuidHalfLength);

  // Transform the first half into a color
  const bgInt = parseInt(firstHalf.padEnd(6, '0').substring(0, 6), 16);
  let r = (bgInt >> 16) & 255;
  let g = (bgInt >> 8) & 255;
  let b = bgInt & 255;

  // Convert to HSL and decrease lightness if it's too high
  let [h, s, l] = rgbToHsl(r, g, b);
  if (h < 0.4)
    h = 0.4 + h;
  [r, g, b] = hslToRgb(h, s, l);

  const backgroundColor = '#' + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);

  // Determine if the text color should be black or white based on the brightness of the background color
  const brightness = (r * 299 + g * 587 + b * 114) / 1000;
  const textColor = brightness > 128 ? 'black' : 'white';

  v.bg = backgroundColor
  v.fg = textColor
  return v
}
