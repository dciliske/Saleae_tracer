package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var addr_file string
var data_file string
var app_file string
var out_file string
var utils_prefix = "m68k-elf-"
var addr2line = utils_prefix + "addr2line"
var addr2line_loading = 1000

type TraceData struct {
	stamp, addr, data string
	byteNum           int
}

type TracePoint struct {
	pc, d0, d1, d2, d3, d4, d5, d6, d7, a0, a1, a2, a3, a4, a5, a6, a7 string
	sourceLine                                                         string
}

func (t *TracePoint) init() {
	t.pc = "0x"
	t.d0 = "0x"
	t.d1 = "0x"
	t.d2 = "0x"
	t.d3 = "0x"
	t.d4 = "0x"
	t.d5 = "0x"
	t.d6 = "0x"
	t.d7 = "0x"
	t.a0 = "0x"
	t.a1 = "0x"
	t.a2 = "0x"
	t.a3 = "0x"
	t.a4 = "0x"
	t.a5 = "0x"
	t.a6 = "0x"
	t.a7 = "0x"
	t.sourceLine = ""
}

func (t TracePoint) Formatted() string {
	out := ""
	out += fmt.Sprintf("%s\n", t.sourceLine)
	out += fmt.Sprintf("pc: %-12.12s\n", t.pc)
	out += fmt.Sprintf("D0=%-12.12s", t.d0)
	out += fmt.Sprintf("D1=%-12.12s", t.d1)
	out += fmt.Sprintf("D2=%-12.12s", t.d2)
	out += fmt.Sprintf("D3=%-12.12s\n", t.d3)
	out += fmt.Sprintf("D4=%-12.12s", t.d4)
	out += fmt.Sprintf("D5=%-12.12s", t.d5)
	out += fmt.Sprintf("D6=%-12.12s", t.d6)
	out += fmt.Sprintf("D7=%-12.12s\n", t.d7)
	out += fmt.Sprintf("A0=%-12.12s", t.a0)
	out += fmt.Sprintf("A1=%-12.12s", t.a1)
	out += fmt.Sprintf("A2=%-12.12s", t.a2)
	out += fmt.Sprintf("A3=%-12.12s\n", t.a3)
	out += fmt.Sprintf("A4=%-12.12s", t.a4)
	out += fmt.Sprintf("A5=%-12.12s", t.a5)
	out += fmt.Sprintf("A6=%-12.12s", t.a6)
	out += fmt.Sprintf("A7=%-12.12s\n", t.a7)
	return out
}

func (t TracePoint) String() string {
	fmt_str := "%-12.12[1]s %-12.12[2]s %-12.12[3]s %-12.12[4]s %-12.12[5]s %-12.12[6]s %-12.12[7]s %-12.12[8]s %-12.12[9]s %-12.12[10]s %-12.12[11]s %-12.12[12]s %-12.12[13]s %-12.12[14]s %-12.12[15]s %-12.12[16]s %-12.12[17]s"
	return fmt.Sprintf(fmt_str, t.pc, t.d0, t.d1, t.d2, t.d3,
		t.d4, t.d5, t.d6, t.d7, t.a0, t.a1,
		t.a2, t.a3, t.a4, t.a5, t.a6, t.a7)
}

func printUsage() {
	fmt.Println("Usage:\n\ttracer <addrs> <data> <app>")
}

func init() {
	flag.StringVar(&out_file, "out", "trace.out", "Destination trace file")
	flag.Parse()
	if len(flag.Args()) != 3 {
		printUsage()
		os.Exit(2)
	}
	addr_file = flag.Arg(0)
	data_file = flag.Arg(1)
	app_file = flag.Arg(2)
}

