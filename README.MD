- Its not safe to store all our data in only one database, so we will use replica for each database instance 
- btw replica is not backup 
    - replica is used if we have a faliure in our database so clients can still reach to the data 
    - backup in case of human error such as deleting all our data in database accediently

- the magic part is distributing the keys, values, requests over all servers


 

# to run the applciation for the first time with falg of the bolt location to create the my.db file : 
```sh
➜ distigo git:(main) ✗ go run .\cmd\main.go -bolt_db_location=$(PWD)/my.db -shard_name=usa -http_addr=127.0.0.1:8080
2023/11/10 18:26:31 opened successfully ..
```