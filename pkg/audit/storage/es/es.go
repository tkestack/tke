package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
	"k8s.io/apimachinery/pkg/util/wait"
	"tkestack.io/tke/pkg/audit/apis/config"
	"tkestack.io/tke/pkg/audit/storage"
	"tkestack.io/tke/pkg/audit/storage/types"
	"tkestack.io/tke/pkg/util/log"
)

const typ = "tke-k8s-audit-event"
const batchSize = 100
const defaultIndices = "auditevent"
const defaultReverveDays = 7

var fieldEnumCache = map[string][]string{}
var lock sync.Mutex

func init() {
	fieldEnumCache["userName"] = []string{}
	fieldEnumCache["clusterName"] = []string{}
	fieldEnumCache["namespace"] = []string{}
	fieldEnumCache["resource"] = []string{}
}

type es struct {
	Addr        string
	Indices     string
	ReserveDays int
	username    string
	password    string
}

func NewStorage(conf *config.ElasticSearchStorage) (storage.AuditStorage, error) {
	cli := &es{
		Addr:        conf.Address,
		Indices:     conf.Indices,
		ReserveDays: conf.ReserveDays,
		username:    conf.Username,
		password:    conf.Password,
	}
	if cli.Indices == "" {
		cli.Indices = defaultIndices
	}
	if cli.ReserveDays < 0 {
		cli.ReserveDays = defaultReverveDays
	}
	err := cli.init()
	if err != nil {
		return nil, err
	}
	go wait.Forever(cli.cleanup, time.Hour)
	go wait.Forever(cli.updateFieldEnumCache, time.Minute)
	return cli, nil
}

func (s *es) init() error {
	if !s.indicesTypeExist() {
		return s.indicesTypeCreate()
	}
	return nil
}

func (s *es) indicesTypeExist() bool {
	req := gorequest.New()
	resp, _, err := req.Get(fmt.Sprintf("%s/%s/_mapping/%s", s.Addr, s.Indices, typ)).SetBasicAuth(s.username, s.password).End()
	if len(err) != 0 {
		return false
	} else if resp.StatusCode != 200 {
		return false
	}
	return true
}

func (s *es) indicesTypeCreate() error {
	keywords := []string{"stage", "verb", "userName", "resource", "namespace", "name", "status", "clusterName"}
	texts := []string{"auditID", "requestURI", "userAgent", "uid", "apiGroup", "apiVersion", "message", "reason", "details", "requestObject", "responseObject", "sourceIPs"}
	req := gorequest.New().Put(fmt.Sprintf("%s/%s", s.Addr, s.Indices)).SetBasicAuth(s.username, s.password)
	req.Header["content-type"] = "application/json"
	properties := map[string]map[string]string{
		"code": {
			"type": "integer",
		},
		"requestReceivedTimestamp": {
			"type": "long",
		},
		"stageTimestamp": {
			"type": "long",
		},
	}
	for _, keyword := range keywords {
		properties[keyword] = map[string]string{
			"type": "keyword",
		}
	}
	for _, text := range texts {
		properties[text] = map[string]string{
			"type": "text",
		}
	}
	reqBody := map[string]interface{}{
		"mappings": map[string]interface{}{
			typ: map[string]interface{}{
				"properties": properties,
			},
		},
	}

	resp, body, err := req.SendStruct(reqBody).End()
	if len(err) > 0 {
		return fmt.Errorf("indicesTypeCreate failed: %v", err)
	} else if resp.StatusCode >= 300 {
		return fmt.Errorf("indicesTypeCreate failed code %d, body %s", resp.StatusCode, body)
	}
	return nil
}

func (s *es) Query(param *storage.QueryParameter) ([]*types.Event, int, error) {
	if param == nil {
		param = &storage.QueryParameter{Size: 10}
	}
	var terms []interface{}

	if param.ClusterName != "" {
		terms = append(terms, map[string]map[string]string{"term": {"clusterName": param.ClusterName}})
	}
	if param.Namespace != "" {
		terms = append(terms, map[string]map[string]string{"term": {"namespace": param.Namespace}})
	}
	if param.Name != "" {
		terms = append(terms, map[string]map[string]string{"term": {"name": param.Name}})
	}
	if param.Resource != "" {
		terms = append(terms, map[string]map[string]string{"term": {"resource": param.Resource}})
	}
	if param.UserName != "" {
		terms = append(terms, map[string]map[string]string{"term": {"userName": param.UserName}})
	}
	if param.StartTime > 0 && param.EndTime > 0 {
		terms = append(terms, map[string]map[string]map[string]int64{"range": {"requestReceivedTimestamp": {"gte": param.StartTime, "lte": param.EndTime}}})
	} else if param.EndTime > 0 {
		terms = append(terms, map[string]map[string]map[string]int64{"range": {"requestReceivedTimestamp": {"lte": param.EndTime}}})
	} else if param.StartTime > 0 {
		terms = append(terms, map[string]map[string]map[string]int64{"range": {"requestReceivedTimestamp": {"gte": param.StartTime}}})
	}
	if param.Query != "" {
		terms = append(terms, map[string]map[string]interface{}{"multi_match": {
			"query":  param.Query,
			"fields": []string{"message", "details", "requestObject", "responseObject"},
		}})
	}
	query := map[string]interface{}{
		"from": param.Offset,
		"size": param.Size,
		"sort": []interface{}{
			map[string]string{"requestReceivedTimestamp": "desc"},
		},
	}
	if len(terms) > 0 {
		query["query"] = map[string]map[string]interface{}{
			"bool": {"filter": terms},
		}
	}
	req := gorequest.New().Get(fmt.Sprintf("%s/%s/%s/_search", s.Addr, s.Indices, typ)).SetBasicAuth(s.username, s.password)
	req.Header["content-type"] = "application/json"
	resp, body, errs := req.SendStruct(query).End()
	if len(errs) > 0 {
		return nil, 0, fmt.Errorf("failed search documents: %v", errs)
	} else if resp.StatusCode >= 300 {
		return nil, 0, fmt.Errorf("failed search document: %s", body)
	}
	res := Result{}
	err := json.Unmarshal([]byte(body), &res)
	if err != nil {
		return nil, 0, err
	}
	events := make([]*types.Event, 0)
	for _, ev := range res.Hits.Hits {
		events = append(events, ev.Event)
	}
	return events, res.Hits.Total, nil
}

