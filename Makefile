build:
	go build -o bin/canvas-sync main.go

record:
	vhs examples/pull_files_demo.tape && vhs examples/update_files_demo.tape