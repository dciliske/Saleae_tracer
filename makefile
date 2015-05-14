.phony: run

all: run

run:
	go run tracer.go Addr_15_39_21_04_2015.txt Data_15_39_21_04_2015.txt app_15_39_21_04_2015.elf

test:
	go run tracer.go Addr_test.txt Data_test.txt app_15_39_21_04_2015.elf > trace.log
