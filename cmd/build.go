package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/opub/scoreplus/model"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Generates application build artifacts.",
	Long:  "Generates SQL scripts that are used to create a new database schema for the application model storage.",
	Run: func(cmd *cobra.Command, args []string) {
		buildSchemaSQL()
	},
}

var models = []interface{}{model.Game{}, model.Member{}, model.Note{}, model.Sport{}, model.Team{}, model.Venue{}}

// generate SQL script to build database tables for the models listed above
func buildSchemaSQL() {
	for _, model := range models {
		buildTableSQL(reflect.TypeOf(model))
	}
}

func buildTableSQL(t reflect.Type) {
	table := strings.ToLower(t.Name())
	fmt.Printf("-- %s\n\nDROP TABLE IF EXISTS %s CASCADE;\nCREATE TABLE %s\n(\n", table, table, table)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.String() == "model.Base" {
			for j := 0; j < field.Type.NumField(); j++ {
				if field.Type.Field(j).Type.Name() != "Model" {
					buildColumnSQL(field.Type.Field(j), i == t.NumField()-1 && j == field.Type.NumField()-1)
				}
			}
		} else {
			buildColumnSQL(field, i == t.NumField()-1)
		}
	}

	fmt.Println(")\nWITH (\n\tOIDS = FALSE\n);")
	fmt.Printf("ALTER TABLE %s OWNER TO scoreplus_owner;\n", table)
	fmt.Printf("GRANT SELECT, USAGE ON SEQUENCE %s_id_seq TO scoreplus_writer;\nGRANT SELECT ON SEQUENCE %s_id_seq TO scoreplus_reader;\n\n", table, table)
}

func buildColumnSQL(f reflect.StructField, last bool) {
	name := strings.ToLower(f.Name)
	gotype := f.Type.Name()
	ctype := gotype

	switch gotype {
	case "int":
		ctype = "integer"
	case "int64":
		ctype = "integer"
	case "string":
		ctype = "text"
	case "bool":
		ctype = "boolean"
	case "Time":
		ctype = "timestamp with time zone"
	case "Int64Array":
		ctype = "integer[] DEFAULT '{}'"
	default:
		//cheating because nested structs have ints as FK
		ctype = "integer NOT NULL DEFAULT 0"
	}

	//couple special cases
	if name == "id" && gotype == "int64" {
		fmt.Printf("\t%s serial PRIMARY KEY", name)
	} else if name == "email" && gotype == "string" {
		fmt.Printf("\t%s email", name)
	} else {
		fmt.Printf("\t%s %s", name, ctype)
	}
	fmt.Printf("%v", f.Tag.Get("sql"))
	if !last {
		fmt.Print(",")
	}
	fmt.Println()
}
