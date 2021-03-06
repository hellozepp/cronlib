package main

// 测试crontab spec的格式

import (
	"log"
	"time"

	"github.com/rfyiamcool/cronlib"
)

//启动多个任务
func main() {
	cron := cronlib.New()

	specList := map[string]string{
		"risk.scan.total.1s":       "*/1 * * * * *",
		"risk.scan.total.2s":       "*/2 * * * * *",
		"risk.scan.total.3s":       "*/3 * * * * *",
		"risk.scan.total.4s":       "*/4 * * * * *",
		"risk.scan.total.5s.to.3s": "*/5 * * * * *",
	}

	for srv, spec := range specList {
		tspec := spec // copy
		ssrv := srv   // copy
		job, err := cronlib.NewJobModel(
			spec,
			func() {
				stdout(ssrv, tspec)
			},
		)
		if err != nil {
			panic(err.Error())
		}

		err = cron.Register(srv, job)
		if err != nil {
			panic(err.Error())
		}
	}

	// update test
	time.AfterFunc(10*time.Second, func() {
		spec := "*/3 * * * * *"
		srv := "risk.scan.total.5s.to.3s"
		log.Println("reset 5s to 3s", srv)
		job, _ := cronlib.NewJobModel(
			spec,
			func() {
				stdout(srv, spec)
			},
		)
		cron.UpdateJobModel(srv, job)
		log.Println("reset finish", srv)

	})

	// kill test
	time.AfterFunc(3*time.Second, func() {

		srv := "risk.scan.total.1s"
		log.Println("stoping", srv)
		cron.StopService(srv)
		log.Println("stop finish", srv)

	})

	time.AfterFunc(11*time.Second, func() {

		srvPrefix := "risk"
		log.Println("stoping srv prefix", srvPrefix)
		cron.StopServicePrefix(srvPrefix)

	})

	cron.Start()
	cron.Join()
}

func stdout(srv, spec string) {
	log.Println(srv, spec)
}

