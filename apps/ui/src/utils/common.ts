export function getCookie(name: string) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts?.pop()?.split(";").shift();
}

export function renameKey(obj: any, oldKey: string, newKey: string) {
  obj[newKey] = obj[oldKey];
  delete obj[oldKey];
}
