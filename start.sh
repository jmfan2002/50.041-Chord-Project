if [ -z "$1" ]; then
    echo "Please pass number of servers"
    exit
fi

NUM_SERVERS=$1
PORTS=4000-$((4000+NUM_SERVERS-1))
sed "s=NUM_SERVERS=$NUM_SERVERS=g;s=PORTS=$PORTS=g" ./docker-compose_template.yml > ./docker-compose.yml

docker compose -f "./docker-compose.yml" up -d --build
