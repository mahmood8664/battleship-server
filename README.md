# battleship-server
Battleship server written by Golang.

This project provide APIs for battleship-client project. It needs Mongodb to store data. 

To build Docker image run: 
docker build -t battleship-server .

To run Docker image run: 

docker run --name battleship-server -e BATTLESHIP_MONGODB_URL=mongodb://mongo:27017 -e BATTLESHIP_MONGODB_USERNAME=mongo -e BATTLESHIP_MONGODB_PASSWORD=123456 --network battleship --restart always -d battleship-server

You can play this game at: http://battleship.mamiri.me/
