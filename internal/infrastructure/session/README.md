
# Flow in /authorize

## 1. Browser hits /authorize with the session cookie automatically attached.

## 2. h.sessions.Current(r) reads the cookie value → returns sid.

## 3. You pass sid to AuthenticationService.Current(ctx, sid) to load the AuthContext from your session store (memory, Redis, DB).

## 4. If no valid AuthContext → redirect to /login.
