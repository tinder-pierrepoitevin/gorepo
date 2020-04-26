package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

func main() {
	// countClassesDeps()
	// testMatchImport()
	packageImportsForDirectory()
}

// input directory
// eg execute go run src/pantsdeps/pantsdeps.go /Users/pierrepoitevin/Documents/palo/backend
// 1. Find and print all the BUILD files in subdirectories
// 2. Print the number of BUILD files found
// 3. Print each dependency encountered with the number of time encountered
// 4. print the number of unique deps
func depsCount() {
	argsWithoutProg := os.Args[1:]
	m := make(map[string]int)
	count := crawlAndMap(argsWithoutProg[0], m)
	fmt.Printf("Number of build found: %d\n", count)
	depsCount := 0
	for key, value := range m {
		//if value > 15 {
		fmt.Println("Key:", key, "Value:", value)
		//}
		depsCount++
	}
	fmt.Printf("Number of distinct deps: %d\n", depsCount)
}

// eg go run src/pantsdeps/pantsdeps.go /Users/pierrepoitevin/Documents/palo/backend/src/java/com/tinder/backend
// print all java class
// print the count of all classes encountered
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

// print all java file found
// print the number of total dependencies (not unique) in all these classes
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

// input: trinmmed is a BUILD content in the string format
// input: m is a modifiable map
// output: modified m with increased value for each key in the dependencies of trimmed
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
	if !strings.HasPrefix(line, "import ") {
		return false
	}
	return matched
}

func matchPackage(line string) bool {
	pattern := "package\\s.*;"
	matched, _ := regexp.MatchString(pattern, line)
	if !strings.HasPrefix(line, "package ") {
		return false
	}
	return matched
}

func packageImportsForDirectory() {
	packageImports(os.Args[1])
}

func packageImports(dir string) {
	// sorted slice of keys
	imports := getAllImports(dir)
	excludes := getPackages(dir)
	// sorted deps for BUILD
	buildDeps := getAllBuildDeps(imports, excludes)
	for _, dep := range buildDeps {
		fmt.Printf("'%v',\n", dep)
	}
}

func getAllImports(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	res := make(map[string]bool)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".java") {
			toParse := dir + "/" + f.Name()
			newImports := parseJavaFileForImports(toParse)
			for _, str := range newImports {
				res[str] = true
			}
		}
	}
	var sortedKeys []string
	for k := range res {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	fmt.Printf("# unique imports: %v\n", len(sortedKeys))
	return sortedKeys
}

func getPackages(dir string) map[string]bool {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	res := make(map[string]bool)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".java") {
			toParse := dir + "/" + f.Name()
			newImports := parseJavaFileForPackages(toParse)
			for _, str := range newImports {
				res[str] = true
			}
		}
	}
	return res
}

func getAllBuildDeps(importKeys []string, excludes map[string]bool) []string {
	res := make(map[string]bool)
	missingKeys := make(map[string]bool)
	for _, importKey := range importKeys {
		newKey := findBuildPackage(importKey, missingKeys, excludes)
		if newKey != "" {
			res[newKey] = true
		}
	}
	var sortedDeps []string
	for k := range res {
		sortedDeps = append(sortedDeps, k)
	}
	sort.Strings(sortedDeps)
	for missingKey := range missingKeys {
		fmt.Printf("Missing key: %v\n", missingKey)
	}
	return sortedDeps
}

func uniqueJavaImportKeysFromArg() {
	uniqueJavaImportKeys(os.Args[1])
}

func uniqueJavaImportKeys(dir string) {
	m := make(map[string]bool)
	crawlJavaUniqueImports(dir, m)
	fmt.Printf("unique keys: %v\n", len(m))

	keyToPackage := make(map[string]string)
	missingMapping := make(map[string]bool)
	for key := range m {
		keyToPackage[key] = findBuildPackage(key, missingMapping, make(map[string]bool))
	}

	// To store the keys in slice in sorted order
	var missingKeys []string
	for k := range missingMapping {
		missingKeys = append(missingKeys, k)
	}
	sort.Strings(missingKeys)
	for _, k := range missingKeys {
		fmt.Printf("Key: %v\n", k)
	}

}

func crawlJavaUniqueImports(dir string, res map[string]bool) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			crawlJavaUniqueImports(dir+"/"+f.Name(), res)
		}
		if strings.HasSuffix(f.Name(), ".java") {
			toParse := dir + "/" + f.Name()
			newImports := parseJavaFileForImports(toParse)
			for _, str := range newImports {
				res[str] = true
			}
		}
	}
}

func parseJavaFileForImports(fileName string) []string {
	res := []string{}
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error reading file!", fileName)
		return []string{}
	}
	reader := bufio.NewReader(file)
	str := ""
	for err == nil {
		str, err = reader.ReadString('\n')
		if err == nil && matchImport(str) {
			if strings.HasPrefix(str, "import static ") {
				res = append(res, strings.TrimPrefix(strings.TrimRight(str, ";\n"), "import static "))
			} else {
				res = append(res, strings.TrimPrefix(strings.TrimRight(str, ";\n"), "import "))
			}
		}
	}
	return res
}

func parseJavaFileForPackages(fileName string) []string {
	res := []string{}
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error reading file!", fileName)
		return []string{}
	}
	reader := bufio.NewReader(file)
	str := ""
	for err == nil {
		str, err = reader.ReadString('\n')
		if err == nil && matchPackage(str) {
			res = append(res, strings.TrimPrefix(strings.TrimRight(str, ";\n"), "package "))
		}
	}
	return res
}

