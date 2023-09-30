build:
	go build -o bin/canvas-sync main.go

record:
	vhs examples/init/run.tape && vhs examples/pull_files/run.tape && vhs examples/update_files/run.tape && vhs examples/pull_videos/run.tape && vhs examples/update_videos/run.tape && vhs examples/view_events/run.tape && vhs examples/view_deadlines/run.tape && vhs examples/view_people/run.tape 