type Result struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	Total int         `json:"total"`
	Hits  []*Document `json:"hits"`
}

type Document struct {
	Event *types.Event `json:"_source"`
}

func (s *es) Save(events []*types.Event) error {
	var batchEvents []*types.Event
	for _, event := range events {
		batchEvents = append(batchEvents, event)
		if len(batchEvents) >= batchSize {
			if err := s.batchSave(batchEvents); err != nil {
				return err
			}
			batchEvents = nil
		}
	}
	if len(batchEvents) > 0 {
		return s.batchSave(batchEvents)
	}
	return nil
}

func (s *es) FieldValues() map[string][]string {
	lock.Lock()
	defer lock.Unlock()
	result := make(map[string][]string)
	for field, values := range fieldEnumCache {
		result[field] = values
	}
	return result
}

func (s *es) batchSave(events []*types.Event) error {
	req := gorequest.New().Post(fmt.Sprintf("%s/%s/%s/_bulk", s.Addr, s.Indices, typ)).SetBasicAuth(s.username, s.password)
	req.Header["content-type"] = "application/x-ndjson"
	req.BounceToRawString = true
	buf := bytes.NewBuffer(nil)
	createEvent := `{"index":{}}`
	for _, event := range events {
		s, _ := json.Marshal(event)
		buf.WriteString(createEvent + "\n")
		buf.WriteString(string(s) + "\n")
	}
	resp, body, err := req.SendString(buf.String()).End()
	if len(err) > 0 {
		return fmt.Errorf("bulk index document failed: %v", err)
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("bulk index document failed: %s", body)
	}
	return nil
}

func (s *es) cleanup() {
	log.Infof("trigger es audit event cleanup")
	req := gorequest.New().Post(fmt.Sprintf("%s/%s/%s/_delete_by_query", s.Addr, s.Indices, typ)).SetBasicAuth(s.username, s.password)
	req.Header["content-type"] = "application/json"
	t := time.Now().Unix()*1000 - int64(s.ReserveDays*24*60*60*1000)
	query := fmt.Sprintf(`{"query":{"bool":{"filter":{"range":{"requestReceivedTimestamp":{"lte":%d}}}}}}`, t)
	_, _, errs := req.SendString(query).End()
	if len(errs) != 0 {
		log.Errorf("failed cleanup older audit events: %v", errs)
	}
}

func (s *es) updateFieldEnumCache() {
	var l sync.Mutex
	var tmpMap = map[string][]string{}
	wg := sync.WaitGroup{}
	for field := range fieldEnumCache {
		wg.Add(1)
		go func(field string) {
			defer wg.Done()
			req := gorequest.New().Get(fmt.Sprintf("%s/%s/%s/_search", s.Addr, s.Indices, typ)).SetBasicAuth(s.username, s.password)
			req.Header["content-type"] = "application/json"
			_, body, errs := req.SendString(fmt.Sprintf(`{"size":0,"aggs":{"distinct_colors":{"terms":{"field":"%s","size":1000}}}}`, field)).End()
			result := struct {
				Aggregations struct {
					DistinctColors struct {
						Buckets []struct {
							Key string `json:"key"`
						} `json:"buckets"`
					} `json:"distinct_colors"`
				} `json:"aggregations"`
			}{}
			if len(errs) > 0 {
				log.Errorf("failed update field %s values: %v", field, errs)
			}
			if err := json.Unmarshal([]byte(body), &result); err != nil {
				log.Errorf("can't get field %s values: %v", field, err)
			}
			var values []string
			for _, bucket := range result.Aggregations.DistinctColors.Buckets {
				if bucket.Key != "" {
					values = append(values, bucket.Key)
				}
			}
			if len(values) > 0 {
				l.Lock()
				tmpMap[field] = values
				l.Unlock()
			} else {
				fmt.Printf("filed %s return %v\n", field, body)
			}
		}(field)
	}
	wg.Wait()
	lock.Lock()
	defer lock.Unlock()
	for field, values := range tmpMap {
		fieldEnumCache[field] = values
	}
}
