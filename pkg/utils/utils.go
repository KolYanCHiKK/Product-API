package utils

import "github.com/lib/pq"

func MapProductImage(image string) pq.StringArray {
	if image == "" {
		return nil
	}
	return pq.StringArray{image}
}
