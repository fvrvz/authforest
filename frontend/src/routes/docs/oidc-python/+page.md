---
title: Python (Authlib)
---

# Python Integration

Use [Authlib](https://authlib.org/) — a comprehensive OAuth/OIDC library for Python.

## Install

```bash
pip install authlib httpx
```

## Flask Example

```python
from flask import Flask, redirect, url_for, session
from authlib.integrations.flask_client import OAuth

app = Flask(__name__)
app.secret_key = 'your-secret-key'

oauth = OAuth(app)
oauth.register(
    name='authforest',
    server_metadata_url='http://localhost:8080/.well-known/openid-configuration',
    client_id='your-client-id',
    client_secret='your-client-secret',
    client_kwargs={'scope': 'openid profile email'},
)


@app.route('/login')
def login():
    redirect_uri = url_for('callback', _external=True)
    return oauth.authforest.authorize_redirect(redirect_uri)


@app.route('/callback')
def callback():
    token = oauth.authforest.authorize_access_token()
    userinfo = token.get('userinfo')
    session['user'] = {
        'name': userinfo.get('name'),
        'email': userinfo.get('email'),
        'roles': userinfo.get('roles', []),
    }
    return redirect('/')


@app.route('/')
def home():
    user = session.get('user')
    if user:
        return f"Hello, {user['name']}! Roles: {user['roles']}"
    return '<a href="/login">Login</a>'


@app.route('/logout')
def logout():
    session.pop('user', None)
    return redirect('/')
```

## FastAPI Example

```python
from fastapi import FastAPI, Request
from fastapi.responses import RedirectResponse
from authlib.integrations.starlette_client import OAuth
from starlette.middleware.sessions import SessionMiddleware

app = FastAPI()
app.add_middleware(SessionMiddleware, secret_key='your-secret-key')

oauth = OAuth()
oauth.register(
    name='authforest',
    server_metadata_url='http://localhost:8080/.well-known/openid-configuration',
    client_id='your-client-id',
    client_secret='your-client-secret',
    client_kwargs={'scope': 'openid profile email'},
)


@app.get('/login')
async def login(request: Request):
    redirect_uri = request.url_for('callback')
    return await oauth.authforest.authorize_redirect(request, redirect_uri)


@app.get('/callback')
async def callback(request: Request):
    token = await oauth.authforest.authorize_access_token(request)
    userinfo = token.get('userinfo')
    request.session['user'] = dict(userinfo)
    return RedirectResponse(url='/')


@app.get('/')
async def home(request: Request):
    user = request.session.get('user')
    if user:
        return {'message': f"Hello, {user['name']}", 'roles': user.get('roles', [])}
    return {'message': 'Not authenticated', 'login': '/login'}
```

## Validate Access Tokens

To validate tokens in a resource server:

```python
from authlib.jose import jwt
import httpx

# Fetch JWKS
jwks_url = 'http://localhost:8080/.well-known/jwks.json'
jwks = httpx.get(jwks_url).json()

# Validate token
claims = jwt.decode(access_token, jwks)
claims.validate()
print(claims['sub'], claims['email'])
```
