package importer

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var createTblStmts map[string]string

var createIndexStmts = []string{
	"CREATE INDEX snomed_is_a_relationships_on_source_id_idx ON snomed_is_a_relationships(source_id)",
	"CREATE INDEX snomed_is_a_relationships_on_destination_id_idx ON snomed_is_a_relationships(destination_id)",
}

const fillConceptsNoHistoryStmt = `
INSERT INTO snomed_concepts_no_history
(concept_id, effective_time, term)
SELECT concept_id, max(effective_time), term
FROM snomed_descriptions GROUP BY concept_id`

const fillIsARelationsipsStmt = `
INSERT INTO snomed_is_a_relationships
(id, source_id, destination_id)
SELECT id, source_id, destination_id FROM snomed_relationships
WHERE type_id = 116680003 AND active = 1`

func init() {
	createTblStmts = make(map[string]string)
	createTblStmts["snomed_concepts"] = `
CREATE TABLE snomed_concepts
(
  id integer,
  effective_time integer,
  active boolean,
  module_id integer,
  definition_status_id integer
)
`

	createTblStmts["snomed_relationships"] = `
CREATE TABLE snomed_relationships
(
  id integer,
  effective_time integer,
  active boolean,
  module_id integer,
  source_id integer,
  destination_id integer,
  relationship_group integer,
  type_id integer,
  characteristic_type_id integer,
  modifier_id integer
)`

	createTblStmts["snomed_descriptions"] = `
CREATE TABLE snomed_descriptions
(
  id integer,
  effective_time integer,
  active integer,
  module_id integer,
  concept_id integer,
  language_code text,
  type_id integer,
  term text,
  case_significance_id integer
)`

	createTblStmts["snomed_is_a_relationships"] = `
CREATE TABLE snomed_is_a_relationships
(
  id bigint,
  source_id bigint,
  destination_id bigint
)`

	createTblStmts["snomed_concepts_no_history"] = `
CREATE TABLE snomed_concepts_no_history
(
  concept_id bigint NOT NULL PRIMARY KEY,
  effective_time integer,
  term text
)`

	createTblStmts["snomed_ancestors_descendants"] = `
CREATE TABLE snomed_ancestors_descendants
(
  concept_id integer primary key,
  ancestors blob,
  descendants blob
)`
}

const insertConceptsStmt = `
INSERT INTO snomed_concepts
VALUES (
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer)
)`

const insertRelsStmt = `
INSERT INTO snomed_relationships
VALUES (
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer)
)`

const insertDescStmt = `
INSERT INTO snomed_descriptions
VALUES (
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
CAST(? AS integer),
?,
CAST(? AS integer),
?,
CAST(? AS integer)
)`

func dirContent(root string) ([]string, error) {
	result := make([]string, 100)

	walker := func(path string, f os.FileInfo, err error) error {
		result = append(result, path)
		return nil
	}

	err := filepath.Walk(root, walker)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func findFile(files []string, rexp string) (string, bool) {
	r, _ := regexp.Compile(rexp)

	for _, file := range files {
		if r.MatchString(file) {
			return file, true
		}
	}

	return "", false
}

func createSnomedTables(db *sql.DB) error {
	for tblName, stmt := range createTblStmts {
		_, err := db.Exec("DROP TABLE IF EXISTS " + tblName)
		if err != nil {
			return err
		}

		_, err = db.Exec(stmt)
		if err != nil {
			return err
		}

		log.Printf("Created %s table", tblName)
	}

	return nil
}

func importSnomedConcepts(db *sql.DB, files []string) error {
	csvPath, found := findFile(files, "SnomedCT_Release_INT_\\d{8}/RF2Release/Full/Terminology/sct2_Concept_Full_INT_\\d{8}.txt$")
	log.Printf("Importing %s", csvPath)

	if !found {
		return fmt.Errorf("Could not find file sct2_Concept_Full_INT_XXXXXXXX.txt in SNOMED archive")
	}

	importedRows, err := importCsv(db, csvPath, '\t', 5, insertConceptsStmt)

	if err != nil {
		return err
	}

	log.Printf("Imported %d rows into snomed_concepts table", importedRows)

	return nil
}

func importSnomedRelationships(db *sql.DB, files []string) error {
	csvPath, found := findFile(files, "SnomedCT_Release_INT_\\d{8}/RF2Release/Full/Terminology/sct2_Relationship_Full_INT_\\d{8}.txt$")
	log.Printf("Importing %s", csvPath)

	if !found {
		return fmt.Errorf("Could not find file sct2_Relationship_Full_INT_XXXXXXXX.txt in SNOMED archive")
	}

	importedRows, err := importCsv(db, csvPath, '\t', 10, insertRelsStmt)

	if err != nil {
		return err
	}

	log.Printf("Imported %d rows into snomed_relationships table", importedRows)

	return nil
}

func escapeQuotes(f string) error {
	fixedCsv, _ := os.Create(f + "-fixed")
	csvFile, _ := os.Open(f)

	scanner := bufio.NewScanner(csvFile)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.ContainsRune(line, '"') {
			fields := strings.Split(line, "\t")
			fields[7] = "\"" + strings.Replace(fields[7], "\"", "\"\"", -1) + "\""
			line = strings.Join(fields, "\t")
		}

		fixedCsv.Write([]byte(line + "\n"))
	}

	fixedCsv.Close()
	csvFile.Close()

	return nil
}