func fileCloser(f *os.File) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func ReadyArgs() []TraceData {
	// open input file
	fi, err := os.Open(addr_file)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer fileCloser(fi)
	// make a read buffer
	scanner := bufio.NewScanner(fi)

	var val string
	byteNum := 0
	var trace []TraceData
	_ = scanner.Scan()
	for scanner.Scan() {
		var data TraceData
		parts := strings.Split(scanner.Text(), ",")
		val_len := len(parts[1])
		val = parts[1][val_len-3 : val_len-1]
		data.stamp = parts[0]
		data.addr = val
		data.byteNum = byteNum
		byteNum = (byteNum + 1) % 4
		trace = append(trace, data)
	}
	return trace
}

func ReadyData(trace []TraceData) []TraceData {
	// open input file
	fi, err := os.Open(data_file)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer fileCloser(fi)
	// make a read buffer
	scanner := bufio.NewScanner(fi)

	var val string
	byteNum := 0
	count := 0
	_ = scanner.Scan()
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		val_len := len(parts[1])
		val = parts[1][val_len-3 : val_len-1]
		trace[count].data = val
		byteNum = (byteNum + 1) % 4
		count++
	}
	return trace

}

func MakeUnified() (unified []TracePoint) {
	trace := ReadyArgs()
	trace = ReadyData(trace)
	unified = make([]TracePoint, len(trace)/4/17)

	var curr *TracePoint
	count := 0
	for _, val := range trace {
		var addr int64
		addr, _ = strconv.ParseInt(val.addr, 16, 32)
		start := int(addr) - val.byteNum
		switch start {
		case 1:
			if val.byteNum == 0 {
				if count == len(unified) {
					unified = append(unified, *new(TracePoint))
				}
				curr = &unified[count]
				curr.init()
			}
			curr.pc += val.data
		case 2:
			curr.d0 += val.data
		case 3:
			curr.d1 += val.data
		case 4:
			curr.d2 += val.data
		case 5:
			curr.d3 += val.data
		case 6:
			curr.d4 += val.data
		case 7:
			curr.d5 += val.data
		case 8:
			curr.d6 += val.data
		case 9:
			curr.d7 += val.data
		case 10:
			curr.a0 += val.data
		case 11:
			curr.a1 += val.data
		case 12:
			curr.a2 += val.data
		case 13:
			curr.a3 += val.data
		case 14:
			curr.a4 += val.data
		case 15:
			curr.a5 += val.data
		case 16:
			curr.a6 += val.data
		case 17:
			curr.a7 += val.data
			if val.byteNum == 3 {
				count++
			}
		}
	}
	return unified
}

func main() {
	fmt.Println(addr_file)
	fmt.Println(data_file)
	fmt.Println(app_file)

	unified := MakeUnified()

	for i := 0; i < len(unified); i += addr2line_loading {
		cmd := exec.Command(addr2line, "-e", app_file)
		for j := i; (j < (i + addr2line_loading)) && (j < len(unified)); j++ {
			cmd.Args = append(cmd.Args, unified[j].pc)
		}
		fmt.Printf("Executing addr2line: %d of %d\n",
			i/addr2line_loading+1, len(unified)/addr2line_loading+1)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("Execution error: ", err)
			fmt.Println("CMD: ", cmd.Args)
			continue
		}
		fmt.Printf("outlen: %d", len(out))
		outstr := fmt.Sprintf("%s", out)
		outstr = strings.Trim(outstr, "{}")
		outlist := strings.Split(outstr, "\n")
		for j := i; (j < (i + addr2line_loading)) && (j < len(unified)); j++ {
			unified[j].sourceLine = outlist[j-i]
		}
	}

	fo, err := os.Create(out_file)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer fileCloser(fo)

	writer := bufio.NewWriter(fo)
	fmt.Println("Starting output")
	count := 1
	for _, val := range unified {
		frame_str := val.Formatted()
		fmt.Fprintf(writer, "Frame %d of %d\n%s\n", count, len(unified), frame_str)
		fmt.Printf("Formatting frame %d of %d\n", count, len(unified))
		count++
	}
	writer.Flush()

	fmt.Println(len(unified))

}
