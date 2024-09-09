# API for Chirpy
This is a dummy API that I made to practice writing web servers in Go. However, if you want to try it, then no problem.

## User Resource
```json
{
  "id": 1,
  "email": "test@example.com",
  "is_chirpy_red": false
}
```

### POST /api/users
Creates a new user

Request Body:
```json
{
  "email": "test@example.com",
  "password": "u21490jf02175j1g09h139htiohaf"
}
```

Response Body:
```json
{
  "id": 1,
  "email": "test@example.com",
  "is_chirpy_red": false
}
```

### POST /api/login
Login to the server

Request Body:
```json
{
  "email": "test@example.com",
  "password": "u21490jf02175j1g09h139htiohaf"
}
```

Response Body:
```json
{
  "id": 1,
  "email": "test@example.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiIxIiwiZXhwIjoxNzI1ODk4NzQwLCJpYXQiOjE3MjU4OTUxNDB9.xAluyBkhGxk_cVp1lLk4wpZlVu4fjb-jw-6EBFByI7A",
  "refresh_token": "39f939f129df5f6580e3725c065f85d4c20a92f133086111427c88309d4e0de2"
}
```

### PUT /api/users
Upgrades a user account to Chirpy Red

Headers:
```Authorization: Bearer {jwtToken}```

Response Body:
```json
{
  "id": 1,
  "email": "test@example.com",
  "is_chirpy_red": true
}
```
