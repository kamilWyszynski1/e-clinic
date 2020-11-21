package main

import (
	"e-clinic/src/drugs"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {
	f, err := os.Open("drugs.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	byteValue, _ := ioutil.ReadAll(f)

	var drugs drugs.ProduktyLecznicze

	if err := xml.Unmarshal(byteValue, &drugs); err != nil {
		panic(err)
	}

	//for _, p := range drugs.ProduktLeczniczy {
	//	fmt.Printf("Drug: %s\n", p.NazwaProduktu)
	//	for _, s := range p.SubstancjeCzynne.SubstancjaCzynna {
	//		fmt.Printf("\tSubs: %s\n", s)
	//	}
	//	fmt.Println("#########3")
	//}

	// configForNeo4j35 := func(conf *neo4j.Config) {}
	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }

	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("drug", "drug", ""), configForNeo4j40)
	if err != nil {
		panic(err)
	}
	// handle driver lifetime based on your application lifetime requirements
	// driver's lifetime is usually bound by the application lifetime, which usually implies one driver instance per application
	defer driver.Close()

	// For multidatabase support, set sessionConfig.DatabaseName to requested database
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count := 0
	for _, p := range drugs.ProduktLeczniczy {
		if len(p.SubstancjeCzynne.SubstancjaCzynna) == 0 {
			continue
		}
		if count >= 10000 {
			return
		}
		if p.NazwaProduktu == "" || len(p.SubstancjeCzynne.SubstancjaCzynna) == 0 {
			continue
		}
		_, err := session.Run("MERGE (d:Drug { name: $name }) ", map[string]interface{}{
			"name": p.NazwaProduktu,
		})
		if err != nil {
			panic(err)
		}

		for _, s := range p.SubstancjeCzynne.SubstancjaCzynna {
			_, err := session.Run("MERGE (s:Sub { name: $sub })", map[string]interface{}{
				"sub": strings.Trim(s, " "),
			})
			if err != nil {
				panic(err)
			}

			_, err = session.Run("MATCH (s:Sub), (d:Drug) WHERE s.name = $sub AND d.name = $name MERGE (s) <- [:HAS] - (d)", map[string]interface{}{
				"name": p.NazwaProduktu,
				"sub":  strings.Trim(s, " "),
			})
			if err != nil {
				panic(err)
			}

		}
		count++
	}

}
