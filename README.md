<h2>Users Endpoints</h2>

```javascript
GET "/users"
POST "/users/register":
    Request:
        Body: {
            "name": "test",
            "email": "test@test.com",
            "password": "1234"
        }

POST "/users/login"
    Request:
        Body: {
            "email": "test@test.com",
            "password": "1234"
        }

GET "/users/{id}"

PUT "/users/update"
    Request:
        Body: {
            "password",
            "name",
            "email"
        }
```

<h2>Movies Endpoint</h2>
<br>
<p>(Todos los metodos necesitan autenticacion)</p>

```javascript
GET "/movies/popular" (Este metodo obtiene las populares de la API)
GET "/movies/popular/internal" (Este metodo obtiene las populares de la DB)
GET "/movies/{id}"

POST "/movies/comment/{movie_id}":
    Request:
        Body:{
            "comment": "Esta muy buena"
        }

PUT "/movies/comment/{comment_id}":
    Request:
        Body: {
            "comment": "Me retracto"
        }

DELETE "/movies/comment/{comment_id}"

GET "/movies/match/comments" (Obtiene todos los comentarios que realizo el usuario registrado)
```
