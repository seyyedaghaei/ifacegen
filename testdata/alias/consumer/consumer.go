package consumer

import (
	shared1 "github.com/seyyedaghaei/ifacegen/testdata/alias/pkgone"
	shared2 "github.com/seyyedaghaei/ifacegen/testdata/alias/pkgtwo"
)

type AliasService struct{}

func (s *AliasService) Combine() (shared1.Thing, shared2.Thing) {
	return shared1.Thing{}, shared2.Thing{}
}
