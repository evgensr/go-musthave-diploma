curl -v -i -X POST http://localhost:8090/api/user/register \
-H 'Content-Type: application/json' \
-d '{"login":"login","password":"password"}'


curl -v -i -X POST http://localhost:8080/api/user/register \
-H 'Content-Type: application/json' \
-d '{"login":"login2","password":"password2"}'

curl -v -i -X POST http://localhost:8090/api/user/login \
-H 'Content-Type: application/json' \
-d '{"login":"login","password":"password"}'

make ; ./gophermart -a=0.0.0.0:8090


curl -v --cookie "education=MTY1NDYxNDA2MnxEdi1CQkFFQ180SUFBUkFCRUFBQUh2LUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FEYVc1MEJBSUFEZz09fG5Tdx-37-d68qK6_u2A_mUQt2iExyoxNfR7QPMucxRn" -H "Content-Type: text/plain" -d 162124202183 -X POST http://localhost:8080/api/user/orders

curl -v -H "Content-Type: text/plain" -d 162124202183 -X POST http://localhost:8080/api/user/orders

curl -v --cookie "education=MTY1NDk0NzU4N3xEdi1CQkFFQ180SUFBUkFCRUFBQUlQLUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FGYVc1ME5qUUVBZ0FDfAwkjQtJ5DA-meF0CRhpD-UIjLRnlL0ssy3KzGsueX5K" -H "Content-Type: text/plain" -d 162124202183 -X POST http://localhost:8090/api/user/orders

curl -v --cookie "education=MTY1NDYxNDA2MnxEdi1CQkFFQ180SUFBUkFCRUFBQUh2LUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FEYVc1MEJBSUFEZz09fG5Tdx-37-d68qK6_u2A_mUQt2iExyoxNfR7QPMucxRn"   -H "Content-Type: text/plain"  http://localhost:8080/api/user/orders

curl -v --cookie "education=MTY1NDYxNDA2MnxEdi1CQkFFQ180SUFBUkFCRUFBQUh2LUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FEYVc1MEJBSUFEZz09fG5Tdx-37-d68qK6_u2A_mUQt2iExyoxNfR7QPMucxRn"   -H "Content-Type: text/plain"  http://localhost:8080/api/user/balance

make ; ./gophermart -a=0.0.0.0:8090 -r=http://0.0.0.0:8080
