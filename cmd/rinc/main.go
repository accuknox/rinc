package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/accuknox/rinc/internal/conf"
	"github.com/accuknox/rinc/internal/db"
	"github.com/accuknox/rinc/internal/job"
	"github.com/accuknox/rinc/internal/kube"
	"github.com/accuknox/rinc/internal/schema"
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

	if conf.GenerateSchema != "" {
		schema, err := schema.Generate(conf.GenerateSchema)
		if err != nil {
			log.Fatalf("generating schema: %s", err.Error())
		}
		fmt.Println(string(schema))
		return
	}

	mongo, err := db.NewMongoDBClient(conf.Mongodb)
	if err != nil {
		log.Fatalf("creating mongo client: %s", err.Error())
	}
	defer func() {
		ctx := context.TODO()
		err := mongo.Disconnect(ctx)
		if err != nil {
			slog.LogAttrs(
				ctx,
				slog.LevelError,
				"closing mongodb client connection",
				slog.String("error", err.Error()),
			)
		}
		slog.LogAttrs(
			ctx,
			slog.LevelDebug,
			"closed mongodb client connection",
		)
	}()

	if conf.RunAsScraper {
		kubeClient, err := kube.NewClient(conf.KubernetesClient)
		if err != nil {
			log.Fatalf("kubernetes client: %s", err.Error())
		}
		job := job.New(*conf, kubeClient, mongo)
		err = job.GenerateAll(context.Background())
		if err != nil {
			log.Fatalf("generating reports: %s", err.Error())
		}
		return
	}

	srv, err := web.NewSrv(*conf, mongo)
	if err != nil {
		log.Fatalf("creating web server instance: %s", err.Error())
	}
	srv.Run(context.Background())
}
