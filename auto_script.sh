clear
sudo docker stop kademlia_kademliaNodes_1
sudo docker stop kademlia_kademliaNodes_2
sudo docker build --tag kadlab .
sudo docker-compose up
