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
 
package hevc

import (
	"errors"
    "fmt"
    "io"
    "bufio"
    "strings"
    "strconv"
    "container/list"
)

const POS_NAME  = 0
const POS_TYPE  = 1
const POS_VALUE = 3

type Field struct {
    Name    string
    Type    string
    Value   int
}

type ParameterSet struct {
    Reader      *bufio.Reader
    FieldList   *list.List
}

func (ps *ParameterSet) Parse() (line string, err error)  {
	if ps.Reader==nil{
		line = ""
		err = errors.New("nil Reader\n")
		return line, err
	}

    eof := false
    for !eof {
        line, err = ps.Reader.ReadString('\n')
        if err == io.EOF {
            err = nil
            eof = true
        }else if err != nil {
            return line, err
        }else if strings.HasPrefix(line, "=========") {
            return line, err
        }else{
        	var value	int64
            var field   Field
            fieldstr   := strings.Fields(strings.TrimSpace(line))
            field.Name  = fieldstr[POS_NAME]
            field.Type  = fieldstr[POS_TYPE]
            value, err = strconv.ParseInt(fieldstr[POS_VALUE], 0, 0)
            if err != nil {
                return line, err
            }else{
            	field.Value = int(value)
                ps.FieldList.PushBack(field)
            }
        }
    }

    return
}

func (ps *ParameterSet) GetValue(name string) (value int, err error){
	err = errors.New(name+" not found\n")
	for e := ps.FieldList.Front(); e != nil; e = e.Next() {
        field := e.Value.(Field)

        if field.Name == name {
            value, err = field.Value, nil
            break
        }
    }

    return value, err
}


func (ps *ParameterSet) GetTypeValue(name string) (t string, v int, err error){
	err = errors.New(name+" not found\n")
	for e := ps.FieldList.Front(); e != nil; e = e.Next() {
        field := e.Value.(Field)

        if field.Name == name {
            t, v, err = field.Type, field.Value, nil
            break
        }
    }

    return t, v, err
}


func (ps *ParameterSet) ShowInfo(){
    fmt.Printf("===========================================================================\n",)
    for e := ps.FieldList.Front(); e != nil; e = e.Next() {
        field := e.Value.(Field)
        fmt.Printf("%-62s %-6s: %4d\n", field.Name, field.Type, field.Value)
    }
    //fmt.Printf("===========================================================================\n")
}
