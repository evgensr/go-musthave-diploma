curl -v -i -X POST http://localhost:8080/api/user/register \
-H 'Content-Type: application/json' \
-d '{"login":"login","password":"password"}'
