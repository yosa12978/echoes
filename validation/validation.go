package validation

import "context"

type Validatable interface {
	Validate(ctx context.Context) (problems map[string]string, ok bool)
}
