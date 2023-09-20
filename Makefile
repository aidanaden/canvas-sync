build:
	go build -o bin/canvas-sync main.go

record:
	vhs examples/init.tape && vhs examples/pull_files_demo.tape && vhs examples/update_files_demo.tape && vhs examples/pull_files_help.tape && vhs examples/update_files_help.tape