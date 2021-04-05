package user

import "github.com/scottcagno/go-blog/pkg/web/templates"

type UserService struct {
	t *templates.TemplateCache
}

func NewUserService(t *templates.TemplateCache) *UserService {
	return &UserService{
		t: t,
	}
}
