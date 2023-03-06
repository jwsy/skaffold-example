# skaffold-example of a mysql-rest API in Go

1. Requires a K8s cluster, I use Rancher desktop as shown here: https://medium.com/macoclock/rancher-desktop-setup-for-k8s-on-your-macos-laptop-6f1c576ceb48
2. Create an example database using the MySQL Tutorial https://www.mysqltutorial.org/mysql-sample-database.aspx
3. Write a simple Go api

## Example Database `db/`
1. Download the sample database

    This tutorial provides a nice multi-table database named `classicmodels` 
    
    https://www.mysqltutorial.org/how-to-load-sample-database-into-mysql-database-server.aspx

2. Run a mysql container and initialize it

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

    Browse to `localhost:3000` and the API returns the JSONified results of the query `SELECT * FROM employees LIMIT 10`
    
    Example `launch.json`
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
    Example log output
    ```
    $ skaffold debug --port-forward
    ...
    building
    ...
    Starting deploy...
    Helm release skaffold-helm not installed. Installing...
    NAME: skaffold-helm
    LAST DEPLOYED: Mon Mar  6 07:43:45 2023
    NAMESPACE: default
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None
    Waiting for deployments to stabilize...
     - deployment/skaffold-helm is ready.
    Deployments stabilized in 4.261 seconds
    Port forwarding service/skaffold-helm-service in namespace default, remote port 3000 -> http://127.0.0.1:3000
    Listing files to watch...
     - skaffold-helm
    Press Ctrl+C to exit
    Not watching for changes...
    [skaffold-helm] API server listening at: [::]:56268
    [skaffold-helm] 2023-03-06T12:43:48Z warning layer=rpc Listening for remote connections (connections are not authenticated nor encrypted)
    [skaffold-helm]
    [skaffold-helm]  ┌───────────────────────────────────────────────────┐
    [skaffold-helm]  │                   Fiber v2.42.0                   │
    [skaffold-helm]  │               http://127.0.0.1:3000               │
    [skaffold-helm]  │       (bound on host 0.0.0.0 and port 3000)       │
    [skaffold-helm]  │                                                   │
    [skaffold-helm]  │ Handlers ............. 4  Processes ........... 1 │
    [skaffold-helm]  │ Prefork ....... Disabled  PID ................ 12 │
    [skaffold-helm]  └───────────────────────────────────────────────────┘
    [skaffold-helm]
    [install-go-debug-support] Installing runtime debugging support files in /dbg
    [install-go-debug-support] Installation complete
    Port forwarding pod/skaffold-helm-f55ddbdf6-96kgl in namespace default, remote port 56268 -> http://127.0.0.1:56268
    [skaffold-helm] / - Aloha!
    [skaffold-helm] <nil>
    [skaffold-helm] 2023-03-06T12:45:49Z error layer=rpc writing response:write tcp 127.0.0.1:56268->127.0.0.1:53180: use of closed network connection
    [skaffold-helm] / - Aloha!
    ```
    
    Expected output:
    ```bash
    $ curl localhost:3000
    {"employees":[{"employeeNumber":{"Int64":1002,"Valid":true},"lastName":{"String":"Murphy","Valid":true},"firstName":{"String":"Diane","Valid":true},"extension":{"String":"x5800","Valid":true},"email":{"String":"dmurphy@classicmodelcars.com","Valid":true},"officeCode":{"String":"1","Valid":true},"reportsTo":{"Int64":0,"Valid":false},"jobTitle":{"String":"President","Valid":true}},{"employeeNumber":{"Int64":1056,"Valid":true},"lastName":{"String":"Patterson","Valid":true},"firstName":{"String":"Mary","Valid":true},"extension":{"String":"x4611","Valid":true},"email":{"String":"mpatterso@classicmodelcars.com","Valid":true},"officeCode":{"String":"1","Valid":true},"reportsTo":{"Int64":1002,"Valid":true},"jobTitle":{"String":"VP Sales","Valid":true}},{"employeeNumber":{"Int64":1076,"Valid":true},"lastName":{"String":"Firrelli","Valid":true},"firstName":{"String":"Jeff","Valid":true},"extension":{"String":"x9273","Valid":true},"email":{"String":"jfirrelli@classicmodelcars.com","Valid":true},"officeCode":{"String":"1","Valid":true},"reportsTo":{"Int64":1002,"Valid":true},"jobTitle":{"String":"VP Marketing","Valid":true}},{"employeeNumber":{"Int64":1088,"Valid":true},"lastName":{"String":"Patterson","Valid":true},"firstName":{"String":"William","Valid":true},"extension":{"String":"x4871","Valid":true},"email":{"String":"wpatterson@classicmodelcars.com","Valid":true},"officeCode":{"String":"6","Valid":true},"reportsTo":{"Int64":1056,"Valid":true},"jobTitle":{"String":"Sales Manager (APAC)","Valid":true}},{"employeeNumber":{"Int64":1102,"Valid":true},"lastName":{"String":"Bondur","Valid":true},"firstName":{"String":"Gerard","Valid":true},"extension":{"String":"x5408","Valid":true},"email":{"String":"gbondur@classicmodelcars.com","Valid":true},"officeCode":{"String":"4","Valid":true},"reportsTo":{"Int64":1056,"Valid":true},"jobTitle":{"String":"Sale Manager (EMEA)","Valid":true}},{"employeeNumber":{"Int64":1143,"Valid":true},"lastName":{"String":"Bow","Valid":true},"firstName":{"String":"Anthony","Valid":true},"extension":{"String":"x5428","Valid":true},"email":{"String":"abow@classicmodelcars.com","Valid":true},"officeCode":{"String":"1","Valid":true},"reportsTo":{"Int64":1056,"Valid":true},"jobTitle":{"String":"Sales Manager (NA)","Valid":true}},{"employeeNumber":{"Int64":1165,"Valid":true},"lastName":{"String":"Jennings","Valid":true},"firstName":{"String":"Leslie","Valid":true},"extension":{"String":"x3291","Valid":true},"email":{"String":"ljennings@classicmodelcars.com","Valid":true},"officeCode":{"String":"1","Valid":true},"reportsTo":{"Int64":1143,"Valid":true},"jobTitle":{"String":"Sales Rep","Valid":true}},{"employeeNumber":{"Int64":1166,"Valid":true},"lastName":{"String":"Thompson","Valid":true},"firstName":{"String":"Leslie","Valid":true},"extension":{"String":"x4065","Valid":true},"email":{"String":"lthompson@classicmodelcars.com","Valid":true},"officeCode":{"String":"1","Valid":true},"reportsTo":{"Int64":1143,"Valid":true},"jobTitle":{"String":"Sales Rep","Valid":true}},{"employeeNumber":{"Int64":1188,"Valid":true},"lastName":{"String":"Firrelli","Valid":true},"firstName":{"String":"Julie","Valid":true},"extension":{"String":"x2173","Valid":true},"email":{"String":"jfirrelli@classicmodelcars.com","Valid":true},"officeCode":{"String":"2","Valid":true},"reportsTo":{"Int64":1143,"Valid":true},"jobTitle":{"String":"Sales Rep","Valid":true}},{"employeeNumber":{"Int64":1216,"Valid":true},"lastName":{"String":"Patterson","Valid":true},"firstName":{"String":"Steve","Valid":true},"extension":{"String":"x4334","Valid":true},"email":{"String":"spatterson@classicmodelcars.com","Valid":true},"officeCode":{"String":"2","Valid":true},"reportsTo":{"Int64":1143,"Valid":true},"jobTitle":{"String":"Sales Rep","Valid":true}}]}%
    ```
