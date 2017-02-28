package main

import (
	"fmt"
	"io"
	"os"
	//"strconv"
	//"strings"
	//"time"

	"github.com/Sirupsen/logrus"
	//"github.com/intelsdi-x/swan/experiments/specjbb-sensitivity-profile/common"
	"github.com/intelsdi-x/swan/pkg/conf"
	"github.com/intelsdi-x/swan/pkg/executor"
	"github.com/intelsdi-x/swan/pkg/experiment"
	//"github.com/intelsdi-x/swan/pkg/experiment/sensitivity"
	//"github.com/intelsdi-x/swan/pkg/experiment/sensitivity/topology"
	"github.com/intelsdi-x/swan/pkg/experiment/sensitivity/validate"
	//"github.com/intelsdi-x/swan/pkg/snap/sessions/specjbb"
	//"github.com/intelsdi-x/swan/pkg/utils/err_collection"
	"github.com/intelsdi-x/swan/pkg/workloads/specjbb"
	"github.com/nu7hatch/gouuid"
	//"github.com/pkg/errors"
	//"golang.org/x/tools/cmd/fiximports/testdata/src/titanic.biz/bar"
	//"gopkg.in/cheggaaa/pb.v1"
	"github.com/intelsdi-x/swan/pkg/kubernetes"
)

var (
	loadGeneratorOneAddress = conf.NewStringFlag(
		"specjbb_load_generator_one",
		"Address of the first SPECjbb Load Generator host",
		"127.0.0.1",
	)
	loadGeneratorTwoAddress = conf.NewStringFlag(
		"specjbb_load_generator_two",
		"Address of the second SPECjbb Load Generator host",
		"127.0.0.1",
	)
)

var (
	appName = os.Args[0]
)

func main() {
	err := conf.ParseFlags()
	if err != nil {
		logrus.Fatalf("Could not parse flags: %q", err.Error())
		os.Exit(experiment.ExSoftware)
	}

	logrus.SetLevel(logrus.DebugLevel)
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2000-01-02 15:04:05"
	formatter.FullTimestamp = true
	logrus.SetFormatter(formatter)

	// Generate an experiment ID and start the metadata session.
	uuid, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("Cannot generate experiment ID: %q", err.Error())
		os.Exit(experiment.ExSoftware)
	}

	logrus.Info("Starting Experiment with uuid ", uuid.String())

	//By default print only UUID of the experiment and nothing more on the stdout
	fmt.Println(uuid.String())

	// Each experiment should have it's own directory to store logs and errors
	experimentDirectory, logFile, err := experiment.CreateExperimentDir(uuid.String(), appName)
	if err != nil {
		logrus.Errorf("IO error: %q", err.Error())
		os.Exit(experiment.ExIOErr)
	}
	// TODO(skonefal): Use experiment directory,
	_ = experimentDirectory

	logrus.SetOutput(io.MultiWriter(logFile, os.Stderr))

	// Validate preconditions: for SPECjbb we only check if CPU governor is set to performance.
	validate.CheckCPUPowerGovernor()

	kubernetesExecutor := executor.NewLocal()
	kubernetesConfig := kubernetes.DefaultConfig()
	kubernetesLauncher := kubernetes.New(kubernetesExecutor, kubernetesExecutor, kubernetesConfig)
	kubernetesHandle, err := kubernetesLauncher.Launch()
	if err != nil {
		logrus.Errorf("could not prepare kubernetes cluster: %s", err)
		os.Exit(experiment.ExSoftware)
	}
	defer kubernetesHandle.Stop()

	specjbbBackendExecutorConfig := executor.DefaultKubernetesConfig()
	specjbbBackendExecutorConfig.PodNamePrefix = "specjbb-backend"
	specjbbBackendExecutorConfig.MemoryLimit = 10000000000
	specjbbBackendExecutorConfig.MemoryRequest = 10000000000
	specjbbBackendExecutorConfig.CPULimit = 8000
	specjbbBackendExecutorConfig.CPURequest = 8000
	specjbbBackendExecutorConfig.Privileged = true
	specjbbBackendExecutorConfig.HostNetwork = true
	specjbbBackendExecutor, err := executor.NewKubernetes(specjbbBackendExecutorConfig)
	if err != nil {
		logrus.Errorf("could not prepare specjbbBackendExecutor: %s", err)
		os.Exit(experiment.ExSoftware)
	}

	// Create launcher for high priority task (in case of SPECjbb it is a backend).
	backendConfig := specjbb.DefaultSPECjbbBackendConfig()
	backendConfig.ControllerAddress = specjbb.ControllerAddress.Value()
	backendConfig.JVMHeapMemoryGBs = 8
	backendConfig.WorkerCount = 8
	backendConfig.ParallelGCThreads = 4
	specjbbBackendLauncher := specjbb.NewBackend(specjbbBackendExecutor, backendConfig)

	// Prepare load generator for hp task (in case of the specjbb it is a controller with transaction injectors).
	txInjectorExecutorOne, err := executor.NewRemoteFromIP(loadGeneratorOneAddress.Value())
	if err != nil {
		logrus.Errorf("could not prepare txInjectorExecutorOne: %s", err)
		os.Exit(experiment.ExSoftware)
	}
	txInjectorExecutorTwo, err := executor.NewRemoteFromIP(loadGeneratorTwoAddress.Value())
	if err != nil {
		logrus.Errorf("could not prepare txInjectorExecutorTwo: %s", err)
		os.Exit(experiment.ExSoftware)
	}
	controllerExecutor, err := executor.NewRemoteFromIP(specjbb.ControllerAddress.Value())
	if err != nil {
		logrus.Errorf("could not prepare controllerExecutor: %s", err)
		os.Exit(experiment.ExSoftware)
	}

	// TODO(skonefal): Use two TxI.
	//loadGeneratorExecutors := []executor.Executor{txInjectorExecutorOne, txInjectorExecutorTwo}
	_ = txInjectorExecutorTwo
	loadGeneratorExecutors := []executor.Executor{txInjectorExecutorOne}
	loadGeneratorConfig := specjbb.DefaultLoadGeneratorConfig()
	loadGeneratorConfig.ControllerAddress = specjbb.ControllerAddress.Value()
	loadGeneratorConfig.JVMHeapMemoryGBs = 3
	specjbbLoadGenerator := specjbb.NewLoadGenerator(controllerExecutor, loadGeneratorExecutors, loadGeneratorConfig)

	// Run tuning.

	backend, err := specjbbBackendLauncher.Launch()
	if err != nil {
		logrus.Errorf("could not prepare specjbbBackendLauncher: %s", err)
		os.Exit(experiment.ExSoftware)
	}
	defer backend.Stop()

	qps, sli, err := specjbbLoadGenerator.Tune(10000)
	if err != nil {
		logrus.Errorf("could not prepare specjbbLoadGenerator: %s", err)
		os.Exit(experiment.ExSoftware)
	}

	logrus.Debugf("qps result: %d", qps)
	logrus.Debugf("load result: %d", sli)

}
