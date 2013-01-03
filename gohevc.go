/* The copyright in this software is being made available under the BSD
 * License, included below. This software may be subject to other third party
 * and contributor rights, including patent rights, and no such rights are
 * granted under this license.
 *
 * Copyright (c) 2012-2013, H265.net
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of the H265.net nor the names of its contributors may
 *    be used to endorse or promote products derived from this software without
 *    specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS
 * BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
 * THE POSSIBILITY OF SUCH DAMAGE.
 */
 
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
