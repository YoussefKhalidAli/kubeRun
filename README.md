"# kubeRun" 

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kuberun-controller .

docker build -t youssefkali/kuberun-controller:v0.2.0 .
