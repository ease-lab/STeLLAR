package visualization

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	log "github.com/sirupsen/logrus"
	"io"
	"lambda-benchmarking/client/experiment/configuration"
	"os"
	"path/filepath"
	"time"
)

func GenerateVisualization(experiment configuration.SubExperiment, deltas []time.Duration, csvFile *os.File, path string) {
	log.Debugf("SubExperiment %d: reading written latencies file %s", experiment.Id, csvFile.Name())
	latenciesDF := readWrittenLatenciesFile(csvFile)

	switch experiment.Visualization {
	case "all":
		log.Infof("SubExperiment %d: generating all visualizations", experiment.Id)
		generateCDFs(experiment, latenciesDF, path)
		generateHistograms(experiment, latenciesDF, path, deltas)
		generateBarCharts(experiment, latenciesDF, path)
	case "bar":
		log.Infof("SubExperiment %d: generating burst bar chart visualization", experiment.Id)
		generateBarCharts(experiment, latenciesDF, path)
	case "CDF":
		log.Infof("SubExperiment %d: generating CDF visualization", experiment.Id)
		generateCDFs(experiment, latenciesDF, path)
	case "histogram":
		log.Infof("SubExperiment %d: generating histograms visualizations (per-burst)", experiment.Id)
		generateHistograms(experiment, latenciesDF, path, deltas)
	case "":
		log.Errorf("SubExperiment %d: no visualization selected, skipping", experiment.Id)
	default:
		log.Errorf("SubExperiment %d: unrecognized visualization `%s`, skipping", experiment.Id, experiment.Visualization)
	}
}

func generateBarCharts(experiment configuration.SubExperiment, latenciesDF dataframe.DataFrame, path string) {
	log.Debugf("SubExperiment %d: plotting characterization bar chart", experiment.Id)
	plotBurstsBarChart(filepath.Join(path, "bursts_characterization.png"), experiment, latenciesDF)
}

func generateHistograms(config configuration.SubExperiment, latenciesDF dataframe.DataFrame, path string, deltas []time.Duration) {
	histogramsDirectoryPath := filepath.Join(path, "histograms")
	log.Infof("Creating directory for histograms at `%s`", histogramsDirectoryPath)
	if err := os.MkdirAll(histogramsDirectoryPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	log.Debugf("SubExperiment %d: plotting latency histograms for each burst", config.Id)
	for burstIndex := 0; burstIndex < config.Bursts; burstIndex++ {
		burstDF := latenciesDF.Filter(dataframe.F{Colname: "Burst ID", Comparator: series.Eq, Comparando: burstIndex})
		plotBurstLatenciesHistogram(
			filepath.Join(histogramsDirectoryPath, fmt.Sprintf("burst%d_delta%v.png", burstIndex, deltas[burstIndex])),
			burstDF.Col("Client Latency (ms)").Float(),
			burstIndex,
			deltas[burstIndex],
		)
	}
}

func generateCDFs(config configuration.SubExperiment, latenciesDF dataframe.DataFrame, path string) {
	log.Debugf("SubExperiment %d: plotting latencies CDF", config.Id)
	plotLatenciesCDF(
		filepath.Join(path, "empirical_CDF.png"),
		latenciesDF.Col("Client Latency (ms)").Float(),
		config,
	)
}

func readWrittenLatenciesFile(csvFile *os.File) dataframe.DataFrame {
	_, err := csvFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Error(err)
	}

	latenciesDF := dataframe.ReadCSV(csvFile)
	return latenciesDF
}
