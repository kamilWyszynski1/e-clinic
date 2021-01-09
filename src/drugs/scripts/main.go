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
	f, err := os.Open("/home/kamil/go/src/e-clinic/src/drugs/scripts/drugs.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	byteValue, _ := ioutil.ReadAll(f)
	var d drugs.ProduktyLecznicze
	if err := xml.Unmarshal(byteValue, &d); err != nil {
		panic(err)
	}
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
		log.Fatal(err)
	}
	defer session.Close()

	count := 0
	for _, p := range d.ProduktLeczniczy {
		if len(p.SubstancjeCzynne.SubstancjaCzynna) == 0 {
			continue
		}
		if count >= 10000 {
			return
		}
		if p.NazwaProduktu == "" || len(p.SubstancjeCzynne.SubstancjaCzynna) == 0 {
			continue
		}
		_, err := session.Run(`MERGE (l:Lek { 
				nazwa: $name, 
				id: $id, 
				typ: $type, 
				nazwa_powszechna: $common_name,
				moc: $strength,
				postac: $shape}) `,
			map[string]interface{}{
				"name":        p.NazwaProduktu,
				"id":          p.ID,
				"type":        p.RodzajPreparatu,
				"common_name": p.NPS,
				"strength":    p.Moc,
				"shape":       p.Postac,
			})
		if err != nil {
			panic(err)
		}

		for _, s := range p.SubstancjeCzynne.SubstancjaCzynna {
			_, err := session.Run("MERGE (s:Sub { nazwa: $sub })", map[string]interface{}{
				"sub": strings.Trim(s, " "),
			})
			if err != nil {
				panic(err)
			}

			_, err = session.Run(
				"MATCH (s:Sub), (l:Lek) WHERE s.nazwa = $sub AND l.nazwa = $name MERGE (s) <- [:POSIADA] - (l)",
				map[string]interface{}{
					"name": p.NazwaProduktu,
					"sub":  strings.Trim(s, " "),
				})
			if err != nil {
				log.Fatal(err)
			}

		}
		count++
	}
}