func importSnomedDescriptions(db *sql.DB, files []string) error {
	csvPath, found := findFile(files, "SnomedCT_Release_INT_\\d{8}/RF2Release/Full/Terminology/sct2_Description_Full-en_INT_\\d{8}.txt$")
	log.Printf("Importing %s", csvPath)

	if !found {
		return fmt.Errorf("Could not find file sct2_Description_Full_INT_XXXXXXXX.txt in SNOMED archive")
	}

	// Escape quotes in TSV file
	// Otherwise, encounding/csv will fail to corretly load this file
	escapeQuotes(csvPath)

	importedRows, err := importCsv(db, csvPath+"-fixed", '\t', 9, insertDescStmt)

	if err != nil {
		return err
	}

	log.Printf("Imported %d rows into snomed_descriptions table", importedRows)

	return nil
}

func execStmt(db *sql.DB, stmt string, logMessage string) error {
	if len(logMessage) > 0 {
		log.Print(logMessage)
	}

	_, err := db.Exec(stmt)

	if err != nil {
		return err
	}

	log.Print("Done")

	return nil
}

func rowsToIntSlice(rows *sql.Rows, slice []int64) ([]int64, error) {
	for rows.Next() {
		var i int64
		err := rows.Scan(&i)
		if err != nil {
			return nil, err
		}

		slice = append(slice, i)
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}

	rows.Close()

	return slice, nil
}

func getAncestorsDescendants(stmt *sql.Stmt, slice *[]int64) {

}

func prewalkSnomedGraph(db *sql.DB) error {
	log.Print("Prewalking SNOMED-CT graph...")

	rows, err := db.Query(`SELECT source_id FROM snomed_is_a_relationships
														UNION
														SELECT destination_id FROM snomed_is_a_relationships`)
	if err != nil {
		return err
	}

	concepts := make([]int64, 0, 355000)
	concepts, _ = rowsToIntSlice(rows, concepts)

	log.Printf("Collected %d concepts", len(concepts))

	var insertTx *sql.Tx
	var insertStmt *sql.Stmt
	var getAncestorsStmt *sql.Stmt
	var getDescendantsStmt *sql.Stmt

	ancestors := make([]int64, 0, 200000)
	descendants := make([]int64, 0, 200000)

	ancestorsBuf := new(bytes.Buffer)
	descendantsBuf := new(bytes.Buffer)

	for index, concept := range concepts {
		if index%20000 == 0 {
			if insertTx != nil {
				insertTx.Commit()
			}

			insertTx, _ = db.Begin()

			insertStmt, err = insertTx.Prepare(`INSERT INTO snomed_ancestors_descendants
                                          (concept_id, ancestors, descendants)
                                          VALUES (?, ?, ?)`)
			if err != nil {
				return err
			}

			getAncestorsStmt, _ =
				insertTx.Prepare(`WITH RECURSIVE t(destination_id) AS (
													SELECT destination_id FROM snomed_is_a_relationships
													WHERE source_id = ?
													UNION
													SELECT sr.destination_id FROM snomed_is_a_relationships AS sr
													JOIN t ON t.destination_id = sr.source_id
													) SELECT destination_id FROM t`)

			getDescendantsStmt, _ =
				insertTx.Prepare(`WITH RECURSIVE t(source_id) AS (
													SELECT source_id FROM snomed_is_a_relationships
													WHERE destination_id = ?
													UNION
													SELECT sr.source_id FROM snomed_is_a_relationships AS sr
													JOIN t ON t.source_id = sr.destination_id
													) SELECT source_id FROM t`)
		}

		if index%20000 == 0 && index != 0 {
			log.Printf("Processed %d concepts...", index)
		}

		rows, err := getAncestorsStmt.Query(concept)
		if err != nil {
			return err
		}
		ancestors = ancestors[0:0]
		ancestors, _ = rowsToIntSlice(rows, ancestors)

		rows, err = getDescendantsStmt.Query(concept)
		if err != nil {
			return err
		}
		descendants = descendants[0:0]
		descendants, _ = rowsToIntSlice(rows, descendants)

		ancestorsBuf.Reset()
		descendantsBuf.Reset()
		err = binary.Write(ancestorsBuf, binary.LittleEndian, ancestors)
		if err != nil {
			return err
		}

		err = binary.Write(descendantsBuf, binary.LittleEndian, descendants)
		if err != nil {
			return err
		}

		_, err = insertStmt.Exec(concept, ancestorsBuf.Bytes(), descendantsBuf.Bytes())
		if err != nil {
			return err
		}
	}
	// final commit
	insertTx.Commit()
	log.Print("Done")

	return nil
}

func ImportSnomed(db *sql.DB, filePath string) error {
	log.Printf("Importing SNOMED-CT dataset")

	error := unpackZipArchive(filePath, func(p string) error {
		files, err := dirContent(p)

		err = createSnomedTables(db)
		if err != nil {
			return err
		}

		err = importSnomedConcepts(db, files)
		if err != nil {
			return err
		}

		err = importSnomedRelationships(db, files)
		if err != nil {
			return err
		}

		err = importSnomedDescriptions(db, files)
		if err != nil {
			return err
		}

		err = execStmt(
			db,
			fillIsARelationsipsStmt,
			"Filling snomed_is_a_relationships table")

		if err != nil {
			return err
		}

		err = execStmt(
			db,
			fillConceptsNoHistoryStmt,
			"Filling snomed_concepts_no_history table")

		if err != nil {
			return err
		}

		log.Print("Creating indices")
		for _, s := range createIndexStmts {
			err = execStmt(db, s, "")

			if err != nil {
				return err
			}
		}

		err = prewalkSnomedGraph(db)
		if err != nil {
			return err
		}

		return nil
	})

	if error != nil {
		return fmt.Errorf("Error during importing SNOMED-CT: %s", error)
	} else {
		return nil
	}
}
