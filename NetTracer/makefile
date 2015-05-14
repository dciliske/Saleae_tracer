logger:
	netcat -l -p 20035 | xxd -c 68 -g 4 > trace.log


post:
	cp ~/nndk/ssh/examples/SSH/SecureSerToEthFactoryApp/SecureSerToEthFactoryApp.map app_$(stamp).map
	cp ~/nndk/ssh/examples/SSH/SecureSerToEthFactoryApp/SecureSerToEthFactoryApp.elf app_$stamp.elf
	m68k-elf-objdump -ds app_$stamp.elf > app_$stamp.ds
	cat trace.log |sed -E 's/^.{9}//' |cut -d\  -f 1-18 > trace_$stamp.log
	go run net_tracer.go -out trace_$stamp.out trace_$stamp.log app_$stamp.elf

view:
	gvim trace_$stamp.out app_$stamp.ds app_$stamp.map

