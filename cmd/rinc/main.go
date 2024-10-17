package main

import (
	"context"
	"log"
	"os"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/job"
	"github.com/accuknox/rinc/internal/kube"
	"github.com/accuknox/rinc/internal/web"
)

func main() {
	conf, err := conf.New(os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
	err = conf.Validate()
	if err != nil {
		log.Fatalf("validating provided config: %s", err.Error())
	}

	if conf.RunAsGenerator {
		kubeClient, err := kube.NewClient(conf.KubernetesClient)
		if err != nil {
			log.Fatalf("kubernetes client: %s", err.Error())
		}
		job := job.New(*conf, kubeClient)
		err = job.GenerateAll(context.Background())
		if err != nil {
			log.Fatalf("generating reports: %s", err.Error())
		}
		return
	}

	srv, err := web.NewSrv(*conf)
	if err != nil {
		log.Fatalf("creating web server instance: %s", err.Error())
	}
	srv.Run(context.Background())
}
