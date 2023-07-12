###  **Building and Testing a REST API in GoLang using Gorilla Mux and MySQL and Templating**

# Install Docker own your system

# Build Docker Image by using this command
`docker build -t myapp-image . `

# Run Docker Image by using this command
`docker run -p 8080:8080 myapp-image`

# Open into your Browser 
localhost:8080
### Home page display there

## Test Login API by using this user_name = Admin && password = Admin123$ for Demo
#### Method : POST 
#### URL : localhost:8080/Login
#### Request Body : JSON { "user_name":"Admin""password":"Admin123$"}

## Add new Routes in routes.go file 
## Connect MYSQL add MYSQL server Address in Config/evn_config.json