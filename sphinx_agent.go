package main

import (
	"github.com/yunge/sphinx"
        "fmt"
"time"
//	"github.com/yvasiyarov/newrelic_platform_go"
)

const (
    MIN_PAUSE_TIME = 30
)

type SphinxStatusData map[string]string
type MetricsDataSource struct {
	SphinxHost        string
	Port              int
	ConnectionTimeout int

        PreviousData SphinxStatusData
        LastData     SphinxStatusData
        LastUpdateTime time.Time
}

func NewMetricsDataSource(sphinxHost string, port int, connectionTimeout int) *MetricsDataSource {
	if port == 0 {
		port = 9312
	}
	ds := &MetricsDataSource{
		SphinxHost:        sphinxHost,
		Port:              port,
		ConnectionTimeout: connectionTimeout,
	}
	return ds
}

func (ds *MetricsDataSource) GetData() (*SphinxStatusData, error) {
        startTime := time.Now()
        if startTime.Sub(ds.LastUpdateTime) > time.Second * MIN_PAUSE_TIME {
            newData, err := ds.QueryData()
            if err != nil {
                return nil, err
            }
            
            if ds.PreviousData == nil {
                ds.PreviousData = newData
            } else {
                ds.PreviousData = ds.LastData
            }
            ds.LastData = newData
        }
        return nil, nil
}

func (ds *MetricsDataSource) QueryData() (*SphinxStatusData, error) {
	client := sphinx.NewClient().SetServer(ds.SphinxHost, ds.Port)

	if ds.ConnectionTimeout != 0 {
		client.SetConnectTimeout(ds.ConnectionTimeout)
	}
	if err := client.Error(); err != nil {
		return nil, err
	}

	defer client.Close()
	status, err := client.Status()
	if err != nil {
		return nil, err
	}

        data := make(SphinxStatusData, len(status))
	for _, row := range status {
		data[row[0]] = row[1]
	}

	return data, nil
}

/*
type Metrica struct {
	Name       string
	Units      string
	DataKey    string
	Datasource *MetricsDataSource
}

func (metrica *Metrica) GetName() string {
	return metrica.Name
}
func (metrica *WaveMetrica) GetUnits() string {
	return metrica.Units
}
func (metrica *WaveMetrica) GetValue() (float64, error) {
	metrica.sawtoothCounter++
	if metrica.sawtoothCounter > metrica.sawtoothMax {
		metrica.sawtoothCounter = 0
	}
	return float64(metrica.sawtoothCounter), nil
}
*/

func main() {
	/*
		plugin := newrelic_platform_go.NewNewrelicPlugin("0.0.1", "7bceac019c7dcafae1ef95be3e3a3ff8866de246", 60)
		component := newrelic_platform_go.NewPluginComponent("Sphinx component", "com.github.yvasiyarov.Sphinx")
		plugin.AddComponent(component)

		m := &WaveMetrica{
			sawtoothMax:     10,
			sawtoothCounter: 5,
		}

		component.AddMetrica(m)
		plugin.Verbose = true
		plugin.Run()
	*/
	ds := NewMetricsDataSource("web-d5.butik.ru", 0, 0)
	ds.QueryData()
}