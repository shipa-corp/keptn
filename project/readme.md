# port forward shipa service

    kubectl -n keptn port-forward service/shipa-keptn 8081:8080

# send action

    curl -X POST -H "Content-Type: application/cloudevents+json" -d @./project/actions/create.framework.json http://localhost:8081/v1/event
    curl -X POST -H "Content-Type: application/cloudevents+json" -d @./project/actions/update.framework.json http://localhost:8081/v1/event
    
    curl -X POST -H "Content-Type: application/cloudevents+json" -d @./project/actions/create.cluster.json http://localhost:8081/v1/event
    curl -X POST -H "Content-Type: application/cloudevents+json" -d @./project/actions/update.cluster.json http://localhost:8081/v1/event
    curl -X POST -H "Content-Type: application/cloudevents+json" -d @./project/actions/remove.cluster.json http://localhost:8081/v1/event
    
    curl -X POST -H "Content-Type: application/cloudevents+json" -d @./project/actions/create.application.json http://localhost:8081/v1/event
    curl -X POST -H "Content-Type: application/cloudevents+json" -d @./project/actions/deploy.application.json http://localhost:8081/v1/event