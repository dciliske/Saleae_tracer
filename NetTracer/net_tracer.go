package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var data_file string
var app_file string
var out_file string
var std bool
var utils_prefix = "m68k-elf-"
var addr2line = utils_prefix + "addr2line"
var addr2line_loading = 1000
var lineLength = 154 // Number of characters per log file line

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
	fmt.Println("Usage:\n\ttracer <data> <app>")
}

func init() {
	flag.StringVar(&out_file, "out", "trace.out", "Destination trace file")
	flag.BoolVar(&std, "std", false, "Output trace to stdout")
	flag.Parse()
	if len(flag.Args()) != 2 {
		printUsage()
		os.Exit(2)
	}
	data_file = flag.Arg(0)
	app_file = flag.Arg(1)
}

func fileCloser(f *os.File) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func ReadyData() (unified []TracePoint) {
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
	count := 0
	_ = scanner.Scan()
	var curr *TracePoint
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		val_len := len(parts[2])
		val = parts[2][val_len-4 : val_len]
		switch count % 34 {
		case 0:
			if count > 0 {
				unified = append(unified, *curr)
			}
			curr = new(TracePoint)
			curr.init()
			fallthrough
		case 1:
			curr.pc += val
		case 2:
			fallthrough
		case 3:
			curr.d0 += val
		case 4:
			fallthrough
		case 5:
			curr.d1 += val
		case 6:
			fallthrough
		case 7:
			curr.d2 += val
		case 8:
			fallthrough
		case 9:
			curr.d3 += val
		case 10:
			fallthrough
		case 11:
			curr.d4 += val
		case 12:
			fallthrough
		case 13:
			curr.d5 += val
		case 14:
			fallthrough
		case 15:
			curr.d6 += val
		case 16:
			fallthrough
		case 17:
			curr.d7 += val
		case 18:
			fallthrough
		case 19:
			curr.a0 += val
		case 20:
			fallthrough
		case 21:
			curr.a1 += val
		case 22:
			fallthrough
		case 23:
			curr.a2 += val
		case 24:
			fallthrough
		case 25:
			curr.a3 += val
		case 26:
			fallthrough
		case 27:
			curr.a4 += val
		case 28:
			fallthrough
		case 29:
			curr.a5 += val
		case 30:
			fallthrough
		case 31:
			curr.a6 += val
		case 32:
			fallthrough
		case 33:
			curr.a7 += val
		}
		count++
	}
	unified = append(unified, *curr)
	return unified

}

func MakeUnified() (unified []TracePoint) {
	return ReadyData()
}

func parseFrameLine(text string) (frame TracePoint, err error) {
	parts := strings.Split(text, " ")
	if len(parts) != 18 {
		return frame, errors.New(fmt.Sprintf("Invalid line - part count: %d", len(parts)))
	}
	frame.pc += parts[0]
	frame.d0 += parts[1]
	frame.d1 += parts[2]
	frame.d2 += parts[3]
	frame.d3 += parts[4]
	frame.d4 += parts[5]
	frame.d5 += parts[6]
	frame.d6 += parts[7]
	frame.d7 += parts[8]
	frame.a0 += parts[9]
	frame.a1 += parts[10]
	frame.a2 += parts[11]
	frame.a3 += parts[12]
	frame.a4 += parts[13]
	frame.a5 += parts[14]
	frame.a6 += parts[15]
	frame.a7 += parts[16]
	return frame, nil
}

func main() {
	fmt.Println(data_file)
	fmt.Println(app_file)

	unified := make([]TracePoint, addr2line_loading)

	fi, err := os.Open(data_file)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer fileCloser(fi)

	fiStat, _ := fi.Stat()
	lineCount := int(fiStat.Size()) / lineLength
	scanner := bufio.NewScanner(fi)
	// Through away the first two readings
	_ = scanner.Scan()
	_ = scanner.Scan()

	fo, err2 := os.Create(out_file)
	if err2 != nil {
		panic(err2)
	}
	// close fi on exit and check for its returned error
	defer fileCloser(fo)

	writer := bufio.NewWriter(fo)
	fmt.Println("Starting output")
	count := 1

	i := 0
	for scanner.Scan() {
		unified[i], err = parseFrameLine(scanner.Text())
		if err != nil {
			continue
		}
		i = (i + 1) % addr2line_loading
		if i == 0 {
			cmd := exec.Command(addr2line, "-e", app_file)
			for j := i; (j < (i + addr2line_loading)) && (j < len(unified)); j++ {
				cmd.Args = append(cmd.Args, unified[j].pc)
			}
			fmt.Printf("Executing addr2line: %d of %d\n",
				count, lineCount/addr2line_loading)
			out, err := cmd.Output()
			if err != nil {
				fmt.Println("Execution error: ", err)
				fmt.Println("CMD: ", cmd.Args)
				continue
			}
			outstr := fmt.Sprintf("%s", out)
			outstr = strings.Trim(outstr, "{}")
			outlist := strings.Split(outstr, "\n")
			for j := i; (j < (i + addr2line_loading)) && (j < len(unified)); j++ {
				unified[j].sourceLine = outlist[j-i]
				frame_str := unified[j].Formatted()
				fmt.Fprintf(writer, "Frame %d of %d\n%s\n", count*addr2line_loading+j, lineCount, frame_str)
			}
			count++
		}
	}
	if i != 0 {
		cmd := exec.Command(addr2line, "-e", app_file)
		for j := i; (j < (i + addr2line_loading)) && (j < len(unified)); j++ {
			cmd.Args = append(cmd.Args, unified[j].pc)
		}
		fmt.Printf("Executing addr2line: %d of %d\n",
			count, lineCount/addr2line_loading)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("Execution error: ", err)
			fmt.Println("CMD: ", cmd.Args)
		} else {
			outstr := fmt.Sprintf("%s", out)
			outstr = strings.Trim(outstr, "{}")
			outlist := strings.Split(outstr, "\n")
			for j := 0; (j < i) && (j < len(outlist)); j++ {
				unified[j].sourceLine = outlist[j]
				frame_str := unified[j].Formatted()
				fmt.Fprintf(writer, "Frame %d of %d\n%s\n", count*addr2line_loading+j, lineCount, frame_str)
			}
			count++
		}
	}
	writer.Flush()

	fmt.Println(len(unified))

}
