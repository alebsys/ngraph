module github.com/alebsys/ngraph

go 1.19

require github.com/prometheus/procfs v0.12.0

require golang.org/x/sys v0.15.0 // indirect

replace github.com/prometheus/procfs v0.12.0 => github.com/alebsys/procfs v0.12.9
