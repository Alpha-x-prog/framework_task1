export function loginFromEmail(email = '') {
  const i = email.indexOf('@')
  return i > 0 ? email.slice(0, i) : email
}


