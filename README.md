## Philosophy: 
WIR3DENGINE is a simple messageboard engine inspired by the defunct [mebious.co.uk](http://mebious.co.uk/)! 
* Written in Golang! ✅
* No admin panel. ✅
* Everything is done manually by hacking away at the sourcecode and sending raw PSQL queries. ✅
* Looks terrible on mobile. ✅
* Only text and images are allowed. No audio files. ✅
  
![Screenshot_20250428_194106](https://github.com/user-attachments/assets/b58611d3-4437-47af-a5f1-970ea94eb8d7)

## Simple setup:
```
sudo apt update 
sudo apt upgrade 
sudo apt install golang postgresql git imagemagick libmagick++-dev
git clone https://github.com/s0nney/WIR3DENGINE.git
cd WIR3DENGINE
$EDITOR .env
go build .
nohup ./WIR3DENGINE >> WIR3DENGINE.log 2>&1 & 
```

## PSQL setup:
```
-- Ideal table structure -- 

CREATE TABLE posts (
        id BIGSERIAL PRIMARY KEY,
        in_name VARCHAR(70),
        in_text VARCHAR(1024),
        ip VARCHAR(512),
        date_posted TIMESTAMP DEFAULT current_timestamp,
);

CREATE TABLE images (
        id BIGSERIAL PRIMARY KEY,
        filename VARCHAR(35),
        checksum VARCHAR(128),
        ip VARCHAR(512)
        date_posted TIMESTAMP DEFAULT current_timestamp,
);

CREATE TABLE bans (
        id BIGSERIAL PRIMARY KEY,
        ip VARCHAR(512),
        date_posted TIMESTAMP DEFAULT current_timestamp,
);
```

## NGINX configuration: 
```
server {
    server_name site.com www.site.com ;

    location / {
        proxy_pass http://localhost:8080;
    	proxy_set_header X-Real-IP $remote_addr;
    	proxy_set_header X-Forwarded-For $proxy_add_x_forward_for;
    	proxy_set_header X-Forwarded-Proto $scheme;
    	proxy_set_header Host $host;
    }

    location /templates/ {
        autoindex off;
    }

    location /tmp/ {
        autoindex off;
    }

    location ~ /\.ht {
        deny all;
    }

}
```

# API:
```
GET /posts => grabs the 20 most recent posts in JSON 
GET /posts/n => grabs the (n < 100) most recent posts in JSON
```

#### (Best ran on Debian!) 
