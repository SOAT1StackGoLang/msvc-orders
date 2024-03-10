# How to import the swagger to Postman

1. Open Postman
2. Click on the Import button 
[![Import button](./imgs/postman-swagger-1.png)](./imgs/postman-swagger-1.png)
3. Select the "Import From Link" tab
4. Paste the link to the swagger documentation
    4.1. Example: http://localhost:8080/swagger/doc.json
5. click in the settings button
[![Settings button](./imgs/postman-swagger-2.png)](./imgs/postman-swagger-2.png)
6. At the settings page, click in the Parameter Generation and Enable Optional Parameters toggle
[![Parameter Generation tab](./imgs/postman-swagger-3.png)](./imgs/postman-swagger-3.png)
[[![Enable Optional Parameters](./imgs/postman-swagger-4.png)](./imgs/postman-swagger-4.png)

This will generate the parameters for the requests based on the swagger documentation
