# 50.041-Chord-Project

## Useful Links
[Chord Original Paper](https://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf)  
[50.041 eDimension](https://edimension.sutd.edu.sg/webapps/blackboard/content/listContent.jsp?course_id=_4766_1&content_id=_163052_1)  
[Jira Board](https://csheiden.atlassian.net/jira/software/projects/TC/boards/2)  
[GitHub](https://github.com/jmfan2002/50.041-Chord-Project)  
[Checkpoint 1](https://docs.google.com/document/d/1egYjJqHyvjDxoG8iJARUEBRDdLf52-gz4wCM5NolD2c/edit)

## Useful commands
Install new dependencies and uninstall unused ones
```Bash
go mod tidy
```

Run server or entry node
```Bash
go run ServerNode
go run EntryNode
```

Start the system:
```Bash
./start.sh [NUM_SERVERS]
```
Bring the system down:
```Bash
docker compose -f "./docker-compose.yml" down
```
