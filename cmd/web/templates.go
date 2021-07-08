package main

import (
	"github.com/navruz-rakhimov/snippetbox/pkg/models"
)

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}
