# Store Auth Sessions in httpOnly Secure Cookies

SMDR uses JWT authentication, but the JWT is stored in an httpOnly secure cookie instead of browser localStorage. This keeps the MVP simple while reducing token exposure to JavaScript if an XSS bug appears; localStorage may still be used only for non-sensitive UI preferences.
