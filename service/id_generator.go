package service

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

func GenerateId() (string, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Generate a snowflake ID.
	id := node.Generate()

	return id.String(), nil
}
