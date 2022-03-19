## wg-webpanel
WireGuard web interface  
This is just configuration managment and you need to have WireGuard installed  

## Features
- Self-serve 
- Clean and simple UI
- Generate client configuration QR code and config file
- Calculate client network traffic

## Running

### Using binary file
web default credential   
Username : wgadmin   
Password : wireguardadmin   

**NOTE**   
Please change default password as soon as possible.after running ```wg-webPanel```,create md5 hash of your password and change it in  ```data/userAccount.json```.


Download the binary file from the release and run it with command:  
```
./wg-webPanel
```   
### Configuration
with ```./wg-webPanel -h ``` you can see help   
there is just one argument to specify json configuration file address  
Json configuration file have the following parameteres:
```
server-port: web ui port   
data-directory: directory address to save config files   (default : "9090")
static-dir: frontend static directory address            (default : "./static")
wireguard-config-path: wireguard config full path        (default : "/etc/wireguard/wg0.conf")
log-file-address: wg-webPanel log file address           (default : "./server.log")
enable-https : enable or disable https                   (default : false)
cert-address : if ssl is enabled,address to cert file    (default : "")
cert-key-address : if ssl is enabled,address to key file (default : "")
```

## Screenshot
---
![Screenshot](screen-1.jpg?raw=true)   
  
![Screenshot](screen-2.jpg?raw=true)   
   
![Screenshot](screen-3.jpg?raw=true)   

## Roadmap
---   
This project is still under development and there might be bugs or unexpected behavior   
### Features
The following is a list of in progress features :   

- [ ] SQLite as storage backend
- [ ] Simple client accounting
## Contributing
Please read our Contributor Guide for more information

## License
MIT. See LICENSE

