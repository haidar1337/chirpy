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
```Authorization: Bearer {JWT_TOKEN}```

Response Body:
```json
{
  "id": 1,
  "email": "test@example.com",
  "is_chirpy_red": true
}
```

### POST /api/refresh
Refreshs a stateless JWT token, which persists for 1 hour, using a stateful refresh token, which persists for 60 days

Headers:
```Authorization: Bearer {REFRESH_TOKEN}```

Response Body:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiIxIiwiZXhwIjoxNzI1OTAwMTAzLCJpYXQiOjE3MjU4OTY1MDN9.Tqc8Zow-MplmxfYYfqZCWCyBrzgDWCZU8lWdhNTDUNI"
}
```

### POST /api/revoke
Revokes a refresh token

Headers:
```Authorization: Bearer {REFRESH_TOKEN}```

Response Status Code: 204 No Content

## Chirp Resource

### POST /api/chirps
Creates a new chirp

Headers:
```Authorization: Bearer {JWT_TOKEN}```
Request Body:
```json
{
  "body": "This is a new chirp!"
}
```
Response Body:
```json
{
  "id": 1,
  "body": "This is a new chirp!",
  "author_id": 1
}
```

### GET /api/chirps
Get all chirps in the database

Response Body:
```json
[
  {
    "id": 1,
    "body": "Chirp!",
    "author_id": 1
  },
  {
    "id": 2,
    "body": "Chirpyyy",
    "author_id": 1
  }
]
```

Optional Query Parameters:
```authorId={ID}```, ```sort={desc || asc}```
Example Request:
```{host}:{port?}/api/chirps?authorId=1&sort=desc```

### GET /api/chirp/{CHIRP_ID}
Get a chirp by its ID

Response Body:
```json
{
    "id": 2,
    "body": "Chirpyyy",
    "author_id": 1
}
```

### DELETE /api/chirps/{CHIRP_ID}

Headers:
```Authorization: Bearer {JWT_TOKEN}```

Response Status Code: 204 No Content

