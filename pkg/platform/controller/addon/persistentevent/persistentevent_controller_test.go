package persistentevent

import (
	"context"
	"gotest.tools/assert"
	"testing"
	v1 "tkestack.io/tke/api/platform/v1"
)

func TestMakeConfigMap(t *testing.T) {
	c := &Controller{}
	backend := &v1.PersistentBackEnd{
		ES: &v1.StorageBackEndES{
			IP:        "127.0.0.1",
			Port:      9200,
			Scheme:    "http",
			IndexName: "fluentd",
		},
	}

	cm, _ := c.makeConfigMap(context.TODO(), backend)

	assert.Equal(t, cm.Data["fluentd.conf"], `<source>
  @type tail
  path /data/log/*
  pos_file /data/pos
  tag host.path.*
  format json
  read_from_head true
  path_key path
</source>
<match **>
  @type elasticsearch
  host 127.0.0.1
  port 9200
  scheme http
  index_name fluentd
  log_es_400_reason true
  type_name _doc
  flush_interval 5s
  <buffer>
    flush_mode interval
    retry_type exponential_backoff
    total_limit_size 32MB
    chunk_limit_size 1MB
    chunk_full_threshold 0.8
    @type file
    path /var/log/td-agent/buffer/ccs.cluster.log_collector.buffer.audit-event-collector.host-path
    overflow_action block
    flush_interval 1s
    flush_thread_burst_interval 0.01
    chunk_limit_records 8000
   </buffer>
</match>
`)

	backend.ES.User = "user"
	backend.ES.Password = "cGFzc3dvcmQK"
	cm, _ = c.makeConfigMap(context.TODO(), backend)
	assert.Equal(t, cm.Data["fluentd.conf"], `<source>
  @type tail
  path /data/log/*
  pos_file /data/pos
  tag host.path.*
  format json
  read_from_head true
  path_key path
</source>
<match **>
  @type elasticsearch
  host 127.0.0.1
  port 9200
  scheme http
  index_name fluentd
  log_es_400_reason true
  type_name _doc
  user user
  password password

  flush_interval 5s
  <buffer>
    flush_mode interval
    retry_type exponential_backoff
    total_limit_size 32MB
    chunk_limit_size 1MB
    chunk_full_threshold 0.8
    @type file
    path /var/log/td-agent/buffer/ccs.cluster.log_collector.buffer.audit-event-collector.host-path
    overflow_action block
    flush_interval 1s
    flush_thread_burst_interval 0.01
    chunk_limit_records 8000
   </buffer>
</match>
`)
}
