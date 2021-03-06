package event

import (
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `event`
}

func (p *Provider) Register() {
	p.Firmeve.Bind(`event`, New(), container.WithShare(true))
}

func (p *Provider) Boot() {

}
