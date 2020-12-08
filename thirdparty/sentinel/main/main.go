// Copyright 1999-2020 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/aliyun/aliyun-ahas-go-sdk/config"
	"log"
	"math/rand"
	"os"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/util"
	ahas "github.com/aliyun/aliyun-ahas-go-sdk"
)

func main() {
	//conf := config.NewDefaultConfig()
	// for testing, logging output to console
	//conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	//err := sentinel.InitWithConfig(conf)
	//if err != nil {
	//	log.Fatal(err)
	//}

	os.Setenv(config.LicenseEnvKey, "7bac187365984241afadce74133c9820")
	os.Setenv(config.EnvironmentEnvKey, config.DeployEnvProd)

	err := ahas.InitAhasDefault()
	if err != nil {
		log.Fatalf("Failed to init AHAS: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	ch := make(chan struct{})

	for i := 0; i < 10; i++ {
		go func() {
			for {
				e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					// Blocked. We could get the block reason from the BlockError.
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					// Passed, wrap the logic here.
					fmt.Println(util.CurrentTimeMillis(), "passed")
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)

					// Be sure the entry is exited finally.
					e.Exit()
				}

			}
		}()
	}
	<-ch
}
