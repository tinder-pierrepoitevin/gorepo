package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	//countClassesDeps()
	testMatchImport()
}

func depsCount() {
	argsWithoutProg := os.Args[1:]
	m := make(map[string]int)
	count := crawlAndMap(argsWithoutProg[0], m)
	fmt.Printf("Number of build found: %d\n", count)
	depsCount := 0
	for key, value := range m {
		if value > 15 {
			fmt.Println("Key:", key, "Value:", value)
		}
		depsCount++
	}
	fmt.Printf("Number of distinct deps: %d\n", depsCount)
}

func countClasses() {
	argsWithoutProg := os.Args[1:]
	count := crawlJava(argsWithoutProg[0])
	fmt.Printf("Number of classes found: %d\n", count)
}

func crawlJava(dir string) int {

	res := 0
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			res += crawlJava(dir + "/" + f.Name())
		}
		if strings.HasSuffix(f.Name(), ".java") {
			toParse := dir + "/" + f.Name()
			fmt.Println("Found java file: " + toParse)
			res++
		}
	}
	return res
}

func countClassesDeps() {
	argsWithoutProg := os.Args[1:]
	count := crawlJavaDeps(argsWithoutProg[0])
	fmt.Printf("Number of classes deps found: %d\n", count)
}

func crawlJavaDeps(dir string) int {

	res := 0
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			res += crawlJavaDeps(dir + "/" + f.Name())
		}
		if strings.HasSuffix(f.Name(), ".java") {
			toParse := dir + "/" + f.Name()
			fmt.Println("Found java file: " + toParse)
			res += strings.Count(parseJavaFile(toParse), "import ")
		}
	}
	return res
}

func crawl(dir string) int {
	//fmt.Println("Crawl directory " + dir)
	res := 0
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			res += crawl(dir + "/" + f.Name())
		}
		if f.Name() == "BUILD" {
			toParse := dir + "/" + f.Name()
			fmt.Println("Found build file: " + toParse)
			slice := parseDependencies(parseFile(toParse))
			res += len(slice)
		}
	}
	return res
}

func crawlAndMap(dir string, m map[string]int) int {
	//fmt.Println("Crawl directory " + dir)
	res := 0
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			res += crawlAndMap(dir+"/"+f.Name(), m)
		}
		if f.Name() == "BUILD" {
			toParse := dir + "/" + f.Name()
			fmt.Println("Found build file: " + toParse)
			parseDependenciesAndMap(parseFile(toParse), m)
			res++
		}
	}
	return res
}

func parseFile(fileName string) string {

	// Read file and remove all spaces
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("error reading file!", fileName)
		return ""
	}
	s := string(file)
	s = strings.Join(strings.Fields(s), "")

	return s
}

func parseJavaFile(fileName string) string {

	// Read file and remove all spaces
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("error reading file!", fileName)
		return ""
	}
	return string(file)
}

func parseDependencies(trimmed string) []string {
	split := strings.Split(trimmed, "dependencies=[")
	if len(split) < 2 {
		return []string{}
	}
	split = split[1:]
	res := []string{}
	for _, str := range split {
		str = (strings.Split(str, "]"))[0]
		res = append(res, strings.Split(str, ",")...)
	}
	return res
}

func parseDependenciesAndMap(trimmed string, m map[string]int) {
	split := strings.Split(trimmed, "dependencies=[")
	if len(split) < 2 {
		return
	}
	split = split[1:]
	res := []string{}
	for _, str := range split {
		str = (strings.Split(str, "]"))[0]
		res = strings.Split(str, ",")
		for _, key := range res {
			if key != "" {
				m[key]++
			}
		}
	}
}

func matchImport(line string) bool {
	pattern := "import\\s.*;"
	matched, _ := regexp.MatchString(pattern, line)
	return matched
}

func parseJavaFileForImports(fileName string) []string {

	return []string{}
}

func testMatchImport() {
	res := matchImport("importerXYZ = new MyImporter();")
	fmt.Println("expected: false; actual: ", res)
	res = matchImport("importer java.util.List;")
	fmt.Println("expected: false; actual: ", res)
	res = matchImport("importjava.util.List;")
	fmt.Println("expected: false; actual: ", res)
	res = matchImport("import java.util.List")
	fmt.Println("expected: false; actual: ", res)
	res = matchImport("import java.util.List;")
	fmt.Println("expected: true; actual: ", res)
	res = matchImport("import java.util.List;\n")
	fmt.Println("expected: true; actual: ", res)
}
