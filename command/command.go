package command

import "github.com/jcarley/harbor/models"

type Command interface {
	Run(job models.Job) int
}
