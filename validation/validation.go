package validation

import (
	"context"
)

// type Problems map[string]string

// func (p Problems) Join(sep string) string {
// 	vals := make([]string, 0, len(p))
// 	for _, v := range p {
// 		vals = append(vals, v)
// 	}
// 	return strings.Join(vals, sep)
// }

type Validatable[T any] interface {
	Validate(ctx context.Context) (T, map[string]string, bool)
}
