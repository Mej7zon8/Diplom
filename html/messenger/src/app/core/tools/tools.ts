export function nullValue<T>(): T | null {
  return null
}

export function isDateValidQuick(v: any) {
  return v != null && v != "0001-01-01T00:00:00Z"
}

export function titleCaseWord(word: string) {
  if (!word) return word;
  return word[0].toUpperCase() + word.substr(1).toLowerCase();
}
