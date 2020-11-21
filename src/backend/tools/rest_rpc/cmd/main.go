package main

import (
	"e-clinic/src/backend/tools/rest_rpc/generator"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] == "--help" {
		fmt.Println(`Usage:
Each program argument must be go source file. Fo each of those file there will file created with generated REST
handlers and client for them. Generated file will have the same name is source file but it's prefix will be .service.go.

Example:
./main /home/me/my_file_with_interfaces.go /home/me/my_other_fils_with_interfaces.go'
This will result in creation of /home/me/my_file_with_interfaces.go /home/me/my_other_fils_with_interfaces.go`)
		return
	}
	prefix := flag.String("prefix", "", "prefix for router")
	interfacesFlag := flag.String("interfaces", "", "comma separated list of interfaces (ex. a,b,c)")
	flag.Parse()
	var offset = 1
	if *prefix != "" {
		if !strings.HasPrefix(*prefix, "/") {
			panic("endpoint prefix has to start with /")
		}
		offset += 1
	}

	interfaces := make([]string, 0, 0)
	if interfacesFlag != nil && *interfacesFlag != "" {
		interfaces = strings.Split(*interfacesFlag, ",")
		fmt.Printf("I will generate this interfaces %v \n", interfaces)
		offset += 1
	}

	fmt.Println("files:", os.Args[offset:])
	for _, filename := range os.Args[offset:] {

		fset := token.NewFileSet()
		dir, file := filepath.Split(filename)
		if !strings.HasSuffix(file, ".go") {
			panic(fmt.Errorf("%s is not valid go source file path", filename))
		}
		f, err := parser.ParseFile(fset, filename, nil, 0)
		if err != nil {
			panic(fmt.Errorf("could not parse file because %w", err))
		}
		generated := generator.Generate(filename, f, *prefix, interfaces...)
		outputFilename := filepath.Join(dir, file[:len(file)-3]+".service.gen.go")
		err = ioutil.WriteFile(outputFilename, []byte(generated), 0644)
		if err != nil {
			panic(fmt.Errorf("could not write result because %w", err))
		}
		if err := exec.Command("goimports", "-w", outputFilename).Run(); err != nil {
			panic(fmt.Errorf("goimport panic: %w", err))
		}
		if err := exec.Command("gofmt", "-w", outputFilename).Run(); err != nil {
			panic(fmt.Errorf("gofmt panic: %w", err))
		}
	}
}
