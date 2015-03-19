package importer

import (
	"fmt"
	fdb "github.com/mlapshin/fhirterm/db"
	"log"
	"path"
)

const createTableStmt = `
CREATE TABLE loinc_loincs
(
  loinc_num character varying NOT NULL primary key,
  component character varying,
  property character varying,
  time_aspect character varying,
  system character varying,
  scale_type character varying,
  method_type character varying,
  class character varying,
  source character varying,
  date_last_changed integer,
  change_type character varying,
  comments text,
  status character varying,
  consumer_name character varying,
  molar_mass character varying,
  classtype integer,
  formula character varying,
  species character varying,
  example_answers text,
  acssym text,
  base_name character varying,
  naaccr_id character varying,
  code_table character varying,
  survey_quest_text text,
  survey_quest_src character varying,
  units_required character varying,
  submitted_units character varying,
  relatednames2 text,
  shortname character varying,
  order_obs character varying,
  cdisc_common_tests character varying,
  hl7_field_subfield_id character varying,
  external_copyright_notice text,
  example_units character varying,
  long_common_name character varying,
  hl7_v2_datatype character varying,
  hl7_v3_datatype character varying,
  curated_range_and_units text,
  document_section character varying,
  example_ucum_units character varying,
  example_si_ucum_units character varying,
  status_reason character varying,
  status_text text,
  change_reason_public text,
  common_test_rank integer,
  common_order_rank integer,
  common_si_test_rank integer,
  hl7_attachment_structure character varying
)
`

const insertRowStmt = `
INSERT INTO loinc_loincs

VALUES
(
  ?,?,?,?,?,?,?,?,?,
  CAST(? AS integer),
  ?,?,?,?,?,
  CAST(? AS integer),
  ?,?,?,?,?,?,?,?,?,?,?,?,?,?,
  ?,?,?,?,?,?,?,?,?,?,?,?,?,?,
  CAST(? AS integer),
  CAST(? AS integer),
  CAST(? AS integer),
  ?
)`

const loincColumnsCount = 48

func createLoincTable(db *fdb.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS loinc_loincs")
	if err != nil {
		return err
	}

	_, err = db.Exec(createTableStmt)
	if err != nil {
		return err
	}

	log.Print("Created loinc_loincs table")
	return nil
}

func importLoincCsv(db *fdb.DB, csvPath string) error {
	insertedRows, err := importCsv(db, csvPath, ',', loincColumnsCount, insertRowStmt)

	if err != nil {
		return err
	}

	log.Printf("Done, imported %d LOINCs", insertedRows)
	return nil
}

func ImportLoinc(db *fdb.DB, filePath string) error {
	log.Printf("Importing LOINC dataset")

	err := unpackZipArchive(filePath, func(p string) error {
		err := createLoincTable(db)
		if err != nil {
			return err
		}

		err = importLoincCsv(db, path.Join(p, "loinc.csv"))
		return err
	})

	if err != nil {
		return fmt.Errorf("Error during importing LOINC: %s", err)
	} else {
		return nil
	}
}
