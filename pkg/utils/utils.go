package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func MapProductImage(image string) pq.StringArray {
	if image == "" {
		return nil
	}
	return pq.StringArray{image}
}

func BuildOrderBy(query string, allowedMap map[string]string) (string, error) {
	if query == "" {
		return "", nil
	}

	fields := strings.Split(query, ",")
	orderMap := make([]string, len(fields))
	for index, value := range fields {
		isDesc := strings.HasPrefix(value, "-")
		direction := "ASC"
		if isDesc {
			value = strings.TrimPrefix(value, "-")
			direction = "DESC"
		}

		column, ok := allowedMap[value]
		if !ok {
			return "", errors.New("Invalid parameters in orderBy")
		}

		orderMap[index] = column + " " + direction
	}
	joinBeforeLast := strings.Join(orderMap, ", ") + ","

	return joinBeforeLast, nil
}

func BuildWhereFilter(parameters map[string]any, whiteList map[string]string) (*string, []any, error) {
	var countNonEmptyKey int
	for key, value := range parameters {
		if key == "" {
			countNonEmptyKey++
			delete(parameters, key)
		}
		switch t := value.(type) {
		case string:
			if t == "" {
				delete(parameters, key)
			}
		case *int:
			if t == nil {
				delete(parameters, key)
			}
		case *float64:
			if t == nil {
				delete(parameters, key)
			}
		}
	}
	if countNonEmptyKey == len(parameters) {
		return nil, nil, errors.New("No parameters for where")
	}

	whereDefinition := make([]string, 0, len(parameters))
	whereParameters := make([]any, 0, len(parameters))
	for key, value := range parameters {
		column, ok := whiteList[key]
		if ok != true {
			return nil, nil, errors.New("Invalid parameters in where")
		}
		whereDefinition = append(whereDefinition, fmt.Sprintf("%s = $%v", column, len(whereDefinition)+1))
		whereParameters = append(whereParameters, value)
	}
	joinWhereRow := strings.Join(whereDefinition, " AND ")

	return &joinWhereRow, whereParameters, nil
}

func BuildSetParams(parameters map[string]any, whiteList map[string]string) (*string, []any, error) {
	var countNonEmptyKey int
	for key, value := range parameters {
		if key == "" {
			countNonEmptyKey++
			delete(parameters, key)
		}
		switch t := value.(type) {
		case string:
			if t == "" {
				delete(parameters, key)
			}
		case *int:
			if t == nil {
				delete(parameters, key)
			}
		case *float64:
			if t == nil {
				delete(parameters, key)
			}
		}
	}
	if countNonEmptyKey == len(parameters) {
		return nil, nil, errors.New("No parameters for set")
	}

	setDefinition := make([]string, 0, len(parameters))
	setParameters := make([]any, 0, len(parameters))
	for key, value := range parameters {
		column, ok := whiteList[key]
		if ok != true {
			return nil, nil, errors.New("Invalid parameters in set")
		}
		setDefinition = append(setDefinition, fmt.Sprintf("%s = $%v", column, len(setDefinition)+1))
		setParameters = append(setParameters, value)
	}
	joinWhereRow := strings.Join(setDefinition, ", ")

	return &joinWhereRow, setParameters, nil
}
