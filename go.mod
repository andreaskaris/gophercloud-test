module example.com/m

go 1.18

replace github.com/gophercloud/gophercloud => ../gophercloud

require (
	github.com/gophercloud/gophercloud v0.20.0
	github.com/gophercloud/utils v0.0.0-20220704184730-55bdbbaec4ba
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/klog/v2 v2.70.1
)

require (
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.6 // indirect
)
