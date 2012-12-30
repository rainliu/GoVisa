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
