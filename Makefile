
purpleair_exporter: main.go collector.go purpleair/purpleair.go
	go build .

.PHONY: clean
clean:
	rm -f purpleair_exporter
