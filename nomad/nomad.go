package nomad

import (
	"encoding/json"

	"github.com/hashicorp/nomad/api"
)

type NomadClient struct {
	*api.Client
}

func New() (NomadClient, error) {
	client, err := api.NewClient(api.DefaultConfig())
	return NomadClient{client}, err
}

func (client NomadClient) ListNs() ([]string, error) {
	var list []string = []string{}
	namespaces, _, err := client.Namespaces().List(&api.QueryOptions{})
	if err != nil {
		return list, err
	}
	for _, ns := range namespaces {
		list = append(list, ns.Name)
	}
	return list, nil
}

func (client NomadClient) ListJobs(ns string) ([]string, error) {
	var list []string = []string{}
	jobs, _, err := client.Jobs().List(&api.QueryOptions{Namespace: ns})

	if err != nil {
		return list, err
	}
	for _, job := range jobs {
		list = append(list, job.ID)
	}
	return list, nil
}

func (client NomadClient) InspectJob(ns string, job string) (string, *api.Job, error) {
	info, _, err := client.Jobs().Info(job, &api.QueryOptions{Namespace: ns})
	out, _ := json.Marshal(&info)
	return string(out), info, err
}
