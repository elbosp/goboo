package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobeam/stringy"
	"github.com/otiai10/copy"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

type GobooConfig struct {
	Config  Config  `yaml:"configuration"`
	History History `yaml:"only_goboo_can_write"`
}

type Config struct {
	ServiceName string   `yaml:"service_name"`
	Domains     []string `yaml:"domains"`
	Aggregate   []string `yaml:"aggregate"`
}

type History struct {
	ServiceName string `yaml:"service_name"`
}

func main() {
	app := cli.NewApp()
	app.Name = "goboo"
	app.Usage = "An artisan framework. Mostly, doing magic to reduce your effort and working time"

	app.Commands = []cli.Command{
		{
			Name:  "config",
			Usage: "Configuring your service by reading config-goboo.yaml",
			Action: func(c *cli.Context) (err error) {
				var gobooConfig GobooConfig

				yamlFile, err := os.Open("config-goboo.yaml")

				if err != nil {
					log.Fatal(err)
				}

				defer yamlFile.Close()

				byteValue, err := ioutil.ReadAll(yamlFile)

				if err != nil {
					log.Fatal(err)
				}

				yaml.Unmarshal(byteValue, &gobooConfig)

				if gobooConfig.Config.ServiceName == "svc-boilerplate-golang" {
					log.Fatal("Ubah configuration.service_name terlebih dahulu!")
				}

				err = filepath.Walk(".", getWalkFunc("*.go", gobooConfig.History.ServiceName, gobooConfig.Config.ServiceName))

				if err != nil {
					log.Fatal(err)
				}

				err = filepath.Walk(".", getWalkFunc("go.mod", gobooConfig.History.ServiceName, gobooConfig.Config.ServiceName))

				if err != nil {
					log.Fatal(err)
				}

				gobooConfig.History.ServiceName = gobooConfig.Config.ServiceName

				newByteValue, err := yaml.Marshal(&gobooConfig)

				if err != nil {
					log.Fatal(err)
				}

				err = ioutil.WriteFile("config-goboo.yaml", newByteValue, 0755)

				if err != nil {
					log.Fatal(err)
				}

				for _, x := range gobooConfig.Config.Domains {
					xWithNoSpace := strings.Replace(x, " ", "", -1)
					xKebabCase := stringy.New(x).KebabCase().Get()
					xSnakeCase := stringy.New(x).SnakeCase().Get()
					xPascalCase := stringy.New(x).CamelCase()
					xCamelCase := stringy.New(xPascalCase).LcFirst()
					xMoreThanOneWord := strings.Contains(x, " ")

					if x == "" {
						continue
					}

					if _, err := os.Stat("domain/" + xKebabCase); !os.IsNotExist(err) {
						continue
					}

					err = copy.Copy("domain/boilerplate", "domain/"+xKebabCase)

					if err != nil {
						log.Fatal(err)
					}

					err = filepath.Walk("domain/"+xKebabCase, getWalkFunc("*.go", "package boilerplate", "package "+xWithNoSpace))

					if err != nil {
						log.Fatal(err)
					}

					if xMoreThanOneWord {
						err = filepath.Walk("domain/"+xKebabCase, getWalkFunc(
							"*.go",
							`"`+gobooConfig.Config.ServiceName+`/domain/boilerplate"`,
							xCamelCase+` "`+gobooConfig.Config.ServiceName+"/domain/"+xKebabCase+`"`,
						))

						if err != nil {
							log.Fatal(err)
						}
					}

					err = filepath.Walk("domain/"+xKebabCase+"/delivery/http/", getWalkFunc("handler.go", `"/boilerplate`, `"/`+xKebabCase))

					if err != nil {
						log.Fatal(err)
					}

					err = filepath.Walk("domain/"+xKebabCase+"/repository/", getWalkFunc("mysql.go", "FROM boilerplate", "FROM "+xSnakeCase))

					if err != nil {
						log.Fatal(err)
					}

					err = filepath.Walk("domain/"+xKebabCase+"/repository/", getWalkFunc("mysql.go", "UPDATE boilerplate", "UPDATE "+xSnakeCase))

					if err != nil {
						log.Fatal(err)
					}

					err = filepath.Walk("domain/"+xKebabCase+"/repository/", getWalkFunc("mysql.go", "INTO boilerplate", "INTO "+xSnakeCase))

					if err != nil {
						log.Fatal(err)
					}

					err = filepath.Walk("domain/"+xKebabCase, getWalkFunc("*.go", "Boilerplate", xPascalCase))

					if err != nil {
						log.Fatal(err)
					}

					err = filepath.Walk("domain/"+xKebabCase, getWalkFunc("*.go", "boilerplate", xCamelCase))

					if err != nil {
						log.Fatal(err)
					}

					if _, err := os.Stat("models/" + xSnakeCase + ".go"); err == nil {
						continue
					}

					input, err := ioutil.ReadFile("models/boilerplate.go")

					if err != nil {
						log.Fatal(err)
					}

					err = ioutil.WriteFile("models/"+xSnakeCase+".go", input, 0755)

					if err != nil {
						log.Fatal(err)
					}

					err = filepath.Walk("models", getWalkFunc(xSnakeCase+".go", "Boilerplate", xPascalCase))

					if err != nil {
						log.Fatal(err)
					}

					if _, err := os.Stat("entity/" + xSnakeCase + ".go"); err == nil {
						continue
					}

					input, err = ioutil.ReadFile("entity/boilerplate.go")

					if err != nil {
						log.Fatal(err)
					}

					err = ioutil.WriteFile("entity/"+xSnakeCase+".go", input, 0755)

					if err != nil {
						log.Fatal(err)
					}

					err = filepath.Walk("entity", getWalkFunc(xSnakeCase+".go", "Boilerplate", xPascalCase))

					if err != nil {
						log.Fatal(err)
					}
				}

				return
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getWalkFunc(file string, oldString string, newString string) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		matched, err := filepath.Match(file, fi.Name())

		if err != nil {
			fmt.Println(err)
			return err
		}

		if matched {
			read, err := ioutil.ReadFile(path)

			if err != nil {
				fmt.Println(err)
			}

			newContents := strings.Replace(string(read), oldString, newString, -1)

			err = ioutil.WriteFile(path, []byte(newContents), 0)

			if err != nil {
				fmt.Println(err)
				return err
			}
		}

		return nil
	}
}
