clear
sudo docker stop $(sudo docker ps -q --filter ancestor=kadlab:latest )
sudo docker build --tag kadlab .
sudo docker-compose up -d
sudo docker attach kademlia_kademliaNodes_1
