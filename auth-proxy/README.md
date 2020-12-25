1. Create new application in auth0
2. Add callback URL to "Allowed Callback URLs", example:
```
http://localhost:8080/callback
```
3. Configure logout URL, example:
```
http://localhost:8080/logout
```
4. setup follow environment variables:
```
TARGET_URL=http://localhost:3000
AUTH0_DOMAIN=<get value from auth0>
AUTH0_CLIENT_ID=<get value from auth0>
AUTH0_CLIENT_SECRET=<get value from auth0>
CALLBACK_URL=http://generic-bank/callback
```
