package main

import (
    "fmt"
    "bufio"
    "os"
    "os/exec"
    "log"
    "gohevc/hevc"
)

func main() {
    fmt.Printf("GoHEVC Analyzer %v for HM\n\n", hevc.Version)

    if len(os.Args)==3 {
        var cmd *exec.Cmd
		const refYUV = "temp.yuv"

        if os.Args[2]=="hm80" {
            cmd = exec.Command("TAppDecoder80", "-b", os.Args[1], "-o", refYUV)
        }else if os.Args[2]=="hm90" {
            cmd = exec.Command("TAppDecoder90", "-b", os.Args[1], "-o", refYUV)
        }else if os.Args[2]=="hm91" {
            cmd = exec.Command("TAppDecoder91", "-b", os.Args[1], "-o", refYUV)
        }else{
            fmt.Printf("Unsupported HM Version\n");
            return
        }
        fmt.Printf("Decoding %s", os.Args[1])

        out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		} else{
		    fmt.Printf("%s\n", out)

	        fmt.Printf("\nAnalyzing %s\n", os.Args[1])

	        traceFile, err := os.Open("TraceDec.txt");
	        if err!=nil {
	            log.Fatal(err)
	        }
	        defer traceFile.Close()

	        yuvFile, err := os.Open(refYUV);
	        if err!=nil {
	            log.Fatal(err)
	        }
	        defer yuvFile.Close()

	        var ha hevc.Analyzer

	        if err = ha.ParseTrace(traceFile); err!=nil {
	            log.Fatal(err)
	        }else{
	        	if err = ha.ParseYUV(yuvFile); err !=nil {
	        		log.Fatal(err)
	        	}else{
	            	ha.ParseCmd(os.Args[1], bufio.NewReader(os.Stdin), string(out))
	        	}
	        }
	    }
    }else{
        fmt.Printf("Usage: gohevc.exe bitstream.265 hmversion(hm80/hm90/hm91)\n")
    }
}
