// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package specjbb

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/intelsdi-x/swan/pkg/executor/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

// TestBackendWithMockedExecutor runs a Backend launcher with the mocked executor.
func TestBackendWithMockedExecutor(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	Convey("When using Backend launcher", t, func() {
		expectedCommand := "java -server -Xms2g -Xmx2g -XX:NativeMemoryTracking=summary -XX:+UseParallelOldGC  -XX:ParallelGCThreads=2 -XX:ConcGCThreads=1 -XX:InitiatingHeapOccupancyPercent=80 -XX:MaxGCPauseMillis=100 -XX:+AlwaysPreTouch  -Djava.library.path=/opt/swan/share/specjbb/lib -Dspecjbb.controller.host=127.0.0.1 -Dspecjbb.forkjoin.workers=1 -jar test -m backend -G GRP1 -J specjbbbackend1 -p /opt/swan/share/specjbb/config/specjbb2015.props"
		expectedHost := "127.0.0.1"

		mockedExecutor := new(mocks.Executor)
		mockedTaskHandle := new(mocks.TaskHandle)
		config := DefaultSPECjbbBackendConfig()
		config.PathToBinary = "test"
		backendLauncher := NewBackend(mockedExecutor, config)

		Convey("While simulating proper execution", func() {
			mockedExecutor.On("Execute", expectedCommand).Return(mockedTaskHandle, nil).Once()
			mockedTaskHandle.On("Address").Return(expectedHost)

			Convey("Build command should create proper command", func() {
				command := getBackendCommand(config)
				So(command, ShouldEqual, expectedCommand)

				Convey("Arguments passed to Executor should be a proper command", func() {
					task, err := backendLauncher.Launch()
					So(err, ShouldBeNil)

					So(task, ShouldNotBeNil)
					So(task, ShouldEqual, mockedTaskHandle)

					Convey("Location of the returned task shall be 127.0.0.1", func() {
						address := task.Address()
						So(address, ShouldEqual, expectedHost)
					})
				})
			})
		})

	})
}