func findBuildPackage(key string, missingMapping map[string]bool, excludes map[string]bool) string {
	arr := strings.Split(key, ".")
	for unicode.IsUpper([]rune(arr[len(arr)-1])[0]) {
		arr = arr[:len(arr)-1]
	}
	nKey := strings.Join(arr, ".")
	if _, ok := excludes[nKey]; ok {
		return ""
	}
	if matchTinderBackend(nKey) {
		return "backend/src/java/" + strings.Join(strings.Split(nKey, "."), "/")
	}
	if matchJava(nKey) {
		// No need to add import in BUILD
		return ""
	}
	if matchJacksonCore(nKey) {
		return "3rdparty/jvm/com/fasterxml/jackson:core"
	}
	if matchHTTPClient(nKey) {
		return "3rdparty/jvm/org/apache/httpcomponents:httpclient"
	}
	if matchLogstach(nKey) {
		return "3rdparty/jvm/net/logstash/logback:logstash-logback-encoder"
	}
	if matchGuava(nKey) {
		return "3rdparty/jvm/com/google/guava"
	}
	if matchTinderUtil(nKey) {
		return "common/src/java/com/tinder/util"
	}
	if matchTinderEnums(nKey) {
		return "common/src/java/com/tinder/enums"
	}
	if matchProtobuf(nKey) {
		return "3rdparty/jvm/com/google/protobuf:protobuf-java"
	}
	if matchLettuce(nKey) {
		return "3rdparty/jvm/biz/paluch:lettuce"
	}
	if matchLang3(nKey) {
		return "3rdparty/jvm/org/apache/commons:commons-lang3"
	}
	if matchFinagleHTTP(nKey) {
		return "3rdparty/jvm/com/twitter/finagle:http"
	}
	if matchElasticsearch(nKey) {
		return "3rdparty/jvm/org/elasticsearch/client:rest"
	}
	if matchGeohash(nKey) {
		return "3rdparty/jvm/ch/hsr/geohash"
	}
	if matchHashIds(nKey) {
		return "3rdparty/jvm/org/hashids"
	}
	if matchHystrix(nKey) {
		return "3rdparty/jvm/com/netflix/hystrix:core"
	}
	if matchNetty(nKey) {
		return "3rdparty/jvm/io/netty"
	}
	if matchDynamoDb(nKey) {
		return "3rdparty/jvm/com/amazonaws:aws-java-sdk-dynamodb"
	}
	if matchPrometheusClient(nKey) {
		return "3rdparty/jvm/io/prometheus:simpleclient"
	}
	missingMapping[nKey] = true
	return ""
}

func matchElasticsearch(key string) bool {
	pattern := "org\\.elasticsearch\\.client\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchFinagleHTTP(key string) bool {
	pattern := "com\\.twitter\\.finagle\\.http\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchLang3(key string) bool {
	pattern := "org\\.apache\\.commons\\.lang3\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchLettuce(key string) bool {
	pattern := "com\\.lambdaworks\\.redis\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchProtobuf(key string) bool {
	pattern := "com\\.google\\.protobuf\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchGuava(key string) bool {
	pattern := "com\\.google\\.common\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchLogstach(key string) bool {
	pattern := "net\\.logstash\\.logback\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchHTTPClient(key string) bool {
	return matchHTTPClientSub(key) || matchHTTPUtil(key)
}

func matchHTTPClientSub(key string) bool {
	pattern := "org\\.apache\\.http\\.client\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchHTTPUtil(key string) bool {
	pattern := "org\\.apache\\.http\\.util\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchJacksonCore(key string) bool {
	return matchJacksonAnnotation(key) || matchJacksonCoreSub(key) || matchJacksonDatabind(key)
}

func matchJacksonAnnotation(key string) bool {
	pattern := "com\\.fasterxml\\.jackson\\.annotation\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchJacksonDatabind(key string) bool {
	pattern := "com\\.fasterxml\\.jackson\\.databind\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchJacksonCoreSub(key string) bool {
	pattern := "com\\.fasterxml\\.jackson\\.core\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchTinderUtil(key string) bool {
	pattern := "com\\.tinder\\.util\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchTinderEnums(key string) bool {
	pattern := "com\\.tinder\\.enums\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchTinderBackend(key string) bool {
	pattern := "com\\.tinder\\.backend\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchJava(key string) bool {
	return matchJavaUtil(key) || matchJavaIO(key) || matchJavaTime(key) || matchJavaNio(key) ||
		matchJavaNet(key) || matchJavax(key)
}

func matchJavaUtil(key string) bool {
	pattern := "java\\.util\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchJavax(key string) bool {
	pattern := "javax\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchJavaNio(key string) bool {
	pattern := "java\\.nio\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchJavaTime(key string) bool {
	pattern := "java\\.time\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchJavaNet(key string) bool {
	pattern := "java\\.net\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchJavaIO(key string) bool {
	pattern := "java\\.io\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchGeohash(key string) bool {
	pattern := "ch\\.hsr\\.geohash\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchHashIds(key string) bool {
	pattern := "org\\.hashids\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchHystrix(key string) bool {
	pattern := "com\\.netflix\\.hystrix\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}
func matchNetty(key string) bool {
	pattern := "org\\.jboss\\.netty\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchDynamoDb(key string) bool {
	pattern := "com\\.amazonaws\\.services\\.dynamodbv2\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
}

func matchPrometheusClient(key string) bool {
	pattern := "io\\.prometheus\\.client\\.*"
	matched, _ := regexp.MatchString(pattern, key)
	return matched
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
	res = matchImport("import static java.util.List;\n")
	fmt.Println("expected: true; actual: ", res)
	res = matchImport("import ch.hsr.geohash.Blah;\n")
	fmt.Println("expected: true; actual: ", res)
	res = matchImport("LOGGER.info(\"success import docs\", keyValue(\"size\", bulkResponse.getItems().length));")
	fmt.Println("expected: false; actual: ", res)
}
