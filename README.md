# mysql-rest

1. Create an example database using the MySQL Tutorial https://www.mysqltutorial.org/mysql-sample-database.aspx
2. Write a simple python api

## Example Database `db/`
1. Download the sample database

    This tutorial provides a nice multi-table database named `classicmodels` 
    
    https://www.mysqltutorial.org/how-to-load-sample-database-into-mysql-database-server.aspx

2. Run a mysql container and initialize it

    * With nerdctl/docker
    
        https://dev.mysql.com/doc/mysql-installation-excerpt/8.0/en/docker-mysql-getting-started.html

        ```bash
        nerdctl run --name=mysql1 -d -v $PWD:/tmp -e MYSQL_ROOT_PASSWORD=root mysql:8.0.32

        # nerdctl logs mysql1 2>&1 | grep GENERATED # to see the passwd if being more secure
        nerdctl exec -it mysql1 mysql -uroot -proot
        nerdctl exec mysql1 -- mysql -uroot -proot < /tmp/mysqlsampledatabase.sql
        nerdctl exec mysql1 -- bash -c "mysql -uroot -proot -e \"show databases;\""
        # this should show the `classicmodels` database
        nerdctl exec mysql1 -- bash -c "mysql -uroot -proot classicmodels  -e \"SELECT * FROM employees LIMIT 10;\""
        ```

    * With helm 
        
        https://stackoverflow.com/questions/66501138/db-migration-in-helm-mysql-initializationfiles-causes-the-pod-to-crash
        
        ```bash
        helm install mysql --set auth.rootPassword=root bitnami/mysql

        kubectl cp mysqlsampledatabase.sql mysql-0:/tmp
        kubectl exec mysql-0 -- sh -c "mysql -uroot -proot < /tmp/mysqlsampledatabase.sql"
        kubectl exec mysql-0 -- mysql -uroot -proot -e "SHOW DATABASES;"
        kubectl exec mysql-0 -- mysql -uroot -proot -e "SHOW TABLES;"
        kubectl exec mysql-0 -- mysql -uroot -proot classicmodels -e "SELECT * FROM employees LIMIT 10;"
        ```

## Simple API `helm-deployment/`
Use Skaffold to create a simple client that connects to the database that you can debug, based off of the official Skaffold examples: https://github.com/GoogleContainerTools/skaffold/tree/bcbdfe043c2f334f919fa2e6ae06aed4a7578486/examples/helm-deployment. Read the README in the `helm-deployment` directory 

Note that as of Feb 2023, Skaffold doesn't seem to support go 1.20 so we use go 1.19 by using
* **go 1.19** in `go.mod` and 
* **golang:1.19-alpine3.17** in the container images in `Dockerfile` 

Get set up with in the `helm-deployment/` directory

1. Install Skaffold, configure it for a local cluster `skaffold  config set --global local-cluster true`
2. Install go locally `brew install go`, and it's OK to have go 1.20 
3. Install Go VS Code Extension and install required modules like `dlv`, use cmd-shift-p to get the modules
4. Install the app dependencies
    ```bash
    go get github.com/gofiber/fiber/v2
    go get github.com/go-sql-driver/mysql
    go mod download
    go mod tidy
    ```
5. Set your go working directory to be `helm-deployment`, note that this step has already been done and `go.work` is checked in.
    ```bash
    go work init
    go work use helm-deployment 
    ```
6. Use `skaffold debug --port-forward` to debug, using `.vscode/launch.json` as shown below. The magic appeared to be to not set the `cwd`
    ```json
    {
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
        "name": "Skaffold debug",
        "type": "go",
        "request": "attach",
        "mode": "remote",
        "port": 56268,
        "host": "127.0.0.1",
        // "substitutePath": [
        //   {
        //       "from": "${workspaceFolder}",
        //       "to": "/code",
        //   },
        // ],
        }
    ]
    }
    ```