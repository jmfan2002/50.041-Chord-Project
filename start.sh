if [ -z "$1" ]; then
    echo "ERROR: Please pass number of servers"
    exit
fi

NUM_SERVERS=$1
PORTS=4000-$((4000+NUM_SERVERS-1))
sed "s=NUM_SERVERS=$NUM_SERVERS=g;s=PORTS=$PORTS=g" ./docker-compose_template.yml > ./docker-compose.yml

# optional: remove -d --build to keep docker compose active in the terminal,
#       to show live container logs
docker compose -f "./docker-compose.yml" -p chord-network up -d --build
