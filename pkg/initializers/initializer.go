package main

type ProjectInitializer interface {
	Init()
}

type ProjectInitialize struct {
}

func (*ProjectInitialize) Initialize(projectType string) {

}
