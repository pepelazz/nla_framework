package exampleProject

import (
	"github.com/pepelazz/projectGenerator/src/exampleProject/city"
	"github.com/pepelazz/projectGenerator/src/exampleProject/client"
	"github.com/pepelazz/projectGenerator/src/types"
)

func GetProject() types.ProjectType {
	p := &types.ProjectType{
		Name: "fourPl",
		Docs: []types.DocType {
			city.GetDoc(),
			client.GetDoc(),
		},
	}

	return *p
}
